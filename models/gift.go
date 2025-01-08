// models/gift.go
package models

type Gift struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PlayerID  int    `gorm:"index" json:"player_id"` // Người gửi quà
	FriendID  int    `gorm:"index" json:"friend_id"` // Người nhận quà
	Gift      string `json:"gift"`                   // Tên món quà
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
