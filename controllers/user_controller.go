package controllers

import (
	"net/http"
	"rakamin/database"
	"rakamin/helpers"
	"rakamin/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController() *UserController {
	db, _ := database.ConnectDB()
	return &UserController{DB: db}
}

func (uc *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = string(hashedPassword)

	if err := uc.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user berhasil melakukan registrasi", "user": user})
}

func (uc *UserController) Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	if err := uc.DB.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password tidak sesuai"})
		return
	}

	// Validasi password
	if !helpers.ValidatePassword(user.Password, existingUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email atau password tidak sesuai"})
		return
	}

	// Generate JWT token
	token, err := helpers.GenerateToken(existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user berhasil login", "token": token})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("userId")

	var user models.User
	if err := uc.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "?user tidak ditemukan"})
		return
	}

	var updatedUser models.User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate input using govalidator
	if _, err := govalidator.ValidateStruct(updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.DB.Model(&user).Updates(updatedUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user berhasil diupdate", "user": updatedUser})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("userId")

	var user models.User
	if err := uc.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user tidak ditemukan"})
		return
	}

	if err := uc.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user berhasil dihapus"})
}
