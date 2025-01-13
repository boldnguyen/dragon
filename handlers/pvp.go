package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type PvPRequest struct {
	PlayerID  int     `json:"player_id"`
	DragonID  uint    `json:"dragon_id"`
	BetAmount float64 `json:"bet_amount"`
}

func PvPMatchHandler(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req PvPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lấy thông tin người chơi và rồng
	var player1 models.Profile
	var dragon1 models.Dragon
	if err := db.Where("player_id = ?", req.PlayerID).First(&player1).Error; err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}
	if err := db.Where("id = ? AND player_id = ?", req.DragonID, req.PlayerID).First(&dragon1).Error; err != nil {
		http.Error(w, "Dragon not found", http.StatusNotFound)
		return
	}

	// Kiểm tra số dư token của người chơi
	if player1.TotalTokens < req.BetAmount {
		http.Error(w, "Insufficient tokens", http.StatusBadRequest)
		return
	}

	// Tìm đối thủ phù hợp
	var dragon2 models.Dragon
	err := db.Where("level = ? AND ABS(attack - ?) <= 10 AND ABS(defense - ?) <= 10 AND player_id != ?",
		dragon1.Level, dragon1.Attack, dragon1.Defense, req.PlayerID).Order("RANDOM()").First(&dragon2).Error
	if err != nil {
		http.Error(w, "No suitable opponent found", http.StatusNotFound)
		return
	}

	var player2 models.Profile
	if err := db.Where("player_id = ?", dragon2.PlayerID).First(&player2).Error; err != nil {
		http.Error(w, "Opponent not found", http.StatusNotFound)
		return
	}

	// Tính toán sức mạnh và kết quả
	player1Strength := dragon1.Attack - dragon1.Defense
	player2Strength := dragon2.Attack - dragon2.Defense

	// Kiểm tra nếu trận đấu hòa
	if player1Strength == player2Strength {
		player1.TotalTokens += req.BetAmount
		player2.TotalTokens += req.BetAmount
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		// Hoàn lại token cho cả hai người chơi
		if err := tx.Model(&models.Profile{}).Where("player_id = ?", player1.PlayerID).Update("total_tokens", player1.TotalTokens).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to refund tokens for player1", http.StatusInternalServerError)
			return
		}
		if err := tx.Model(&models.Profile{}).Where("player_id = ?", player2.PlayerID).Update("total_tokens", player2.TotalTokens).Error; err != nil {
			tx.Rollback()
			http.Error(w, "Failed to refund tokens for player2", http.StatusInternalServerError)
			return
		}

		tx.Commit()

		// Trả về kết quả hòa
		response := map[string]interface{}{
			"message": "It's a draw",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	var winner *models.Profile
	var loser *models.Profile
	if player1Strength > player2Strength {
		winner = &player1
		loser = &player2
	} else {
		winner = &player2
		loser = &player1
	}

	// Cập nhật token và lưu kết quả
	reward := req.BetAmount * 2
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	winner.TotalTokens += reward
	loser.TotalTokens -= req.BetAmount

	if err := tx.Model(&models.Profile{}).Where("player_id = ?", winner.PlayerID).Update("total_tokens", winner.TotalTokens).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update winner tokens", http.StatusInternalServerError)
		return
	}
	if err := tx.Model(&models.Profile{}).Where("player_id = ?", loser.PlayerID).Update("total_tokens", loser.TotalTokens).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update loser tokens", http.StatusInternalServerError)
		return
	}

	match := models.PvPMatch{
		Player1ID:       player1.PlayerID,
		Player2ID:       player2.PlayerID,
		Player1DragonID: dragon1.ID,
		Player2DragonID: dragon2.ID,
		BetAmount:       req.BetAmount,
		WinnerID:        winner.PlayerID,
		Reward:          "Tokens",
		CreatedAt:       time.Now(),
	}
	if err := tx.Save(&match).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to save match", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	// Trả về kết quả
	response := map[string]interface{}{
		"message":    "PvP match completed",
		"winner_id":  winner.PlayerID,
		"loser_id":   loser.PlayerID,
		"reward":     reward,
		"player1_id": player1.PlayerID,
		"player2_id": player2.PlayerID,
		"player1_dragon": map[string]interface{}{
			"id":    dragon1.ID,
			"level": dragon1.Level,
		},
		"player2_dragon": map[string]interface{}{
			"id":    dragon2.ID,
			"level": dragon2.Level,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
