package config

import (
	"fmt"
	"log"
	"os"

	"candyzflorist-go/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "3306"
	}
	if user == "" {
		user = "root"
	}
	if dbname == "" {
		dbname = "candyzflorist"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Printf("⚠️ Gagal terhubung ke database MySQL (%s): %v\n", dsn, err)
		log.Println("💡 Pastikan MySQL aktif dan database 'candyzflorist' sudah dibuat.")
		return nil
	}

	log.Println("✅ Berhasil terhubung ke database MySQL!")

	// Auto Migration
	err = DB.AutoMigrate(
		&models.User{},
		&models.Kategori{},
		&models.Produk{},
		&models.ProdukImage{},
		&models.Pengaturan{},
	)
	if err != nil {
		log.Printf("⚠️ AutoMigrate error: %v\n", err)
	} else {
		log.Println("✅ Database schema migrated successfully!")
	}

	return DB
}
