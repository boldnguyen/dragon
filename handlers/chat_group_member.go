package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type AddMemberRequest struct {
	GroupID  uint `json:"group_id"`  // ID của nhóm chat
	PlayerID int  `json:"player_id"` // ID của người chơi
}

func AddMemberToChatGroup(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Kiểm tra nếu player đã tồn tại trong group
	var existingMember models.ChatGroupMember
	if err := db.Where("group_id = ? AND player_id = ?", req.GroupID, req.PlayerID).First(&existingMember).Error; err == nil {
		http.Error(w, "Player is already a member of the group", http.StatusConflict)
		return
	} else if err != gorm.ErrRecordNotFound { // Kiểm tra lỗi khác ngoài "không tìm thấy"
		http.Error(w, "Failed to check group membership", http.StatusInternalServerError)
		return
	}

	// Thêm player vào group
	member := models.ChatGroupMember{
		GroupID:   req.GroupID,
		PlayerID:  req.PlayerID,
		JoinedAt:  time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := db.Create(&member).Error; err != nil {
		http.Error(w, "Failed to add member to chat group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(member)
}

func GetChatGroupMembers(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	groupID := r.URL.Query().Get("group_id")
	if groupID == "" {
		http.Error(w, "group_id is required", http.StatusBadRequest)
		return
	}

	var members []models.ChatGroupMember
	if err := db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		http.Error(w, "Failed to fetch chat group members", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(members)
}

func RemoveMemberFromChatGroup(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := db.Where("group_id = ? AND player_id = ?", req.GroupID, req.PlayerID).Delete(&models.ChatGroupMember{}).Error; err != nil {
		http.Error(w, "Failed to remove member from chat group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Member removed successfully",
	})
}
