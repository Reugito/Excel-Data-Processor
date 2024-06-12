package services

import (
	"bytes"
	"context"
	"dataProcessor/cache"
	"dataProcessor/database"
	"dataProcessor/models"
	"dataProcessor/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"log"
	"sync"
)

type ContactService struct {
	MySQLRepo *database.MySQLRepo
	RedisRepo *cache.RedisRepo
}

var (
	validate = validator.New()
	ctx      = context.Background()
)

func (serv *ContactService) ProcessFile(fileData []byte, errChan chan error) {
	log.Println("Processing file starts")
	defer log.Println("Processing file ends")

	xlFile, err := excelize.OpenReader(bytes.NewReader(fileData))
	if err != nil {
		errChan <- fmt.Errorf("failed to open file: %w", err)
		return
	}

	sheetName := xlFile.GetSheetName(0)
	rows, err := xlFile.GetRows(sheetName)
	if err != nil {
		errChan <- fmt.Errorf("failed to get rows: %w", err)
		return
	}

	contactsChan := make(chan models.Contact, 100)
	doneChan := make(chan struct{})
	processingErrChan := make(chan error, 1)

	const numWorkers = 10
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for contact := range contactsChan {
				contact.RefId = uuid.New().String()
				if err := validate.Struct(contact); err != nil {
					processingErrChan <- fmt.Errorf("validation failed: %s", utils.FormatValidationError(err))
					return
				}
				// Store in DB
				if err := serv.MySQLRepo.InsertContact(contact); err != nil {
					processingErrChan <- fmt.Errorf("failed to insert contact into DB: %w", err)
					return
				}
				// Store in Redis
				if err := serv.RedisRepo.CacheContact(contact); err != nil {
					processingErrChan <- fmt.Errorf("failed to cache contact in Redis: %w", err)
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	for rowIdx, row := range rows {
		if rowIdx == 0 {
			if !utils.ValidateHeaders(row) {
				errChan <- fmt.Errorf("invalid headers structure received in file")
				return
			}
			continue
		}

		contact := utils.ParseRow(row)
		contactsChan <- contact
	}

	close(contactsChan)
	<-doneChan

	select {
	case err := <-processingErrChan:
		errChan <- err
	default:
		errChan <- nil
	}
}

func (serv *ContactService) GetAllContacts() ([]models.Contact, error) {
	log.Println("Getting all contacts")
	defer log.Println("Getting all contacts ends")
	var contacts []models.Contact
	var cursor uint64
	pattern := "contacts_*"

	for {
		keys, nextCursor, err := serv.RedisRepo.Client.Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			contact, err := serv.RedisRepo.GetContactFromCache(key)
			if err != nil {
				fmt.Printf("Failed to get contact for key %s: %s\n", key, err)
				continue
			}
			contacts = append(contacts, contact)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	if len(contacts) == 0 {
		log.Println("No contacts found in cache, fetching from database...")
		contactsInDB, err := serv.MySQLRepo.FindAllContacts()
		if err != nil {
			return nil, err
		}
		contacts = contactsInDB
	}

	return contacts, nil
}

func (serv *ContactService) UpdateContact(refId string, updatedContact *models.Contact) error {
	log.Println("Updating contact starts")
	defer log.Println("Updating contact ends")

	key := "contacts_" + refId

	contact, err := serv.RedisRepo.GetContactFromCache(key)
	if err != nil && err != redis.Nil {
		return err
	} else if err == redis.Nil {
		contact, err = serv.MySQLRepo.FindContactByRefId(refId)
		if err != nil {
			return err
		}
	}
	if err := validate.Struct(updatedContact); err != nil {
		return fmt.Errorf("validation failed:: %s", utils.FormatValidationError(err))
	}

	if contact.RefId == "" {
		return fmt.Errorf("contact not found")
	}

	if err := serv.MySQLRepo.UpdateContactByRefID(refId, updatedContact); err != nil {
		return fmt.Errorf("failed to update contact in DB: %w", err)
	}

	// Update in Redis
	updatedContact.RefId = refId
	if err := serv.RedisRepo.UpdateCacheByRefID(*updatedContact); err != nil {
		return fmt.Errorf("failed to update contact in Redis: %w", err)
	}
	return nil
}

func (serv *ContactService) FindContactByRefIdContact(refId string) (models.Contact, error) {
	log.Println("Finding contact starts")
	defer log.Println("Finding contact ends")

	key := "contacts_" + refId

	contact, err := serv.RedisRepo.GetContactFromCache(key)
	if err != nil && err != redis.Nil {
		return contact, fmt.Errorf("contact not found: %w", err)
	} else if err == redis.Nil {
		contact, err = serv.MySQLRepo.FindContactByRefId(refId)
		if err != nil {
			return contact, fmt.Errorf("contact not found: %s", err.Error())
		}
		serv.RedisRepo.CacheContact(contact)
		if err != nil {
			return contact, fmt.Errorf("failed to insert data in cache: %w", err)
		}
	}
	return contact, nil
}
