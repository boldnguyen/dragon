package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type PvERequest struct {
	PlayerID    int    `json:"player_id"`
	DragonLevel int    `json:"dragon_level"`
	EnemyID     int    `json:"enemy_id"`
	DragonID    uint   `json:"dragon_id,omitempty"`   // Hoặc dùng DragonName nếu cần
	DragonName  string `json:"dragon_name,omitempty"` // Dùng để chỉ định con rồng
}

func FightEnemy(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req PvERequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lấy thông tin người chơi từ bảng Profile
	var profile models.Profile
	if err := db.Where("player_id = ?", req.PlayerID).First(&profile).Error; err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	// Tìm con rồng mà người chơi chọn
	var dragon models.Dragon
	if req.DragonID != 0 {
		// Nếu có DragonID, tìm rồng theo ID
		if err := db.Where("id = ? AND player_id = ?", req.DragonID, req.PlayerID).First(&dragon).Error; err != nil {
			http.Error(w, "Dragon not found", http.StatusNotFound)
			return
		}
	} else if req.DragonName != "" {
		// Nếu có DragonName, tìm rồng theo tên
		if err := db.Where("name = ? AND player_id = ?", req.DragonName, req.PlayerID).First(&dragon).Error; err != nil {
			http.Error(w, "Dragon not found", http.StatusNotFound)
			return
		}
	} else {
		http.Error(w, "Dragon ID or name must be provided", http.StatusBadRequest)
		return
	}

	// Lấy thông tin đối thủ từ bảng Enemy
	var enemy models.Enemy
	if err := db.Where("id = ?", req.EnemyID).First(&enemy).Error; err != nil {
		http.Error(w, "Enemy not found", http.StatusNotFound)
		return
	}

	// Tính toán sức mạnh của người chơi từ con rồng đã chọn
	playerPower := dragon.Attack + dragon.Defense

	// Tính toán sức mạnh của đối thủ
	enemyPower := enemy.Level + enemy.Attack

	// Tính kết quả trận đấu
	var win bool
	if playerPower > enemyPower {
		win = true
		profile.Wins++
		profile.Experience += 100
		profile.TotalTokens += 50
	} else {
		win = false
		profile.Losses++
	}

	// Cập nhật profile
	if err := db.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Trả về kết quả trận đấu
	response := map[string]interface{}{
		"message":       "Battle completed",
		"win":           win,
		"player_exp":    profile.Experience,
		"player_tokens": profile.TotalTokens,
		"enemy":         enemy.Name,
		"reward":        enemy.Reward,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Định dạng kết quả trả về với indent
	indentData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	w.Write(indentData)
}
