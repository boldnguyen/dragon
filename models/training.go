package models

type Training struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	PlayerID     int    `json:"player_id"`     // ID của người chơi
	DragonID     uint   `json:"dragon_id"`     // ID của rồng được huấn luyện
	TrainingType string `json:"training_type"` // Loại huấn luyện (ví dụ: "attack", "defense")
	StartTime    int64  `json:"start_time"`    // Thời gian bắt đầu huấn luyện
	EndTime      int64  `json:"end_time"`      // Thời gian kết thúc huấn luyện
	UseToken     bool   `json:"use_token"`     // Có sử dụng token để tăng tốc không
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}
