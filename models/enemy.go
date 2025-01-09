package models

type Enemy struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	Level     int    `json:"level"`
	Health    int    `json:"health"`
	Attack    int    `json:"attack"`
	Defense   int    `json:"defense"`
	Reward    string `json:"reward"` // Reward could be "tokens, items, experience"
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
