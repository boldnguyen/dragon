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
	// Endpoint nạp tiền vào ví
	router.HandleFunc("/api/wallet/deposit", func(w http.ResponseWriter, r *http.Request) {
		handlers.DepositFunds(db, w, r)
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

	// Endpoint gửi tin nhắn
	router.HandleFunc("/api/message/send", func(w http.ResponseWriter, r *http.Request) {
		handlers.SendMessage(db, w, r)
	}).Methods("POST")

	// Endpoint lấy tin nhắn
	router.HandleFunc("/api/message/get", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMessages(db, w, r)
	}).Methods("GET")

	// Thêm nhóm chat
	router.HandleFunc("/api/chat/group/create", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateChatGroup(db, w, r)
	}).Methods("POST")

	// Thêm thành viên vào nhóm
	router.HandleFunc("/api/chat/group/member/add", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddMemberToChatGroup(db, w, r)
	}).Methods("POST")

	// Lấy danh sách thành viên nhóm
	router.HandleFunc("/api/chat/group/member/list", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetChatGroupMembers(db, w, r)
	}).Methods("GET")

	// Xóa thành viên khỏi nhóm
	router.HandleFunc("/api/chat/group/member/remove", func(w http.ResponseWriter, r *http.Request) {
		handlers.RemoveMemberFromChatGroup(db, w, r)
	}).Methods("DELETE")

	// Endpoint lấy danh sách vật phẩm
	router.HandleFunc("/api/store/items", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetItems(db, w, r)
	}).Methods("GET")

	// Endpoint mua vật phẩm
	router.HandleFunc("/api/store/buy", func(w http.ResponseWriter, r *http.Request) {
		handlers.BuyItem(db, w, r)
	}).Methods("POST")

	// Enpoint mua trứng
	router.HandleFunc("/api/egg/buy", func(w http.ResponseWriter, r *http.Request) {
		handlers.BuyEgg(db, w, r)
	}).Methods("POST")

	// Endpoint ấp trứng
	router.HandleFunc("/api/egg/hatch", func(w http.ResponseWriter, r *http.Request) {
		handlers.HatchEgg(db, w, r)
	}).Methods("POST")

	// Endpoint hoàn thành ấp trứngtrứng
	router.HandleFunc("/api/egg/complete", func(w http.ResponseWriter, r *http.Request) {
		handlers.CompleteHatching(db, w, r)
	}).Methods("GET")

	// Các route cho bản đồ
	router.HandleFunc("/api/map/stages", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetStages(db, w, r)
	}).Methods("GET")

	router.HandleFunc("/api/map/stages/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetStageByID(db, w, r)
	}).Methods("GET")

	// Các route cho Marketplace
	router.HandleFunc("/api/marketplace/listings", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateMarketplaceListing(db, w, r) // Đăng bán vật phẩm
	}).Methods("POST")

	router.HandleFunc("/api/marketplace/listings", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetMarketplaceListings(db, w, r) // Xem danh sách vật phẩm rao bán
	}).Methods("GET")

	router.HandleFunc("/api/marketplace/listings/{id}/purchase", func(w http.ResponseWriter, r *http.Request) {
		handlers.PurchaseItem(db, w, r) // Mua vật phẩm
	}).Methods("POST")

	// Endpoint chiến đấu với đối thủ AI
	router.HandleFunc("/api/pve/fight", func(w http.ResponseWriter, r *http.Request) {
		handlers.FightEnemy(db, w, r)
	}).Methods("POST")

	return router
}
