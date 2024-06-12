package cache

import (
	"context"
	"dataProcessor/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisService struct {
	Client *redis.Client
}

func (repo *RedisRepo) CacheContact(contact models.Contact) error {
	key := fmt.Sprintf("contacts_%s", contact.RefId)
	jsonData, err := json.Marshal(contact)
	if err != nil {
		return err
	}
	return repo.Client.Set(context.Background(), key, jsonData, 5*time.Minute).Err()
}

func (repo *RedisRepo) GetContactFromCache(key string) (models.Contact, error) {
	var contact models.Contact

	contactJSON, err := repo.Client.Get(context.Background(), key).Result()
	if err != nil {
		return contact, err
	}

	err = json.Unmarshal([]byte(contactJSON), &contact)
	if err != nil {
		return contact, err
	}

	return contact, nil
}

func (repo *RedisRepo) UpdateCacheByRefID(updatedContact models.Contact) error {

	key := fmt.Sprintf("contacts_%s", updatedContact.RefId)

	contact, err := repo.GetContactFromCache(key)
	if err != nil && err != redis.Nil {
		return err
	} else if err == redis.Nil {
		return errors.New("contact not found")
	}

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

	contactJSON, err := json.Marshal(contact)
	if err != nil {
		return fmt.Errorf("failed to marshal contact: %w", err)
	}
	err = repo.Client.Set(context.Background(), key, contactJSON, 5*time.Minute).Err()
	if err != nil {
		return fmt.Errorf("failed to set contact in Redis: %w", err)
	}
	return nil
}
