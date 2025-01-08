// models/chat_group_member.go
package models

type ChatGroupMember struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	GroupID   uint  `gorm:"index" json:"group_id"`  // Tham chiếu đến chat group
	PlayerID  int   `gorm:"index" json:"player_id"` // Tham chiếu đến người chơi
	JoinedAt  int64 `json:"joined_at"`
	UpdatedAt int64 `json:"updated_at"`
}
