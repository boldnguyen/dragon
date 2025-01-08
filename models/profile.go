package models

type Profile struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	PlayerID     int     `gorm:"uniqueIndex" json:"player_id"`
	Name         string  `json:"name"`
	Level        int     `gorm:"default:1" json:"level"`
	Avatar       string  `json:"avatar"`
	Wins         int     `json:"wins"`
	Losses       int     `json:"losses"`
	DragonsOwned int     `json:"dragons_owned"`
	TotalTokens  float64 `json:"total_tokens"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
}
