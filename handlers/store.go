// handlers/store.go
package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type BuyItemRequest struct {
	PlayerID int  `json:"player_id"`
	ItemID   uint `json:"item_id"`
}

func GetItems(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var items []models.Item
	query := db
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Tìm kiếm các item trong cơ sở dữ liệu
	if err := query.Find(&items).Error; err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}

	// Đảm bảo trả về dữ liệu ở định dạng JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Sử dụng MarshalIndent để tạo ra JSON với thụt lề
	indentData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	// Gửi kết quả JSON đã thụt lề ra client
	w.Write(indentData)
}

// BuyItem xử lý việc mua đồ của người chơi
func BuyItem(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req BuyItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lấy thông tin ví của người chơi
	var wallet models.Wallet
	if err := db.Where("user_id = ?", req.PlayerID).First(&wallet).Error; err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Lấy thông tin món đồ từ bảng Items
	var item models.Item
	if err := db.Where("id = ?", req.ItemID).First(&item).Error; err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	// Kiểm tra số dư trong ví trước khi mua
	if wallet.Balance < item.Price {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	// Trừ tiền từ ví người chơi
	wallet.Balance -= item.Price
	if err := db.Save(&wallet).Error; err != nil {
		http.Error(w, "Failed to update wallet", http.StatusInternalServerError)
		return
	}

	// Thêm item vào mục items của ví
	wallet.Items = append(wallet.Items, item.Name)
	if err := db.Save(&wallet).Error; err != nil {
		http.Error(w, "Failed to update wallet items", http.StatusInternalServerError)
		return
	}

	// Trả về thông báo thành công với dữ liệu đẹp hơn
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Item purchased successfully",
		"balance": wallet.Balance,
		"items":   wallet.Items,
	}

	// Sử dụng MarshalIndent để định dạng kết quả JSON
	indentData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	// Gửi kết quả JSON đã thụt lề ra client
	w.Write(indentData)
}

// BuyEggRequest là yêu cầu để mua trứng
type BuyEggRequest struct {
	PlayerID int  `json:"player_id"`
	EggID    uint `json:"egg_id"`
}

// BuyEgg xử lý việc mua trứng của người chơi
func BuyEgg(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req BuyEggRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch the egg
	var egg models.Egg
	if err := db.First(&egg, req.EggID).Error; err != nil {
		http.Error(w, "Egg not found", http.StatusNotFound)
		return
	}

	// Verify player's balance
	var profile models.Profile
	if err := db.Where("player_id = ?", req.PlayerID).First(&profile).Error; err != nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	if egg.Currency == "GOLD" && profile.TotalTokens < egg.Price {
		http.Error(w, "Insufficient GOLD", http.StatusBadRequest)
		return
	} else if egg.Currency == "DIAMOND" {
		// Placeholder for DIAMOND currency logic
	}

	// Deduct the price from the player's balance
	profile.TotalTokens -= egg.Price
	if err := db.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update player's balance", http.StatusInternalServerError)
		return
	}

	// Record the transaction
	transaction := models.Transaction{
		PlayerID:  req.PlayerID,
		ItemID:    egg.ID,
		Price:     egg.Price,
		Currency:  egg.Currency,
		CreatedAt: time.Now().Unix(),
	}

	if err := db.Create(&transaction).Error; err != nil {
		http.Error(w, "Failed to record transaction", http.StatusInternalServerError)
		return
	}

	// Create the egg in the player's wallet (or inventory)
	var wallet models.Wallet
	if err := db.Where("user_id = ?", req.PlayerID).First(&wallet).Error; err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	wallet.Items = append(wallet.Items, egg.Name) // Add the egg to the wallet's items
	if err := db.Save(&wallet).Error; err != nil {
		http.Error(w, "Failed to update wallet", http.StatusInternalServerError)
		return
	}

	// Respond back with a beautiful formatted JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Egg purchased successfully",
	}

	// Sử dụng MarshalIndent để định dạng kết quả JSON
	indentData, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	// Gửi kết quả JSON đã thụt lề ra client
	w.Write(indentData)
}
