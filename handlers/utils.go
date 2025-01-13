package handlers

import (
	"math/rand"
	"time"
)

// Helper function để tạo tên rồng ngẫu nhiên
func generateRandomDragonName() string {
	names := []string{"Firestorm", "Thunderstrike", "Frostbite", "Blazewing", "Shadowflame", "Stoneheart"}
	rand.Seed(time.Now().UnixNano())
	return names[rand.Intn(len(names))]
}

// Helper function để chọn độ hiếm ngẫu nhiên
func generateRandomRarity() string {
	rareTypes := []string{"Common", "Uncommon", "Rare", "Epic", "Legendary"}
	rand.Seed(time.Now().UnixNano())
	return rareTypes[rand.Intn(len(rareTypes))]
}

// Helper function để tạo giá trị ngẫu nhiên cho sức mạnh tấn công và phòng thủ
func generateRandomAttributes(rarity string) (int, int) {
	var attack int
	var defense int

	// Random attack and defense based on rarity
	switch rarity {
	case "Common":
		attack = rand.Intn(10) + 5
		defense = rand.Intn(6) + 3
	case "Uncommon":
		attack = rand.Intn(15) + 10
		defense = rand.Intn(8) + 5
	case "Rare":
		attack = rand.Intn(20) + 15
		defense = rand.Intn(10) + 7
	case "Epic":
		attack = rand.Intn(25) + 20
		defense = rand.Intn(12) + 9
	case "Legendary":
		attack = rand.Intn(30) + 25
		defense = rand.Intn(15) + 12
	}

	return attack, defense
}
