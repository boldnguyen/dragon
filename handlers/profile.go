package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type CreateProfileRequest struct {
	PlayerID int    `json:"player_id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar,omitempty"`
}

func CreateProfile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req CreateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if profile with the same PlayerID already exists
	var existingProfile models.Profile
	if err := db.Where("player_id = ?", req.PlayerID).First(&existingProfile).Error; err == nil {
		http.Error(w, "Profile already exists", http.StatusConflict)
		return
	}

	profile := models.Profile{
		PlayerID:     req.PlayerID,
		Name:         req.Name,
		Avatar:       req.Avatar,
		Level:        1, // Default level
		Wins:         0,
		Losses:       0,
		DragonsOwned: 0,
		TotalTokens:  0,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	if err := db.Create(&profile).Error; err != nil {
		http.Error(w, "Failed to create profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Profile created successfully",
		"profileID": profile.ID,
	})
}

func GetProfile(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	var profile models.Profile
	if err := db.Where("player_id = ?", playerID).First(&profile).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profile)
}
