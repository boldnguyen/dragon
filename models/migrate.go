package models

import (
	"log"

	"gorm.io/gorm"
)

// AutoMigrate thực hiện migration cho tất cả các model
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Wallet{},
		&Profile{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
