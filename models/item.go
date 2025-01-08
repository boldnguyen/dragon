// models/item.go
package models

type Item struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"` // e.g., featured, eggs, equipment, etc.
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"` // GOLD or DIAMOND
	CreatedAt   int64   `json:"created_at"`
	UpdatedAt   int64   `json:"updated_at"`
}
