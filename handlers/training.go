package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type StartTrainingRequest struct {
	PlayerID     int    `json:"player_id"`     // ID của người chơi
	DragonID     uint   `json:"dragon_id"`     // ID của rồng
	TrainingType string `json:"training_type"` // Loại huấn luyện (attack, defense)
	UseToken     bool   `json:"use_token"`     // Có sử dụng token để tăng tốc không
}

// StartTraining xử lý việc bắt đầu huấn luyện
func StartTraining(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req StartTrainingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Kiểm tra xem người chơi có sở hữu rồng không
	var dragon models.Dragon
	if err := db.Where("id = ? AND player_id = ?", req.DragonID, req.PlayerID).First(&dragon).Error; err != nil {
		http.Error(w, "Dragon not found or does not belong to the player", http.StatusNotFound)
		return
	}

	// Nếu sử dụng token, hoàn thành huấn luyện ngay lập tức
	if req.UseToken {
		// Kiểm tra số token của người chơi
		var profile models.Profile
		if err := db.Where("player_id = ?", req.PlayerID).First(&profile).Error; err != nil {
			http.Error(w, "Player profile not found", http.StatusNotFound)
			return
		}

		// Giả sử huấn luyện tốn 10 token
		if profile.TotalTokens < 10 {
			http.Error(w, "Not enough tokens", http.StatusBadRequest)
			return
		}

		// Trừ token của người chơi
		profile.TotalTokens -= 10
		if err := db.Save(&profile).Error; err != nil {
			http.Error(w, "Failed to update player profile", http.StatusInternalServerError)
			return
		}

		// Tăng thuộc tính của rồng ngay lập tức
		switch req.TrainingType {
		case "attack":
			dragon.Attack += 2 // Tăng attack thêm 2 điểm
		case "defense":
			dragon.Defense += 2 // Tăng defense thêm 2 điểm
		default:
			http.Error(w, "Invalid training type", http.StatusBadRequest)
			return
		}

		// Lưu thông tin cập nhật của rồng
		if err := db.Save(&dragon).Error; err != nil {
			http.Error(w, "Failed to update dragon", http.StatusInternalServerError)
			return
		}

		// Trả về phản hồi thành công
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"message":      "Training completed instantly using tokens",
			"dragon_id":    dragon.ID,
			"new_attack":   dragon.Attack,
			"new_defense":  dragon.Defense,
			"tokens_spent": 10,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	// Nếu không sử dụng token, tạo bản ghi huấn luyện và đợi thời gian
	trainingTime := int64(300) // Mặc định 5 phút (300 giây)

	// Tạo bản ghi huấn luyện
	training := models.Training{
		PlayerID:     req.PlayerID,
		DragonID:     req.DragonID,
		TrainingType: req.TrainingType,
		StartTime:    time.Now().Unix(),
		EndTime:      time.Now().Unix() + trainingTime,
		UseToken:     req.UseToken,
		CreatedAt:    time.Now().Unix(),
		UpdatedAt:    time.Now().Unix(),
	}

	// Lưu bản ghi huấn luyện vào cơ sở dữ liệu
	if err := db.Create(&training).Error; err != nil {
		http.Error(w, "Failed to start training", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message":       "Training started",
		"training_id":   training.ID,
		"end_time":      training.EndTime,
		"training_type": training.TrainingType,
	}

	json.NewEncoder(w).Encode(response)
}
