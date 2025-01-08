// models/friendship.go
package models

type Friendship struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	PlayerID  int    `gorm:"index" json:"player_id"` // Tham chiếu đến PlayerID trong bảng Profile
	FriendID  int    `gorm:"index" json:"friend_id"` // Tham chiếu đến PlayerID của người bạn
	Status    string `json:"status"`                 // "pending", "accepted", "blocked", v.v.
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
