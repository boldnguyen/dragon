package models

type Clan struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `json:"name"`        // Tên clan
	LeaderID    int    `json:"leader_id"`   // ID của người chơi là leader
	Description string `json:"description"` // Mô tả clan
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type ClanMember struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ClanID    uint   `json:"clan_id"`   // ID của clan
	PlayerID  int    `json:"player_id"` // ID của người chơi
	Role      string `json:"role"`      // Vai trò trong clan (leader, member, v.v.)
	JoinedAt  int64  `json:"joined_at"` // Thời gian tham gia clan
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type ClanTask struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ClanID      uint   `json:"clan_id"`     // ID của clan
	TaskName    string `json:"task_name"`   // Tên nhiệm vụ
	Description string `json:"description"` // Mô tả nhiệm vụ
	Reward      string `json:"reward"`      // Phần thưởng
	Completed   bool   `json:"completed"`   // Trạng thái hoàn thành
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type ClanChat struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ClanID    uint   `json:"clan_id"`   // ID của clan
	PlayerID  int    `json:"player_id"` // ID của người chơi
	Message   string `json:"message"`   // Nội dung tin nhắn
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
