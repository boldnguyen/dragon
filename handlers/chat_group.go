package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type CreateChatGroupRequest struct {
	Name string `json:"name"` // Tên nhóm
	Type string `json:"type"` // Loại nhóm (world, clan, etc.)
}

func CreateChatGroup(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req CreateChatGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	group := models.ChatGroup{
		Name:      req.Name,
		Type:      req.Type,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := db.Create(&group).Error; err != nil {
		http.Error(w, "Failed to create chat group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
}
