package database

import (
	"dataProcessor/config"
	"dataProcessor/models"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func ConnectMySQL(cfg config.MySQLConfig) (*MySQLRepo, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	// Open a connection to MySQL using GORM
	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&models.Contact{}); err != nil {
		return nil, err
	}

	fmt.Printf("Connected to MySQL database %s\n", cfg.Database)

	return &MySQLRepo{DB: db}, nil
}
