package controllers

import (
	"errors"
	"net/http"
	"path/filepath"
	"rakamin/database"
	"rakamin/helpers"
	"rakamin/models"
	"time"

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

	// Cek dan upload photo
	var updatedPhoto *models.Photo
	if _, _, err := c.Request.FormFile("file"); err == nil {
		updatedPhoto, err = uploadAndCreatePhoto(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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

func uploadAndCreatePhoto(c *gin.Context) (*models.Photo, error) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		return nil, errors.New("error")
	}
	defer file.Close()

	filename := filepath.Base(fileHeader.Filename)
	if err := c.SaveUploadedFile(fileHeader, filepath.Join("./uploads", filename)); err != nil {
		return nil, errors.New("gagal menyimpan gambar")
	}

	photo := &models.Photo{
		Title:     filename,
		Caption:   "caption " + filename,
		PhotoURL:  filepath.Join("./uploads", filename),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if _, err := govalidator.ValidateStruct(photo); err != nil {
		return nil, errors.New("invalid photo data")
	}

	return photo, nil
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
