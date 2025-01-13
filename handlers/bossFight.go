package handlers

import (
	"dragon/models"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type PvEBossRequest struct {
	PlayerID int  `json:"player_id"`
	BossID   int  `json:"boss_id"`
	DragonID uint `json:"dragon_id"`
}

func FightBoss(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req PvEBossRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lấy thông tin người chơi và rồng
	var profile models.Profile
	if err := db.Where("player_id = ?", req.PlayerID).First(&profile).Error; err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	var dragon models.Dragon
	if err := db.Where("id = ? AND player_id = ?", req.DragonID, req.PlayerID).First(&dragon).Error; err != nil {
		http.Error(w, "Dragon not found", http.StatusNotFound)
		return
	}

	// Lấy thông tin boss
	var boss models.Boss
	if err := db.Where("id = ?", req.BossID).First(&boss).Error; err != nil {
		http.Error(w, "Boss not found", http.StatusNotFound)
		return
	}

	// Tính toán kết quả trận đấu
	playerPower := dragon.Attack + dragon.Defense
	bossPower := boss.Attack + boss.Defense

	var win bool
	var reward string
	if playerPower > bossPower {
		win = true
		profile.Wins++
		profile.TotalTokens += 500 // Boss battle thưởng lớn hơn PvE
		reward = "500 tokens"

		// Thêm phần thưởng rồng ngẫu nhiên với tỉ lệ 1/20
		rand.Seed(time.Now().UnixNano()) // Khởi tạo bộ sinh số ngẫu nhiên

		if rand.Intn(20) == 0 { // Tỉ lệ 1/20
			// Tạo rồng ngẫu nhiên
			dragonName := generateRandomDragonName()
			attack, defense := generateRandomAttributes(generateRandomRarity())
			rarity := generateRandomRarity()

			newDragon := models.Dragon{
				PlayerID:  req.PlayerID,
				Name:      dragonName,
				Rarity:    rarity,
				Level:     1,
				Attack:    attack,
				Defense:   defense,
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			}

			if err := db.Create(&newDragon).Error; err != nil {
				http.Error(w, "Failed to create special dragon", http.StatusInternalServerError)
				return
			}

			reward = "Special Dragon Egg and 500 tokens"
		} else {
			reward = "500 tokens" // Không có phần thưởng rồng, chỉ có tokens
		}

		// Thêm kinh nghiệm cho rồng
		dragon.Experience += 200

		// Kiểm tra và nâng cấp cấp độ của rồng nếu đạt đủ kinh nghiệm
		for dragon.Experience >= dragon.Level*300 { // Cấp độ boss có thể yêu cầu nhiều kinh nghiệm hơn PvE
			dragon.Experience -= dragon.Level * 300
			dragon.Level++

			// Tăng cường các chỉ số của rồng khi lên cấp
			dragon.Attack += 20  // Tăng điểm tấn công nhiều hơn
			dragon.Defense += 10 // Tăng điểm phòng thủ nhiều hơn
		}

		// Lưu thông tin cập nhật của rồng
		if err := db.Save(&dragon).Error; err != nil {
			http.Error(w, "Failed to update dragon level", http.StatusInternalServerError)
			return
		}

	} else {
		win = false
		profile.Losses++
		reward = "No reward"
	}

	// Cập nhật lại profile
	if err := db.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Trả về kết quả chi tiết
	response := map[string]interface{}{
		"message":     "Boss battle completed",
		"win":         win,
		"playerPower": playerPower,
		"bossPower":   bossPower,
		"reward":      reward,
		"playerProfile": map[string]interface{}{
			"wins":        profile.Wins,
			"losses":      profile.Losses,
			"totalTokens": profile.TotalTokens,
		},
		"dragonProfile": map[string]interface{}{
			"level":      dragon.Level,
			"experience": dragon.Experience,
			"attack":     dragon.Attack,
			"defense":    dragon.Defense,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Định dạng kết quả trả về
	indentData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	w.Write(indentData)
}
