package jobs

import (
	"log"
	"tokokecil/model"

	"gorm.io/gorm"
)

// ResetStockJob: mengatur stock semua produk ke 100 setiap hari
// func ResetStockJob(db *gorm.DB) {
// 	result := db.Model(&model.Product{}).Update("stock", 100)
// 	if result.Error != nil {
// 		log.Println("[CronJob] Failed to reset stock:", result.Error)
// 	} else {
// 		log.Println("[CronJob] Stock reset to 100 for all products")
// 	}
// }

// ResetStockJob: mengatur stock semua produk ke 100 setiap interval cron
func ResetStockJob(db *gorm.DB) {
	// Tambahkan .Session agar update global diperbolehkan
	result := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Model(&model.Product{}).Update("stock", 100)
	if result.Error != nil {
		log.Println("[CronJob] Failed to reset stock:", result.Error)
	} else {
		log.Println("[CronJob] Stock reset to 100 for all products")
	}
}
