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

	// Thêm bạn bè
	router.HandleFunc("/api/friend/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddFriend(db, w, r)
	}).Methods("POST")

	// Endpoint chấp nhận bạn bè
	router.HandleFunc("/api/friend/accept", func(w http.ResponseWriter, r *http.Request) {
		handlers.AcceptFriend(db, w, r)
	}).Methods("POST")

	// Xóa bạn bè
	router.HandleFunc("/api/friend/remove", func(w http.ResponseWriter, r *http.Request) {
		handlers.RemoveFriend(db, w, r)
	}).Methods("DELETE")

	// Lấy danh sách bạn bè
	router.HandleFunc("/api/friend/list", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetFriends(db, w, r)
	}).Methods("GET")

	// Gửi quà
	router.HandleFunc("/api/friend/gift", func(w http.ResponseWriter, r *http.Request) {
		handlers.SendGift(db, w, r)
	}).Methods("POST")

	return router
}
