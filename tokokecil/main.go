package main

import (
	"os"

	"tokokecil/config"
	"tokokecil/handler"
	"tokokecil/jobs"
	"tokokecil/middleware"
	"tokokecil/model"
	"tokokecil/repository"
	"tokokecil/service"

	"github.com/robfig/cron/v3"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()

	// --- AutoMigrate ---
	if err := db.AutoMigrate(&model.Product{}); err != nil {
		panic("Gagal auto-migrate tabel: " + err.Error())
	}

	// --- START CRON JOB ---

	c := cron.New()
	// Setiap hari jam 00:00, jalankan ResetStockJob
	c.AddFunc("0 0 * * *", func() {
		jobs.ResetStockJob(db)
	})

	// Ekspresi cron: "* * * * *" → setiap menit
	// c.AddFunc("* * * * *", func() {
	// 	jobs.ResetStockJob(db)
	// })

	// Ekspresi cron: "*/2 * * * *" → setiap 2 menit
	// c.AddFunc("*/2 * * * *", func() {
	// 	jobs.ResetStockJob(db)
	// })

	c.Start()

	// --- END CRON JOB ---

	// Dependency injection Product saja
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	e := echo.New()

	// Pasang Logrus custom logger middleware (jika masih dipakai)
	e.Use(middleware.MiddlewareLogging)

	// Integrasi custom error handler
	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	// Product routes (protected by JWT middleware)
	products := e.Group("/products")
	products.Use(middleware.JWTAuth) // Apply JWT middleware!
	products.GET("", productHandler.GetAllProducts)
	products.POST("", productHandler.CreateProduct)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

// package main

// import (
// 	"log"
// 	// "os"

// 	"tokokecil/config"
// 	grpcdelivery "tokokecil/delivery/grpc"

// 	// "tokokecil/handler"
// 	// "tokokecil/jobs"
// 	// "tokokecil/middleware"
// 	"tokokecil/model"
// 	"tokokecil/repository"
// 	"tokokecil/service"
// 	// "github.com/robfig/cron/v3"
// 	// _ "tokokecil/docs"
// 	// "github.com/labstack/echo/v4"
// 	// echoSwagger "github.com/swaggo/echo-swagger"
// )

// func main() {
// 	config.LoadEnv()
// 	db := config.DBInit()

// 	// --- AutoMigrate ---
// 	if err := db.AutoMigrate(&model.Product{}); err != nil {
// 		panic("Gagal auto-migrate tabel: " + err.Error())
// 	}

// 	// // --- START CRON JOB ---
// 	// c := cron.New()
// 	// c.AddFunc("0 0 * * *", func() {
// 	//     jobs.ResetStockJob(db)
// 	// })
// 	// c.Start()
// 	// // --- END CRON JOB ---

// 	// Dependency injection Product saja
// 	productRepo := repository.NewProductRepository(db)
// 	productService := service.NewProductService(productRepo)
// 	// productHandler := handler.NewProductHandler(productService)

// 	// // Echo HTTP server & Swagger (Dikomentari, nonaktif)
// 	// e := echo.New()
// 	// e.Use(middleware.MiddlewareLogging)
// 	// e.HTTPErrorHandler = handler.CustomHTTPErrorHandler
// 	// e.GET("/swagger/*", echoSwagger.WrapHandler)
// 	// products := e.Group("/products")
// 	// products.Use(middleware.JWTAuth)
// 	// products.GET("", productHandler.GetAllProducts)
// 	// products.POST("", productHandler.CreateProduct)
// 	// port := os.Getenv("PORT")
// 	// if port == "" {
// 	//     port = "8080"
// 	// }
// 	// e.Logger.Fatal(e.Start(":" + port))

// 	// Jalankan gRPC server (as main service)
// 	log.Println("Starting gRPC Product Service on :50051 ...")
// 	grpcdelivery.StartGRPCServer(productService)
// }
