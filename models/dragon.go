// models/dragon.go
package models

type Dragon struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PlayerID  int    `json:"player_id"`
	Name      string `json:"name"`
	Rarity    string `json:"rarity"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
