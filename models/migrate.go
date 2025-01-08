package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

func SeedItems(db *gorm.DB) {
	items := []Item{
		{Name: "Starter Dragon Egg", Category: "Eggs", Price: 100, Currency: "Gold"},
		{Name: "Advanced Equipment", Category: "Equipment", Price: 300, Currency: "Gold"},
		{Name: "Energy Potion", Category: "Energy", Price: 50, Currency: "Diamond"},
		{Name: "Basic Magic Book", Category: "Magic Book", Price: 200, Currency: "Gold"},
		{Name: "10,000 Gold", Category: "GOLD & DIAMOND", Price: 10, Currency: "USD"},
		{Name: "1,000 Diamonds", Category: "GOLD & DIAMOND", Price: 15, Currency: "USD"},
	}

	for _, item := range items {
		if err := db.FirstOrCreate(&item, Item{Name: item.Name}).Error; err != nil {
			panic("Failed to seed items: " + err.Error())
		}
	}
}

func SeedEggs(db *gorm.DB) {
	eggs := []Egg{
		{Name: "Starter Dragon Egg", Rarity: "Common", Price: 100, Currency: "GOLD", HatchTime: 300, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
		{Name: "Advanced Dragon Egg", Rarity: "Rare", Price: 500, Currency: "GOLD", HatchTime: 600, CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()},
	}

	for _, egg := range eggs {
		if err := db.FirstOrCreate(&egg, Egg{Name: egg.Name}).Error; err != nil {
			panic("Failed to seed eggs: " + err.Error())
		}
	}
}

// AutoMigrate thực hiện migration cho tất cả các model
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&Wallet{},
		&Profile{},
		&Friendship{},
		&Gift{},
		&Message{},
		&ChatGroup{},
		&ChatGroupMember{},
		&Item{},
		&Transaction{},
		&Egg{},
		&Dragon{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed dữ liệu khởi tạo
	SeedItems(db)
	SeedEggs(db) // Gọi SeedEggs để seed trứng vào cơ sở dữ liệu
}
