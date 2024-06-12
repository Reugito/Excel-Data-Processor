package controllers

import (
	"bytes"
	"dataProcessor/models"
	"dataProcessor/services"
	"dataProcessor/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type ContactController struct {
	ContactService *services.ContactService
}
type Application struct {
	ContactController *ContactController
}

func (contactCnt *ContactController) UploadFileHandler(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "xlsx file not found")
		return
	}

	if filepath.Ext(file.Filename) != ".xlsx" {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid file format, only .xlsx files are accepted")
		return
	}

	fileData, err := file.Open()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Opening file failed")
		return
	}
	defer fileData.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(fileData)

	errChan := make(chan error)

	go contactCnt.ContactService.ProcessFile(buf.Bytes(), errChan)

	// Wait for the error response from the goroutine
	if err := <-errChan; err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "File uploaded successfully", nil)
}

func (contactCnt *ContactController) GetAllContactsHandler(ctx *gin.Context) {
	contacts, err := contactCnt.ContactService.GetAllContacts()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve data from database")
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved data", contacts)
}

func (contactCnt *ContactController) UpdateContactHandler(ctx *gin.Context) {
	refID := ctx.Param("refId")
	var updatedContact models.Contact

	if err := ctx.ShouldBindJSON(&updatedContact); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := contactCnt.ContactService.UpdateContact(refID, &updatedContact); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Contact updated successfully", nil)

}

func (contactCnt *ContactController) GetContactsByRefIdHandler(ctx *gin.Context) {
	refID := ctx.Param("refId")

	contact, err := contactCnt.ContactService.FindContactByRefIdContact(refID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Contact updated successfully", contact)
}

func (contactCnt *ContactController) SayHello(ctx *gin.Context) {
	utils.SuccessResponse(ctx, http.StatusOK, "Application Running", nil)
}
