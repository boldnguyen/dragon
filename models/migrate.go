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

func SeedMap(db *gorm.DB) {
	stages := []Stage{
		{Name: "Forest", Description: "A mystical forest with many secrets"},
		{Name: "Mountain", Description: "A high mountain with dangerous cliffs"},
		{Name: "Desert", Description: "A vast desert with hidden treasures"},
	}

	for _, stage := range stages {
		if err := db.FirstOrCreate(&stage, Stage{Name: stage.Name}).Error; err != nil {
			panic("Failed to seed stages: " + err.Error())
		}

		// Thêm vòng cho mỗi stage
		rounds := []Round{
			{Name: "Round 1", Description: "The first challenge of the stage"},
			{Name: "Round 2", Description: "The second challenge of the stage"},
		}

		for _, round := range rounds {
			round.StageID = stage.ID
			if err := db.Create(&round).Error; err != nil {
				panic("Failed to seed rounds: " + err.Error())
			}

			// Thêm nhiệm vụ cho mỗi vòng
			missions := []Mission{
				{Name: "Collect 10 herbs", Description: "Find and collect 10 herbs", Reward: "50 Gold"},
				{Name: "Defeat the monster", Description: "Defeat the monster guarding the area", Reward: "100 Gold"},
			}

			for _, mission := range missions {
				mission.RoundID = round.ID
				if err := db.Create(&mission).Error; err != nil {
					panic("Failed to seed missions: " + err.Error())
				}
			}
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
		&Stage{},
		&Round{},
		&Mission{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed dữ liệu khởi tạo
	SeedItems(db)
	SeedEggs(db) // Gọi SeedEggs để seed trứng vào cơ sở dữ liệu
	SeedMap(db)  // Gọi SeedMap để seed dữ liệu bản đồ vào cơ sở dữ liệu
}
