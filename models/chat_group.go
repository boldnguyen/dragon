// models/chat_group.go
package models

type ChatGroup struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `json:"name"` // Tên nhóm chat (World chat, Clan chat, etc.)
	Type      string `json:"type"` // Loại nhóm (world, clan, etc.)
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
