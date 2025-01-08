// handlers/message.go
package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type SendMessageRequest struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"` // Có thể là bạn bè hoặc nhóm chat
	Content    string `json:"content"`
	ChatType   string `json:"chat_type"` // Loại chat: world, clan, friend, team
}

func SendMessage(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	message := models.Message{
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		ChatType:   req.ChatType,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	if err := db.Create(&message).Error; err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Message sent successfully",
	})
}

func GetMessages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	chatType := r.URL.Query().Get("chat_type")
	receiverID := r.URL.Query().Get("receiver_id")

	if chatType == "" || receiverID == "" {
		http.Error(w, "chat_type and receiver_id are required", http.StatusBadRequest)
		return
	}

	// Lấy danh sách tin nhắn theo chat_type và receiver_id
	var messages []models.Message
	if err := db.Where("chat_type = ? AND receiver_id = ?", chatType, receiverID).Find(&messages).Error; err != nil {
		http.Error(w, "Failed to fetch messages", http.StatusInternalServerError)
		return
	}

	// Lấy thông tin profile cho người gửi và người nhận
	for i, message := range messages {
		// Lấy thông tin người gửi
		var senderProfile models.Profile
		if err := db.Where("player_id = ?", message.SenderID).First(&senderProfile).Error; err == nil {
			messages[i].SenderName = senderProfile.Name
			messages[i].SenderAvatar = senderProfile.Avatar
		}

		// Lấy thông tin người nhận
		var receiverProfile models.Profile
		if err := db.Where("player_id = ?", message.ReceiverID).First(&receiverProfile).Error; err == nil {
			messages[i].ReceiverName = receiverProfile.Name
			messages[i].ReceiverAvatar = receiverProfile.Avatar
		}
	}

	// Trả về danh sách tin nhắn
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}
