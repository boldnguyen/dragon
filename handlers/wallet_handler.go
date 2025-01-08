package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type ConnectWalletRequest struct {
	UserID    int    `json:"user_id"`
	PublicKey string `json:"public_key"`
}

func ConnectWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req ConnectWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Kiểm tra xem ví đã tồn tại chưa
	var existingWallet models.Wallet
	err := db.Where("public_key = ?", req.PublicKey).First(&existingWallet).Error

	if err == nil {
		// Ví đã tồn tại, trả về thông báo
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wallet_id": existingWallet.ID,
			"message":   "Wallet already connected",
		})
		return
	} else if err != gorm.ErrRecordNotFound {
		// Có lỗi trong quá trình kiểm tra
		http.Error(w, "Failed to connect wallet", http.StatusInternalServerError)
		return
	}

	// Tạo ví mới
	wallet := models.Wallet{
		UserID:    req.UserID,
		PublicKey: req.PublicKey,
	}

	if err := db.Create(&wallet).Error; err != nil {
		http.Error(w, "Failed to connect wallet", http.StatusInternalServerError)
		return
	}

	// Trả về thành công
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"wallet_id": wallet.ID,
		"message":   "Wallet connected successfully",
	})
}

func SyncWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	var wallet models.Wallet
	if err := db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}
