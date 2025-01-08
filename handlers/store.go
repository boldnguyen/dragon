// handlers/store.go
package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"

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

	if err := query.Find(&items).Error; err != nil {
		http.Error(w, "Failed to fetch items", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(items)
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

	// Trả về thông báo thành công
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Item purchased successfully",
		"balance": wallet.Balance,
		"items":   wallet.Items,
	})
}
