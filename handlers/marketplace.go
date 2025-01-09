package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// CreateMarketplaceListing handles creating a new marketplace listing
func CreateMarketplaceListing(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req struct {
		SellerID int     `json:"seller_id"`
		ItemID   uint    `json:"item_id"`
		Price    float64 `json:"price"`
		Currency string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	listing := models.MarketplaceListing{
		SellerID:  req.SellerID,
		ItemID:    req.ItemID,
		Price:     req.Price,
		Currency:  req.Currency,
		CreatedAt: time.Now().Unix(),
	}

	if err := db.Create(&listing).Error; err != nil {
		http.Error(w, "Failed to create listing", http.StatusInternalServerError)
		return
	}

	// Return the created listing in pretty JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response, _ := json.MarshalIndent(listing, "", "  ")
	w.Write(response)
}

// GetMarketplaceListings handles fetching all marketplace listings
func GetMarketplaceListings(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var listings []models.MarketplaceListing
	if err := db.Find(&listings).Error; err != nil {
		http.Error(w, "Failed to fetch listings", http.StatusInternalServerError)
		return
	}

	// Return the listings in pretty JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, _ := json.MarshalIndent(listings, "", "  ")
	w.Write(response)
}

// PurchaseItem xử lý việc mua vật phẩm từ marketplace
func PurchaseItem(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	listingID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid listing ID", http.StatusBadRequest)
		return
	}

	var req struct {
		BuyerID int `json:"buyer_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var listing models.MarketplaceListing
	if err := db.First(&listing, listingID).Error; err != nil {
		http.Error(w, "Listing not found", http.StatusNotFound)
		return
	}

	// Kiểm tra nếu người mua có đủ tiền trong ví
	var buyerWallet models.Wallet
	if err := db.Where("user_id = ?", req.BuyerID).First(&buyerWallet).Error; err != nil {
		http.Error(w, "Buyer wallet not found", http.StatusNotFound)
		return
	}

	if buyerWallet.Balance < listing.Price {
		http.Error(w, "Insufficient balance", http.StatusPaymentRequired)
		return
	}

	// Kiểm tra ví của người bán và cập nhật
	var sellerWallet models.Wallet
	if err := db.Where("user_id = ?", listing.SellerID).First(&sellerWallet).Error; err != nil {
		http.Error(w, "Seller wallet not found", http.StatusNotFound)
		return
	}

	// Giao dịch: trừ tiền của người mua và cộng vào người bán
	tx := db.Begin()
	buyerWallet.Balance -= listing.Price
	sellerWallet.Balance += listing.Price * 0.95 // 5% phí giao dịch

	if err := tx.Save(&buyerWallet).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update buyer wallet", http.StatusInternalServerError)
		return
	}

	if err := tx.Save(&sellerWallet).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update seller wallet", http.StatusInternalServerError)
		return
	}

	// Thêm item vào ví của người mua
	var item models.Item
	if err := db.Where("id = ?", listing.ItemID).First(&item).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	buyerWallet.Items = append(buyerWallet.Items, item.Name)

	// Cập nhật lại ví của người mua
	if err := tx.Save(&buyerWallet).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update buyer wallet items", http.StatusInternalServerError)
		return
	}

	// Xóa listing sau khi giao dịch thành công
	if err := tx.Delete(&listing).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete listing", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	// Trả về kết quả thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":       "Purchase successful",
		"buyer_balance": buyerWallet.Balance,
		"items":         buyerWallet.Items,
	}
	indentData, _ := json.MarshalIndent(response, "", "  ")
	w.Write(indentData)
}
