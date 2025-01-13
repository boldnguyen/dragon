package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type PvERequest struct {
	PlayerID int  `json:"player_id"`
	EnemyID  int  `json:"enemy_id"`
	DragonID uint `json:"dragon_id"` // ID rồng cần được cung cấp rõ ràng
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

	// Tìm con rồng mà người chơi chọn dựa vào DragonID
	var dragon models.Dragon
	if err := db.Where("id = ? AND player_id = ?", req.DragonID, req.PlayerID).First(&dragon).Error; err != nil {
		http.Error(w, "Dragon not found", http.StatusNotFound)
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
		profile.TotalTokens += 50

		// Thêm kinh nghiệm cho rồng
		dragon.Experience += 100

		// Kiểm tra và nâng cấp cấp độ của rồng nếu đạt đủ kinh nghiệm
		for dragon.Experience >= dragon.Level*200 {
			dragon.Experience -= dragon.Level * 200
			dragon.Level++

			// Tăng cường các chỉ số của rồng khi lên cấp
			dragon.Attack += 10 // Tăng điểm tấn công
			dragon.Defense += 5 // Tăng điểm phòng thủ
		}

		// Lưu thông tin cập nhật của rồng
		if err := db.Save(&dragon).Error; err != nil {
			http.Error(w, "Failed to update dragon level", http.StatusInternalServerError)
			return
		}
	} else {
		win = false
		profile.Losses++
	}

	// Cập nhật profile người chơi
	if err := db.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Trả về kết quả trận đấu
	response := map[string]interface{}{
		"message":           "Battle completed",
		"win":               win,
		"player_exp":        profile.Experience,
		"player_tokens":     profile.TotalTokens,
		"enemy":             enemy.Name,
		"reward":            enemy.Reward,
		"dragon_level":      dragon.Level,
		"dragon_experience": dragon.Experience,
		"dragon_attack":     dragon.Attack,
		"dragon_defense":    dragon.Defense,
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
