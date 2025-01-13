package models

type Boss struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"`
	Health    int    `json:"health"`
	Attack    int    `json:"attack"`
	Defense   int    `json:"defense"`
	Special   string `json:"special"` // Kỹ năng đặc biệt của boss
	Level     int    `json:"level"`   // Cấp độ của boss
	Reward    string `json:"reward"`  // Phần thưởng khi đánh bại boss
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
