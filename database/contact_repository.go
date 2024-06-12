package database

import (
	"dataProcessor/models"
	"gorm.io/gorm"
)

type MySQLRepo struct {
	DB *gorm.DB
}

func (repo *MySQLRepo) InsertContact(data models.Contact) error {
	return repo.DB.Create(&data).Error
}

func (repo *MySQLRepo) FindAllContacts() ([]models.Contact, error) {
	var data []models.Contact

	err := repo.DB.Find(&data).Error
	return data, err
}

func (repo *MySQLRepo) FindContactByRefId(refID string) (models.Contact, error) {
	var contact models.Contact
	if err := repo.DB.Where("ref_id = ?", refID).First(&contact).Error; err != nil {
		return contact, err
	}
	return contact, nil
}

func (repo *MySQLRepo) UpdateContactByRefID(refID string, updatedContact *models.Contact) error {
	var contact models.Contact
	if err := repo.DB.Where("ref_id = ?", refID).First(&contact).Error; err != nil {
		return err
	}

	// Update the contact fields
	contact.FirstName = updatedContact.FirstName
	contact.LastName = updatedContact.LastName
	contact.CompanyName = updatedContact.CompanyName
	contact.Address = updatedContact.Address
	contact.City = updatedContact.City
	contact.County = updatedContact.County
	contact.Postal = updatedContact.Postal
	contact.Phone = updatedContact.Phone
	contact.Email = updatedContact.Email
	contact.Web = updatedContact.Web

	return repo.DB.Save(&contact).Error
}

func (repo *MySQLRepo) DeleteData(id uint) error {
	return repo.DB.Delete(&models.Contact{}, id).Error
}
