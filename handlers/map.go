package handlers

import (
	"dragon/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Lấy danh sách các Stage (khu vực)
func GetStages(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var stages []models.Stage

	// Lấy danh sách các Stage từ database
	if err := db.Preload("Rounds.Missions").Find(&stages).Error; err != nil {
		http.Error(w, "Unable to fetch stages", http.StatusInternalServerError)
		return
	}

	// Định dạng dữ liệu trả về với thụt lề
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Sử dụng MarshalIndent để thụt lề dữ liệu JSON
	indentData, err := json.MarshalIndent(stages, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	// Gửi dữ liệu đã định dạng ra client
	w.Write(indentData)
}

// GetStageByID trả về thông tin một khu vực cụ thể
func GetStageByID(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	stageID := mux.Vars(r)["id"]
	var stage models.Stage
	if err := db.Preload("Rounds.Missions").First(&stage, stageID).Error; err != nil {
		http.Error(w, "Stage not found", http.StatusNotFound)
		return
	}

	// Trả về kết quả với định dạng JSON đẹp
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Sử dụng MarshalIndent để định dạng kết quả JSON
	indentData, err := json.MarshalIndent(stage, "", "  ")
	if err != nil {
		http.Error(w, "Error formatting JSON", http.StatusInternalServerError)
		return
	}

	// Gửi kết quả JSON đã thụt lề ra client
	w.Write(indentData)
}
