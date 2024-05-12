package routes

import (
	"rakamin/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	userController := controllers.NewUserController()
	photoController := controllers.NewPhotoController()

	// User Endpoints
	router.POST("/users/register", userController.Register)
	router.POST("/users/login", userController.Login)
	router.PUT("/users/:userId", userController.UpdateUser)
	router.DELETE("/users/:userId", userController.DeleteUser)

	// Photo Endpoints
	router.POST("/photos", photoController.UploadPhoto)
	router.GET("/photos", photoController.GetPhotos)
	router.PUT("/photos/:photoId", photoController.UpdatePhoto)
	router.DELETE("/photos/:photoId", photoController.DeletePhoto)

	return router
}
