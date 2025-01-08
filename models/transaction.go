// models/transaction.go
package models

type Transaction struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	PlayerID  int     `gorm:"index" json:"player_id"`
	ItemID    uint    `gorm:"index" json:"item_id"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
	CreatedAt int64   `json:"created_at"`
}
