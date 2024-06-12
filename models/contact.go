package models

import (
	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	RefId       string `json:"ref_id" validate:"required,uuid4"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	CompanyName string `json:"company_name" validate:"required"`
	Address     string `json:"address" validate:"required"`
	City        string `json:"city" validate:"required"`
	County      string `json:"county" validate:"required"`
	Postal      string `json:"postal" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Web         string `json:"web" validate:"required,url"`
}
