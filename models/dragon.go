package models

type Dragon struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PlayerID  int    `json:"player_id"`
	Name      string `json:"name"`
	Rarity    string `json:"rarity"`
	Level     int    `gorm:"default:1" json:"level"`   // Cấp độ của rồng
	Attack    int    `gorm:"default:10" json:"attack"` // Mạnh mẽ (tấn công) của rồng
	Defense   int    `gorm:"default:5" json:"defense"` // Phòng thủ của rồng
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
