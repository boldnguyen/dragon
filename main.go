package main

import (
	"dragon/config"
	"dragon/models"
	"dragon/routers"
	"log"
	"net/http"
)

func main() {
	// Kết nối database
	db := config.ConnectDB()

	// Tự động tạo bảng
	models.AutoMigrate(db)

	// Định nghĩa router
	router := routers.SetupRouter(db)

	// Chạy server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
