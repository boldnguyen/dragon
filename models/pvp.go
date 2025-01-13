package models

import (
	"time"
)

type PvPMatch struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Player1ID       int       `json:"player1_id"`
	Player2ID       int       `json:"player2_id"`
	Player1DragonID uint      `json:"player1_dragon_id"`
	Player2DragonID uint      `json:"player2_dragon_id"`
	BetAmount       float64   `json:"bet_amount"`
	WinnerID        int       `json:"winner_id"` // ID người chơi thắng
	Reward          string    `json:"reward"`    // Thưởng (token/vật phẩm)
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
