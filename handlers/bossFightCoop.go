package handlers

import (
	"dragon/models"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type PvEBossCoopRequest struct {
	Player1ID int  `json:"player1_id"`
	Player2ID int  `json:"player2_id"`
	BossID    int  `json:"boss_id"`
	Dragon1ID uint `json:"dragon1_id"`
	Dragon2ID uint `json:"dragon2_id"`
}

func FightBossCoop(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req PvEBossCoopRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get information for the first player (Player 1)
	var profile1 models.Profile
	if err := db.Where("player_id = ?", req.Player1ID).First(&profile1).Error; err != nil {
		http.Error(w, "Player 1 not found", http.StatusNotFound)
		return
	}

	var dragon1 models.Dragon
	if err := db.Where("id = ? AND player_id = ?", req.Dragon1ID, req.Player1ID).First(&dragon1).Error; err != nil {
		http.Error(w, "Dragon 1 not found", http.StatusNotFound)
		return
	}

	// Get information for the second player (Player 2)
	var profile2 models.Profile
	if err := db.Where("player_id = ?", req.Player2ID).First(&profile2).Error; err != nil {
		http.Error(w, "Player 2 not found", http.StatusNotFound)
		return
	}

	var dragon2 models.Dragon
	if err := db.Where("id = ? AND player_id = ?", req.Dragon2ID, req.Player2ID).First(&dragon2).Error; err != nil {
		http.Error(w, "Dragon 2 not found", http.StatusNotFound)
		return
	}

	// Get boss information
	var boss models.Boss
	if err := db.Where("id = ?", req.BossID).First(&boss).Error; err != nil {
		http.Error(w, "Boss not found", http.StatusNotFound)
		return
	}

	// Calculate the combined power of both players
	playerPower1 := dragon1.Attack + dragon1.Defense
	playerPower2 := dragon2.Attack + dragon2.Defense
	totalPower := playerPower1 + playerPower2

	// Calculate the boss's power
	bossPower := boss.Attack + boss.Defense

	var win bool
	var reward string
	if totalPower > bossPower {
		win = true
		profile1.Wins++
		profile2.Wins++
		profile1.TotalTokens += 250 // Each player receives 250 tokens
		profile2.TotalTokens += 250
		reward = "500 tokens"

		// Add a random dragon reward with a 1/20 chance
		rand.Seed(time.Now().UnixNano()) // Initialize the random number generator

		if rand.Intn(20) == 0 { // 1/20 chance
			// Create a random dragon for both players
			dragonName := generateRandomDragonName()
			attack, defense := generateRandomAttributes(generateRandomRarity())
			rarity := generateRandomRarity()

			newDragon1 := models.Dragon{
				PlayerID:  req.Player1ID,
				Name:      dragonName,
				Rarity:    rarity,
				Level:     1,
				Attack:    attack,
				Defense:   defense,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}

			newDragon2 := models.Dragon{
				PlayerID:  req.Player2ID,
				Name:      dragonName,
				Rarity:    rarity,
				Level:     1,
				Attack:    attack,
				Defense:   defense,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}

			// Create the new dragons for both players
			if err := db.Create(&newDragon1).Error; err != nil {
				http.Error(w, "Failed to create special dragon for player 1", http.StatusInternalServerError)
				return
			}
			if err := db.Create(&newDragon2).Error; err != nil {
				http.Error(w, "Failed to create special dragon for player 2", http.StatusInternalServerError)
				return
			}

			reward = "Special Dragon Egg for both players and 500 tokens"
		} else {
			reward = "500 tokens"
		}

		// Add experience to both dragons
		dragon1.Experience += 200
		dragon2.Experience += 200

		// Check and level up both dragons if they have enough experience
		for dragon1.Experience >= dragon1.Level*300 {
			dragon1.Experience -= dragon1.Level * 300
			dragon1.Level++
			dragon1.Attack += 20
			dragon1.Defense += 10
		}

		for dragon2.Experience >= dragon2.Level*300 {
			dragon2.Experience -= dragon2.Level * 300
			dragon2.Level++
			dragon2.Attack += 20
			dragon2.Defense += 10
		}

		// Save the updated information for both dragons
		if err := db.Save(&dragon1).Error; err != nil {
			http.Error(w, "Failed to update dragon 1 level", http.StatusInternalServerError)
			return
		}

		if err := db.Save(&dragon2).Error; err != nil {
			http.Error(w, "Failed to update dragon 2 level", http.StatusInternalServerError)
			return
		}

	} else {
		win = false
		profile1.Losses++
		profile2.Losses++
		reward = "No reward"
	}

	// Update the profiles of both players
	if err := db.Save(&profile1).Error; err != nil {
		http.Error(w, "Failed to update player 1 profile", http.StatusInternalServerError)
		return
	}

	if err := db.Save(&profile2).Error; err != nil {
		http.Error(w, "Failed to update player 2 profile", http.StatusInternalServerError)
		return
	}

	// Return detailed results
	response := map[string]interface{}{
		"message":      "Boss battle completed (Co-op)",
		"win":          win,
		"player1Power": playerPower1,
		"player2Power": playerPower2,
		"bossPower":    bossPower,
		"reward":       reward,
		"player1Profile": map[string]interface{}{
			"wins":        profile1.Wins,
			"losses":      profile1.Losses,
			"totalTokens": profile1.TotalTokens,
		},
		"player2Profile": map[string]interface{}{
			"wins":        profile2.Wins,
			"losses":      profile2.Losses,
			"totalTokens": profile2.TotalTokens,
		},
		"dragon1Profile": map[string]interface{}{
			"level":      dragon1.Level,
			"experience": dragon1.Experience,
			"attack":     dragon1.Attack,
			"defense":    dragon1.Defense,
		},
		"dragon2Profile": map[string]interface{}{
			"level":      dragon2.Level,
			"experience": dragon2.Experience,
			"attack":     dragon2.Attack,
			"defense":    dragon2.Defense,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Format the response
	indentData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	w.Write(indentData)
}
