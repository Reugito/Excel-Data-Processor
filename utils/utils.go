package utils

import (
	"dataProcessor/models"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) string {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err.Error()
	}

	validationErrors := err.(validator.ValidationErrors)
	firstError := validationErrors[0]
	return fmt.Sprintf("Field '%s' is '%s'", firstError.Field(), firstError.Tag())
}

func ParseRow(row []string) models.Contact {
	contact := models.Contact{}
	for i, cell := range row {
		switch i {
		case 0:
			contact.FirstName = cell
		case 1:
			contact.LastName = cell
		case 2:
			contact.CompanyName = cell
		case 3:
			contact.Address = cell
		case 4:
			contact.City = cell
		case 5:
			contact.County = cell
		case 6:
			contact.Postal = cell
		case 7:
			contact.Phone = cell
		case 8:
			contact.Email = cell
		case 9:
			contact.Web = cell
		}
	}
	return contact
}

func ValidateHeaders(row []string) bool {
	expectedHeaders := []string{"first_name", "last_name", "company_name", "address", "city", "county", "postal", "phone", "email", "web"}
	for i, cell := range row {
		if cell != expectedHeaders[i] {
			return false
		}
	}
	return true
}
