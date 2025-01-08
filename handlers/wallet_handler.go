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

// DepositFundsRequest là cấu trúc yêu cầu nạp tiền vào ví
type DepositFundsRequest struct {
	PlayerID int     `json:"player_id"`
	Amount   float64 `json:"amount"`
}

// DepositFunds xử lý việc nạp tiền vào ví của người chơi
func DepositFunds(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req DepositFundsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Tìm ví của người chơi
	var wallet models.Wallet
	if err := db.Where("user_id = ?", req.PlayerID).First(&wallet).Error; err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Kiểm tra xem số tiền có hợp lệ không
	if req.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	// Cập nhật số dư ví
	wallet.Balance += req.Amount

	// Lưu lại ví
	if err := db.Save(&wallet).Error; err != nil {
		http.Error(w, "Failed to deposit funds", http.StatusInternalServerError)
		return
	}

	// Cập nhật số token trong profile (nếu cần)
	var profile models.Profile
	if err := db.Where("player_id = ?", req.PlayerID).First(&profile).Error; err == nil {
		// Cập nhật số token bằng số tiền trong ví (hoặc theo tỉ lệ chuyển đổi)
		profile.TotalTokens += req.Amount // Giả sử mỗi 1 đơn vị tiền trong ví tương ứng với 1 token
		if err := db.Save(&profile).Error; err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Funds deposited successfully",
		"balance": wallet.Balance,
	})
}

// SyncWallet đồng bộ số dư ví và token trong profile
func SyncWallet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	var wallet models.Wallet
	if err := db.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Tìm thông tin profile của người chơi
	var profile models.Profile
	if err := db.Where("player_id = ?", userID).First(&profile).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Cập nhật lại số token trong profile dựa trên số dư ví
	profile.TotalTokens = wallet.Balance // Giả sử mỗi đơn vị tiền trong ví tương ứng với 1 token
	if err := db.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to sync profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Wallet and profile synced successfully",
		"balance": wallet.Balance,
		"tokens":  profile.TotalTokens,
	})
}
