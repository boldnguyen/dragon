package models

import (
	"database/sql/driver"
	"encoding/json"
	"log"

	"gorm.io/gorm"
)

// StringArray là kiểu tùy chỉnh để xử lý JSONB
type StringArray []string

// Value chuyển StringArray thành JSON để lưu vào database
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan đọc JSON từ database và chuyển thành StringArray
func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return json.Unmarshal([]byte("[]"), s) // Gán giá trị mặc định nếu không đúng kiểu
	}
	return json.Unmarshal(bytes, s)
}

type Wallet struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	UserID    int         `json:"user_id"`
	PublicKey string      `json:"public_key"`
	Balance   float64     `json:"balance" gorm:"default:0"`
	Dragons   StringArray `gorm:"type:jsonb;default:'[]'" json:"dragons"`
	Items     StringArray `gorm:"type:jsonb;default:'[]'" json:"items"`
	CreatedAt int64       `json:"created_at"`
	UpdatedAt int64       `json:"updated_at"`
}

// AutoMigrate đảm bảo rằng bảng được tạo ra trong cơ sở dữ liệu
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&Wallet{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
