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

	return router
}
