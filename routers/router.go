package routers

import (
	"dragon/handlers"
	"net/http"

	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	router := mux.NewRouter()

	// Endpoint kết nối ví
	router.HandleFunc("/api/wallet/connect", func(w http.ResponseWriter, r *http.Request) {
		handlers.ConnectWallet(db, w, r)
	}).Methods("POST")

	// Endpoint đồng bộ tài sản
	router.HandleFunc("/api/wallet/sync", func(w http.ResponseWriter, r *http.Request) {
		handlers.SyncWallet(db, w, r)
	}).Methods("GET")

	// Endpoint tạo thông tin cá nhân
	router.HandleFunc("/api/profile/create", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateProfile(db, w, r)
	}).Methods("POST")

	// Endpoint lấy thông tin profile
	router.HandleFunc("/api/profile/get", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetProfile(db, w, r)
	}).Methods("GET")

	return router
}
