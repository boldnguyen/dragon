package models

type Dragon struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	PlayerID   int    `json:"player_id"`
	Name       string `json:"name"`
	Rarity     string `json:"rarity"`
	Level      int    `gorm:"default:1" json:"level"`
	Attack     int    `gorm:"default:10" json:"attack"`
	Defense    int    `gorm:"default:5" json:"defense"`
	Experience int    `gorm:"default:0" json:"experience"` // Kinh nghiệm của rồng
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}
