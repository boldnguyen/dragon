package models

type MarketplaceListing struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	SellerID  int     `gorm:"index" json:"seller_id"` // PlayerID của người bán
	ItemID    uint    `gorm:"index" json:"item_id"`   // ID của vật phẩm hoặc rồng
	Price     float64 `json:"price"`                  // Giá bán
	Currency  string  `json:"currency"`               // GOLD hoặc DIAMOND
	CreatedAt int64   `json:"created_at"`
}
