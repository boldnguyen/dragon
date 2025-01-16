package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type CreateClanRequest struct {
	PlayerID    int    `json:"player_id"`    // ID của người chơi tạo clan
	Name        string `json:"name"`        // Tên clan
	Description string `json:"description"` // Mô tả clan
}

// CreateClan xử lý việc tạo clan
func CreateClan(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req CreateClanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Tạo clan mới
	clan := models.Clan{
		Name:        req.Name,
		LeaderID:    req.PlayerID,
		Description: req.Description,
		CreatedAt:   time.Now().Unix(),
		UpdatedAt:   time.Now().Unix(),
	}

	// Lưu clan vào cơ sở dữ liệu
	if err := db.Create(&clan).Error; err != nil {
		http.Error(w, "Failed to create clan", http.StatusInternalServerError)
		return
	}

	// Thêm người chơi vào clan với vai trò leader
	clanMember := models.ClanMember{
		ClanID:    clan.ID,
		PlayerID:  req.PlayerID,
		Role:      "leader",
		JoinedAt:  time.Now().Unix(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := db.Create(&clanMember).Error; err != nil {
		http.Error(w, "Failed to add leader to clan", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Clan created successfully",
		"clan_id": clan.ID,
	}

	json.NewEncoder(w).Encode(response)
}

type JoinClanRequest struct {
	PlayerID int  `json:"player_id"` // ID của người chơi
	ClanID   uint `json:"clan_id"`   // ID của clan
}

// JoinClan xử lý việc tham gia clan
func JoinClan(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req JoinClanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Kiểm tra xem clan có tồn tại không
	var clan models.Clan
	if err := db.First(&clan, req.ClanID).Error; err != nil {
		http.Error(w, "Clan not found", http.StatusNotFound)
		return
	}

	// Thêm người chơi vào clan với vai trò member
	clanMember := models.ClanMember{
		ClanID:    req.ClanID,
		PlayerID:  req.PlayerID,
		Role:      "member",
		JoinedAt:  time.Now().Unix(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	if err := db.Create(&clanMember).Error; err != nil {
		http.Error(w, "Failed to join clan", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Joined clan successfully",
		"clan_id": req.ClanID,
	}

	json.NewEncoder(w).Encode(response)
}

// GetClanInfo xử lý việc lấy thông tin clan
func GetClanInfo(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	clanID := r.URL.Query().Get("clan_id")

	// Lấy thông tin clan
	var clan models.Clan
	if err := db.First(&clan, clanID).Error; err != nil {
		http.Error(w, "Clan not found", http.StatusNotFound)
		return
	}

	// Lấy danh sách thành viên trong clan
	var members []models.ClanMember
	if err := db.Where("clan_id = ?", clanID).Find(&members).Error; err != nil {
		http.Error(w, "Failed to fetch clan members", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"clan":    clan,
		"members": members,
	}

	json.NewEncoder(w).Encode(response)
}

type ClanChatRequest struct {
	ClanID   uint   `json:"clan_id"`   // ID của clan
	PlayerID int    `json:"player_id"` // ID của người chơi
	Message  string `json:"message"`   // Nội dung tin nhắn
}

// ClanChat xử lý việc gửi tin nhắn trong clan
func ClanChat(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req ClanChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Tạo tin nhắn mới
	clanChat := models.ClanChat{
		ClanID:    req.ClanID,
		PlayerID:  req.PlayerID,
		Message:   req.Message,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Lưu tin nhắn vào cơ sở dữ liệu
	if err := db.Create(&clanChat).Error; err != nil {
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Message sent successfully",
	}

	json.NewEncoder(w).Encode(response)
}

// GetClanRanking xử lý việc lấy xếp hạng clan
func GetClanRanking(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	// Lấy danh sách clan và xếp hạng dựa trên số thành viên hoặc điểm số
	var clans []models.Clan
	if err := db.Order("id ASC").Find(&clans).Error; err != nil {
		http.Error(w, "Failed to fetch clan ranking", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"ranking": clans,
	}

	json.NewEncoder(w).Encode(response)
}

type CompleteClanTaskRequest struct {
	ClanID uint `json:"clan_id"` // ID của clan
	TaskID uint `json:"task_id"` // ID của nhiệm vụ
}

// CompleteClanTask xử lý việc hoàn thành nhiệm vụ clan
func CompleteClanTask(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req CompleteClanTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lấy thông tin nhiệm vụ
	var task models.ClanTask
	if err := db.First(&task, req.TaskID).Error; err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Đánh dấu nhiệm vụ đã hoàn thành
	task.Completed = true
	if err := db.Save(&task).Error; err != nil {
		http.Error(w, "Failed to complete task", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Task completed successfully",
		"reward":  task.Reward,
	}

	json.NewEncoder(w).Encode(response)
}