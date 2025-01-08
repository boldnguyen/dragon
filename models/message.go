package models

type Message struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	SenderID       int    `json:"sender_id"`
	ReceiverID     int    `json:"receiver_id"`
	Content        string `json:"content"`
	ChatType       string `json:"chat_type"` // Loại chat: world, clan, friend, team
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
	SenderName     string `json:"sender_name,omitempty"`
	SenderAvatar   string `json:"sender_avatar,omitempty"`
	ReceiverName   string `json:"receiver_name,omitempty"`   // Thêm trường này
	ReceiverAvatar string `json:"receiver_avatar,omitempty"` // Thêm trường này
}
