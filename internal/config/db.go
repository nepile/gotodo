package config

import (
	"fmt"
	"log"

	"github.com/nepile/gotodo/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Kredensial sesuai docker-compose.yml
	dsn := "host=localhost user=neville password=235314166 dbname=devcon_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database! \n", err)
	}

	fmt.Println("Koneksi Database Berhasil!")

	// Auto Migration untuk membuat tabel otomatis
	err = db.AutoMigrate(&model.User{}, &model.Todo{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi! \n", err)
	}

	DB = db
}
