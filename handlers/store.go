// handlers/store.go
package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type BuyItemRequest struct {
	PlayerID int  `json:"player_id"`
	ItemID   uint `json:"item_id"`
}

func GetItems(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var items []models.Item
	query := db
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&items).Error; err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
}

func BuyItem(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req BuyItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch the item
	var item models.Item
	if err := db.First(&item, req.ItemID).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Verify player's profile
	var profile models.Profile
	if err := db.Where("player_id = ?", req.PlayerID).First(&profile).Error; err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	// Verify player's wallet
	var wallet models.Wallet
	if err := db.Where("user_id = ?", req.PlayerID).First(&wallet).Error; err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Check if player has enough balance
	if item.Currency == "GOLD" && profile.TotalTokens < item.Price {
		http.Error(w, "Insufficient GOLD", http.StatusBadRequest)
		return
	} else if item.Currency == "DIAMOND" {
		// Logic for DIAMOND currency (if applicable)
		// Placeholder for now
	}

	// Deduct the price from the player's balance (profile and wallet)
	profile.TotalTokens -= item.Price
	if err := db.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update player's profile", http.StatusInternalServerError)
		return
	}

	// Deduct the price from the wallet balance
	wallet.Balance -= item.Price
	if err := db.Save(&wallet).Error; err != nil {
		http.Error(w, "Failed to update wallet balance", http.StatusInternalServerError)
		return
	}

	// Add item to the player's wallet items
	// Check if the item already exists in the wallet's items array
	var itemExists bool
	for _, existingItem := range wallet.Items {
		if existingItem == item.Name { // or you can use item.ID
			itemExists = true
			break
		}
	}

	if !itemExists {
		// Add item to wallet items array
		wallet.Items = append(wallet.Items, item.Name) // Use item.ID if you prefer to store IDs
		if err := db.Save(&wallet).Error; err != nil {
			http.Error(w, "Failed to update wallet items", http.StatusInternalServerError)
			return
		}
	}

	// Record the transaction
	transaction := models.Transaction{
		PlayerID:  req.PlayerID,
		ItemID:    item.ID,
		Price:     item.Price,
		Currency:  item.Currency,
		CreatedAt: time.Now().Unix(),
	}

	if err := db.Create(&transaction).Error; err != nil {
		http.Error(w, "Failed to record transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Item purchased successfully",
	})
}
