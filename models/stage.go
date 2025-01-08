package models

import (
	"time"
)

// Stage đại diện cho một khu vực trên bản đồ
type Stage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Rounds []Round `json:"rounds"`
}

type Round struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	StageID     uint      `json:"stage_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Missions []Mission `json:"missions"`
}

type Mission struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	RoundID     uint      `json:"round_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Reward      string    `json:"reward"` // Phần thưởng khi hoàn thành nhiệm vụ
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
