// models/egg.go
package models

type Egg struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Name      string  `json:"name"`
	Rarity    string  `json:"rarity"`     // Độ hiếm của trứng (common, rare, legendary, etc.)
	Price     float64 `json:"price"`      // Giá của trứng, nếu mua bằng token
	Currency  string  `json:"currency"`   // Loại tiền tệ (e.g., GOLD, DIAMOND)
	HatchTime int64   `json:"hatch_time"` // Thời gian để ấp trứng (theo giây)
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}
