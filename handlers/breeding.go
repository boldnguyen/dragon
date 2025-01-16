package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type StartBreedingRequest struct {
	PlayerID  int  `json:"player_id"`  // ID của người chơi
	Dragon1ID uint `json:"dragon1_id"` // ID của rồng thứ nhất
	Dragon2ID uint `json:"dragon2_id"` // ID của rồng thứ hai
	UseToken  bool `json:"use_token"`  // Có sử dụng token để tăng tốc không
}

// StartBreeding xử lý việc bắt đầu lai tạo
func StartBreeding(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req StartBreedingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Kiểm tra xem người chơi có sở hữu cả hai rồng không
	var dragon1 models.Dragon
	if err := db.Where("id = ? AND player_id = ?", req.Dragon1ID, req.PlayerID).First(&dragon1).Error; err != nil {
		http.Error(w, "Dragon 1 not found or does not belong to the player", http.StatusNotFound)
		return
	}

	var dragon2 models.Dragon
	if err := db.Where("id = ? AND player_id = ?", req.Dragon2ID, req.PlayerID).First(&dragon2).Error; err != nil {
		http.Error(w, "Dragon 2 not found or does not belong to the player", http.StatusNotFound)
		return
	}

	// Kiểm tra xem hai rồng có thể lai tạo được không (ví dụ: khác giới tính, cùng loại, v.v.)
	// Ở đây giả sử không có ràng buộc nào, chỉ cần hai rồng thuộc về cùng một người chơi

	// Tính toán thời gian lai tạo
	breedingTime := int64(600) // Mặc định 10 phút (600 giây)
	if req.UseToken {
		// Nếu sử dụng token, giảm thời gian lai tạo còn 1 phút
		breedingTime = 60
	}

	// Tạo một trứng mới từ hai rồng
	eggName := "Egg from " + dragon1.Name + " and " + dragon2.Name
	egg := models.Egg{
		Name:      eggName,
		Rarity:    calculateRarity(dragon1.Rarity, dragon2.Rarity), // Tính độ hiếm của trứng
		Price:     0,                                               // Trứng lai tạo không có giá
		Currency:  "GOLD",
		HatchTime: 300, // Thời gian ấp trứng mặc định
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Lưu trứng vào cơ sở dữ liệu
	if err := db.Create(&egg).Error; err != nil {
		http.Error(w, "Failed to create egg", http.StatusInternalServerError)
		return
	}

	// Tạo bản ghi lai tạo
	breeding := models.Breeding{
		PlayerID:  req.PlayerID,
		Dragon1ID: req.Dragon1ID,
		Dragon2ID: req.Dragon2ID,
		EggID:     egg.ID,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + breedingTime,
		UseToken:  req.UseToken,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Lưu bản ghi lai tạo vào cơ sở dữ liệu
	if err := db.Create(&breeding).Error; err != nil {
		http.Error(w, "Failed to start breeding", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message":     "Breeding started",
		"breeding_id": breeding.ID,
		"end_time":    breeding.EndTime,
		"egg_id":      egg.ID,
		"egg_name":    egg.Name,
		"egg_rarity":  egg.Rarity,
	}

	json.NewEncoder(w).Encode(response)
}

type CompleteBreedingRequest struct {
	BreedingID uint `json:"breeding_id"` // ID của bản ghi lai tạo
}

// CompleteBreeding xử lý việc hoàn thành lai tạo
func CompleteBreeding(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var req CompleteBreedingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Lấy thông tin lai tạo
	var breeding models.Breeding
	if err := db.First(&breeding, req.BreedingID).Error; err != nil {
		http.Error(w, "Breeding not found", http.StatusNotFound)
		return
	}

	// Kiểm tra xem lai tạo đã hoàn thành chưa
	if time.Now().Unix() < breeding.EndTime {
		http.Error(w, "Breeding is not yet complete", http.StatusBadRequest)
		return
	}

	// Lấy thông tin trứng
	var egg models.Egg
	if err := db.First(&egg, breeding.EggID).Error; err != nil {
		http.Error(w, "Egg not found", http.StatusNotFound)
		return
	}

	// Tạo một con rồng mới từ trứng
	dragonName := generateRandomDragonName()
	attack, defense := generateRandomAttributes(egg.Rarity)
	rarity := egg.Rarity

	newDragon := models.Dragon{
		PlayerID:  breeding.PlayerID,
		Name:      dragonName,
		Rarity:    rarity,
		Level:     1,
		Attack:    attack,
		Defense:   defense,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Thêm con rồng vào bảng Dragon
	if err := db.Create(&newDragon).Error; err != nil {
		http.Error(w, "Failed to create new dragon", http.StatusInternalServerError)
		return
	}

	// Xóa bản ghi lai tạo
	if err := db.Delete(&breeding).Error; err != nil {
		http.Error(w, "Failed to complete breeding", http.StatusInternalServerError)
		return
	}

	// Trả về phản hồi thành công
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message":       "Breeding completed",
		"dragon_id":     newDragon.ID,
		"dragon_name":   newDragon.Name,
		"dragon_rarity": newDragon.Rarity,
		"new_attack":    newDragon.Attack,
		"new_defense":   newDragon.Defense,
	}

	json.NewEncoder(w).Encode(response)
}

// Hàm tính độ hiếm của trứng dựa trên độ hiếm của hai rồng
func calculateRarity(rarity1, rarity2 string) string {
	// Logic tính độ hiếm (ví dụ: nếu cả hai rồng đều hiếm, trứng sẽ rất hiếm)
	if rarity1 == "Rare" && rarity2 == "Rare" {
		return "Very Rare"
	}
	return "Common"
}
