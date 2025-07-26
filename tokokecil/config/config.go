package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, pakai environment bawaan OS")
	}
}

func DBInit() *gorm.DB {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=require",
		host, user, pass, name, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// Connection Pooling Config
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error mengambil database object:", err)
	}

	maxOpen, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_OPEN_CONNS", "20"))
	maxIdle, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_IDLE_CONNS", "10"))
	maxLife, _ := strconv.Atoi(getEnvOrDefault("DB_CONN_MAX_LIFETIME", "60"))

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLife) * time.Minute)

	log.Printf("✅ Pooling PostgreSQL: maxOpen=%d, maxIdle=%d, maxLife=%d menit\n", maxOpen, maxIdle, maxLife)
	return db
}

// Helper
func getEnvOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
