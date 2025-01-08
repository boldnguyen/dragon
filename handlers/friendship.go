package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type AddFriendRequest struct {
	PlayerID int `json:"player_id"`
	FriendID int `json:"friend_id"`
}

type AcceptFriendRequest struct {
	PlayerID int `json:"player_id"`
	FriendID int `json:"friend_id"`
}

func AddFriend(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Kiểm tra xem người chơi đã có bạn chưa
	var existingFriendship models.Friendship
	if err := db.Where("player_id = ? AND friend_id = ?", req.PlayerID, req.FriendID).First(&existingFriendship).Error; err == nil {
		http.Error(w, "Already friends or request already sent", http.StatusConflict)
		return
	}

	// Thêm mối quan hệ kết bạn
	friendship := models.Friendship{
		PlayerID:  req.PlayerID,
		FriendID:  req.FriendID,
		Status:    "pending", // Chờ đợi chấp nhận
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := db.Create(&friendship).Error; err != nil {
		http.Error(w, "Failed to send friend request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Friend request sent successfully",
	})
}

func RemoveFriend(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Xóa mối quan hệ bạn bè
	if err := db.Where("player_id = ? AND friend_id = ?", req.PlayerID, req.FriendID).Delete(&models.Friendship{}).Error; err != nil {
		http.Error(w, "Failed to remove friend", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Friend removed successfully",
	})
}
func AcceptFriend(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req AcceptFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Kiểm tra xem mối quan hệ đã có tồn tại và có trạng thái 'pending'
	var friendship models.Friendship
	if err := db.Where("player_id = ? AND friend_id = ? AND status = ?", req.FriendID, req.PlayerID, "pending").First(&friendship).Error; err != nil {
		http.Error(w, "Friend request not found or already accepted", http.StatusNotFound)
		return
	}

	// Cập nhật trạng thái mối quan hệ từ 'pending' thành 'accepted'
	friendship.Status = "accepted"
	if err := db.Save(&friendship).Error; err != nil {
		http.Error(w, "Failed to accept friend request", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Friend request accepted successfully",
	})
}

func GetFriends(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "player_id is required", http.StatusBadRequest)
		return
	}

	// Lấy danh sách bạn bè
	var friends []models.Profile
	if err := db.Table("profiles").
		Joins("JOIN friendships ON friendships.friend_id = profiles.player_id").
		Where("friendships.player_id = ? AND friendships.status = ?", playerID, "accepted").
		Find(&friends).Error; err != nil {
		http.Error(w, "Failed to fetch friends", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(friends)
}

type GiftRequest struct {
	PlayerID int    `json:"player_id"`
	FriendID int    `json:"friend_id"`
	Gift     string `json:"gift"`
}

func SendGift(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req GiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Kiểm tra mối quan hệ bạn bè
	var friendship models.Friendship
	if err := db.Where("player_id = ? AND friend_id = ? AND status = ?", req.PlayerID, req.FriendID, "accepted").First(&friendship).Error; err != nil {
		http.Error(w, "Not friends or pending request", http.StatusBadRequest)
		return
	}

	// Lưu thông tin quà tặng vào bảng Gifts
	gift := models.Gift{
		PlayerID:  req.PlayerID,
		FriendID:  req.FriendID,
		Gift:      req.Gift,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := db.Create(&gift).Error; err != nil {
		http.Error(w, "Failed to send gift", http.StatusInternalServerError)
		return
	}

	// Phản hồi thành công
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Gift sent successfully",
	})
}
