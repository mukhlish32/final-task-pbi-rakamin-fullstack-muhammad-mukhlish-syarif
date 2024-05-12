package controllers

import (
	"net/http"
	"rakamin/database"
	"rakamin/helpers"
	"rakamin/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PhotoController struct {
	DB *gorm.DB
}

func NewPhotoController() *PhotoController {
	db, _ := database.ConnectDB()
	return &PhotoController{DB: db}
}

func (pc *PhotoController) UploadPhoto(c *gin.Context) {
	var photo models.Photo
	if err := c.BindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "foto berhasil diupload", "photo": photo})
}

func (pc *PhotoController) GetPhotos(c *gin.Context) {
	var photos []models.Photo
	if err := pc.DB.Find(&photos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"photos": photos})
}

func (pc *PhotoController) UpdatePhoto(c *gin.Context) {
	photoID := c.Param("photoId")

	var photo models.Photo
	if err := pc.DB.Where("id = ?", photoID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "foto tidak ditemukan"})
		return
	}

	// Mengecek permission pada user
	if err := helpers.CheckPhotoPermission(c, photo.UserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var updatedPhoto models.Photo
	if err := c.BindJSON(&updatedPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := govalidator.ValidateStruct(updatedPhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.DB.Model(&photo).Updates(updatedPhoto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "foto berhasil diubah", "photo": updatedPhoto})
}

func (pc *PhotoController) DeletePhoto(c *gin.Context) {
	photoID := c.Param("photoId")

	var photo models.Photo
	if err := pc.DB.Where("id = ?", photoID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "foto tidak ditemukan"})
		return
	}

	// Mengecek permission pada user
	if err := helpers.CheckPhotoPermission(c, photo.UserID); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := pc.DB.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "foto berhasil dihapus"})
}
