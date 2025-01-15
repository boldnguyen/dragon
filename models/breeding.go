package models

type Breeding struct {
	ID        uint  `gorm:"primaryKey" json:"id"`
	PlayerID  int   `json:"player_id"`  // ID của người chơi
	Dragon1ID uint  `json:"dragon1_id"` // ID của rồng thứ nhất
	Dragon2ID uint  `json:"dragon2_id"` // ID của rồng thứ hai
	EggID     uint  `json:"egg_id"`     // ID của trứng được tạo ra
	StartTime int64 `json:"start_time"` // Thời gian bắt đầu lai tạo
	EndTime   int64 `json:"end_time"`   // Thời gian kết thúc lai tạo
	UseToken  bool  `json:"use_token"`  // Có sử dụng token để tăng tốc không
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}
