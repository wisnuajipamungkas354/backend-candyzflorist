package main

import (
	"fmt"
	"log"
	"os"

	"candyzflorist-go/config"
	"candyzflorist-go/middleware"
	"candyzflorist-go/routes"
	"candyzflorist-go/seeders"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ File .env tidak ditemukan, menggunakan environment variables sistem")
	}

	// 2. Ensure Uploads Folder Exists
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Printf("⚠️ Gagal membuat folder uploads: %v\n", err)
	}

	// 3. Connect to Database & Auto-Migrate
	db := config.ConnectDatabase()

	// 4. Run Seeders
	if db != nil {
		seeders.SeedDatabase(db)
	}

	// 5. Initialize Gin Engine
	r := gin.Default()

	// 6. Attach CORS Middleware
	r.Use(middleware.CORSMiddleware())

	// 7. Setup Application Routes
	routes.SetupRoutes(r)

	// 8. Run Server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 CandyzFlorist Backend Server berjalan di http://localhost:%s\n", port)
	log.Printf("📁 Direktori penyimpanan upload: %s\n", uploadDir)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
