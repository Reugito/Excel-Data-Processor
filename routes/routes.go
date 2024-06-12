package routes

import (
	"dataProcessor/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, dataProcesses *controllers.Application) {

	basePath := router.Group("/v1")

	route := basePath.Group("/contact")

	route.POST("/upload", dataProcesses.ContactController.UploadFileHandler)

	route.GET("/getAll", dataProcesses.ContactController.GetAllContactsHandler)

	route.PUT("/update/:refId", dataProcesses.ContactController.UpdateContactHandler)

	route.GET("/getByRefId/:refId", dataProcesses.ContactController.GetContactsByRefIdHandler)

	route.GET("/", dataProcesses.ContactController.SayHello)
}
