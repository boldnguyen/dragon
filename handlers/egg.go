package handlers

import (
	"dragon/models"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// HatchEggRequest là yêu cầu để ấp trứng
type HatchEggRequest struct {
	PlayerID int  `json:"player_id"`
	EggID    uint `json:"egg_id"`
	UseToken bool `json:"use_token"` // Nếu true, trứng sẽ ấp nhanh hơn và tốn token
}

// HatchEgg xử lý việc ấp trứng
func HatchEgg(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req HatchEggRequest
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

	// Verify player's wallet or items
	var wallet models.Wallet
	if err := db.Where("user_id = ?", req.PlayerID).First(&wallet).Error; err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Check if the player has the egg
	eggFound := false
	for _, item := range wallet.Items {
		if item == egg.Name {
			eggFound = true
			break
		}
	}

	if !eggFound {
		http.Error(w, "Player does not own this egg", http.StatusBadRequest)
		return
	}

	// Calculate hatch time (if using token, speed up the process)
	hatchTime := egg.HatchTime
	if req.UseToken {
		// Giả sử sử dụng token sẽ giảm thời gian ấp (ví dụ: giảm 50%)
		hatchTime /= 2
	}

	// Ghi lại thời gian ấp trứng trong cơ sở dữ liệu hoặc hệ thống
	hatchEndTime := time.Now().Unix() + hatchTime

	// Thêm thông tin về thời gian ấp vào hệ thống
	var profile models.Profile
	if err := db.Where("player_id = ?", req.PlayerID).First(&profile).Error; err != nil {
		http.Error(w, "Player profile not found", http.StatusNotFound)
		return
	}

	// Giảm số token của người chơi (nếu cần)
	if req.UseToken {
		if profile.TotalTokens < 10 {
			http.Error(w, "Not enough tokens", http.StatusBadRequest)
			return
		}
		profile.TotalTokens -= 10 // Giả sử ấp nhanh hơn tốn 10 token
		if err := db.Save(&profile).Error; err != nil {
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}
	}

	// Trả về kết quả thành công với định dạng JSON đẹp
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message":     "Egg is hatching",
		"hatch_time":  hatchEndTime,
		"use_token":   req.UseToken,
		"new_balance": profile.TotalTokens,
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

// Helper function để tạo tên rồng ngẫu nhiên
func generateRandomDragonName() string {
	names := []string{"Firestorm", "Thunderstrike", "Frostbite", "Blazewing", "Shadowflame", "Stoneheart"}
	rand.Seed(time.Now().UnixNano())
	return names[rand.Intn(len(names))]
}

// Helper function để tạo giá trị ngẫu nhiên cho sức mạnh tấn công và phòng thủ
func generateRandomAttributes(rarity string) (int, int) {
	// Random attack power between 5 and 20
	var attack int
	var defense int

	// Random attack and defense based on rarity
	switch rarity {
	case "Common":
		attack = rand.Intn(10) + 5
		defense = rand.Intn(6) + 3
	case "Uncommon":
		attack = rand.Intn(15) + 10
		defense = rand.Intn(8) + 5
	case "Rare":
		attack = rand.Intn(20) + 15
		defense = rand.Intn(10) + 7
	case "Epic":
		attack = rand.Intn(25) + 20
		defense = rand.Intn(12) + 9
	case "Legendary":
		attack = rand.Intn(30) + 25
		defense = rand.Intn(15) + 12
	}

	return attack, defense
}

// Helper function để chọn độ hiếm ngẫu nhiên
func generateRandomRarity() string {
	rareTypes := []string{"Common", "Uncommon", "Rare", "Epic", "Legendary"}
	rand.Seed(time.Now().UnixNano())
	return rareTypes[rand.Intn(len(rareTypes))]
}

func CompleteHatching(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	playerIDStr := r.URL.Query().Get("player_id")

	// Chuyển đổi playerID từ string sang int
	playerID, err := strconv.Atoi(playerIDStr)
	if err != nil {
		http.Error(w, "Invalid player ID", http.StatusBadRequest)
		return
	}

	// Tìm thông tin người chơi và trạng thái trứng
	var profile models.Profile
	if err := db.Where("player_id = ?", playerID).First(&profile).Error; err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Tạo một con rồng ngẫu nhiên
	dragonName := generateRandomDragonName()
	attack, defense := generateRandomAttributes(generateRandomRarity())
	rarity := generateRandomRarity()

	newDragon := models.Dragon{
		PlayerID:  playerID,   // Sử dụng playerID kiểu int
		Name:      dragonName, // Tạo tên ngẫu nhiên
		Rarity:    rarity,     // Tạo độ hiếm ngẫu nhiên
		Level:     1,          // Mặc định rồng mới ở cấp độ 1
		Attack:    attack,     // Sức mạnh tấn công ngẫu nhiên
		Defense:   defense,    // Sức mạnh phòng thủ ngẫu nhiên
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Thêm con rồng vào bảng Dragon
	if err := db.Create(&newDragon).Error; err != nil {
		http.Error(w, "Failed to create new dragon", http.StatusInternalServerError)
		return
	}

	// Lấy ví của người chơi và thêm rồng vào ví
	var wallet models.Wallet
	if err := db.Where("user_id = ?", playerID).First(&wallet).Error; err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	// Thêm rồng vào ví
	wallet.Dragons = append(wallet.Dragons, newDragon.Name)
	if err := db.Save(&wallet).Error; err != nil {
		http.Error(w, "Failed to update wallet", http.StatusInternalServerError)
		return
	}

	// Cập nhật lại số lượng rồng trong Profile
	profile.DragonsOwned++
	if err := db.Save(&profile).Error; err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Trả về kết quả thành công với định dạng JSON đẹp
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message": "Egg hatched, new dragon acquired!",
		"dragon":  newDragon,
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

