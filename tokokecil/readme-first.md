# Tokokecil: Migrasi Go + Echo dari MySQL ke PostgreSQL + GORM

Panduan ini menjelaskan langkah-langkah migrasi project Tokokecil (Go + Echo) dari MySQL ke PostgreSQL menggunakan GORM, mulai dari setup database, environment, hingga running aplikasi tanpa error AutoMigrate.

## Pendekatan SOLID & Scope Perubahan

Project Tokokecil didesain dengan prinsip **SOLID**:
- **Single Responsibility:** Setiap layer punya tugas spesifik.
- **Open/Closed Principle:** Mudah di-extend tanpa mengubah handler.
- **Lainnya:** Dependensi jelas, perubahan hanya perlu di lapisan _model_, _repository_, _service_, dan _main.go_.

**Hasilnya:**  
Saat migrasi database, **hanya 4 bagian ini yang perlu diubah**:
- Model (`model/`)
- Repository (`repository/`)
- Service (`service/`)
- Main entrypoint (`main.go`)

> **Catatan:**  
> Layer **handler** tidak perlu diubah karena tetap berkomunikasi lewat interface Service yang tidak berubah secara kontrak, sehingga controller/api layer tetap stabil dan tidak terganggu perubahan DB.
---

## 1. Siapkan Database PostgreSQL

- Pastikan PostgreSQL sudah ter-install dan berjalan di local.
- Buat database baru (misal: `tokoku_abcde`) agar tidak bentrok dengan database lama:
  ```sql
  CREATE DATABASE tokoku_abcde;
  ```
- Pastikan user (`postgres`) punya akses penuh ke database ini.  
  Uji koneksi via terminal:
  ```sh
  psql -U postgres -d tokoku_abcde
  ```
---

## 2. Update Environment Variable `.env`

- Edit file `.env` di root project:
  ```env
  DB_USER=postgres
  DB_HOST=127.0.0.1
  DB_PORT=5432
  DB_NAME=tokoku_abcde
  JWT_SECRET=mysecretkey123
  ```
- Jika tidak pakai password, field `DB_PASS` boleh dihapus/dikosongkan.

---

## 3. Update Koneksi di `config.go` (GORM)

- Ubah kode koneksi ke PostgreSQL, tanpa password jika tidak digunakan:
  ```go
  user := os.Getenv("DB_USER")
  host := os.Getenv("DB_HOST")
  port := os.Getenv("DB_PORT")
  name := os.Getenv("DB_NAME")

  dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", host, user, name, port)
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
  if err != nil {
      log.Fatal("Gagal koneksi ke database:", err)
  }
  ```
- Tambahkan log print nama database untuk memastikan:
  ```go
  var dbName string
  db.Raw("SELECT current_database()").Scan(&dbName)
  fmt.Println("PostgreSQL benar-benar connect ke:", dbName)
  ```

---

## 4. Pastikan Database Benar-benar Kosong

- Cek database sebelum menjalankan migrasi:
  ```sql
  SELECT table_schema, table_name
  FROM information_schema.tables
  WHERE table_catalog = 'tokoku_abcde';
  ```
- Pastikan tidak ada tabel `users`, `products`, atau tabel sisa lain di schema apapun.

---

## 5. Update Model: GORM Struct

File: `model/user.go`
```go
package model

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"type:varchar(100);not null"`
    Email     string    `gorm:"type:varchar(100);unique;not null"`
    Password  string    `gorm:"type:text;not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

File: `model/product.go`
```go
package model

import "time"

type Product struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"type:varchar(100);not null"`
    Price     float64   `gorm:"type:numeric(12,2);not null"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
---
## 6. Update Repository: GORM CRUD Implementation

File: `repository/user_repository.go`
```go
package repository

import (
    "tokokecil/model"
    "gorm.io/gorm"
)

type UserRepository interface {
    GetAll() ([]model.User, error)
    GetByID(id uint) (model.User, error)
    GetByEmail(email string) (model.User, error)
    Create(user model.User) (model.User, error)
    Update(id uint, user model.User) (model.User, error)
    Delete(id uint) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db}
}

func (r *userRepository) GetAll() ([]model.User, error) {
    var users []model.User
    err := r.db.Find(&users).Error
    return users, err
}

func (r *userRepository) GetByID(id uint) (model.User, error) {
    var user model.User
    err := r.db.First(&user, id).Error
    return user, err
}

func (r *userRepository) GetByEmail(email string) (model.User, error) {
    var user model.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return user, err
}

func (r *userRepository) Create(user model.User) (model.User, error) {
    err := r.db.Create(&user).Error
    return user, err
}

func (r *userRepository) Update(id uint, user model.User) (model.User, error) {
    var u model.User
    err := r.db.First(&u, id).Error
    if err != nil {
        return u, err
    }
    u.Name = user.Name
    u.Email = user.Email
    // Untuk update password, tambahkan jika perlu
    err = r.db.Save(&u).Error
    return u, err
}

func (r *userRepository) Delete(id uint) error {
    return r.db.Delete(&model.User{}, id).Error
}
```

File: `repository/product_repository.go`
```go
package repository

import (
    "tokokecil/model"
    "gorm.io/gorm"
)

type ProductRepository interface {
    GetAll() ([]model.Product, error)
    Create(product model.Product) (model.Product, error)
}

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
    return &productRepository{db}
}

func (r *productRepository) GetAll() ([]model.Product, error) {
    var products []model.Product
    err := r.db.Find(&products).Error
    return products, err
}

func (r *productRepository) Create(product model.Product) (model.Product, error) {
    err := r.db.Create(&product).Error
    return product, err
}
```

---

## 7. Update Service Layer: No Logic Change, Only Type (uint) Fix

File: `service/user_service.go`
```go
package service

import (
    "errors"
    "tokokecil/model"
    "tokokecil/repository"
)

type UserService interface {
    GetAllUsers() ([]model.User, error)
    GetUserByID(id int) (model.User, error)
    GetByEmail(email string) (model.User, error)
    CreateUser(user model.User) (model.User, error)
    UpdateUser(id int, user model.User) (model.User, error)
    DeleteUser(id int) error
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(r repository.UserRepository) UserService {
    return &userService{repo: r}
}

func (s *userService) GetAllUsers() ([]model.User, error) {
    return s.repo.GetAll()
}

func (s *userService) GetUserByID(id int) (model.User, error) {
    return s.repo.GetByID(uint(id))
}

func (s *userService) GetByEmail(email string) (model.User, error) {
    return s.repo.GetByEmail(email)
}

func (s *userService) CreateUser(user model.User) (model.User, error) {
    if user.Name == "" || user.Email == "" || user.Password == "" {
        return model.User{}, errors.New("name, email, password wajib diisi")
    }
    return s.repo.Create(user)
}

func (s *userService) UpdateUser(id int, user model.User) (model.User, error) {
    if user.Name == "" || user.Email == "" {
        return model.User{}, errors.New("name & email wajib diisi")
    }
    return s.repo.Update(uint(id), user)
}

func (s *userService) DeleteUser(id int) error {
    return s.repo.Delete(uint(id))
}
```

File: `service/product_service.go`
```go
package service

import (
    "errors"
    "tokokecil/model"
    "tokokecil/repository"
)

type ProductService interface {
    GetAllProducts() ([]model.Product, error)
    CreateProduct(product model.Product) (model.Product, error)
}

type productService struct {
    repo repository.ProductRepository
}

func NewProductService(r repository.ProductRepository) ProductService {
    return &productService{repo: r}
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
    return s.repo.GetAll()
}

func (s *productService) CreateProduct(product model.Product) (model.Product, error) {
    if product.Name == "" || product.Price <= 0 {
        return model.Product{}, errors.New("name & price wajib diisi")
    }
    return s.repo.Create(product)
}
```
---
## 8. Update main.go: Dependency Injection & AutoMigrate

File: `main.go`
```go
package main

import (
    "os"
    "fmt"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "tokokecil/config"
    "tokokecil/handler"
    custommw "tokokecil/middleware"
    "tokokecil/repository"
    "tokokecil/service"
    "tokokecil/model"
)

func main() {
    config.LoadEnv()
    db := config.DBInit()

    // Print DB name (debug)
    var dbName string
    db.Raw("SELECT current_database()").Scan(&dbName)
    fmt.Println("PostgreSQL benar-benar connect ke:", dbName)

    // AutoMigrate semua model
    if err := db.AutoMigrate(&model.User{}, &model.Product{}); err != nil {
        panic("Gagal auto-migrate tabel: " + err.Error())
    }

    // Dependency injection
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    userHandler := handler.NewUserHandler(userService)
    authHandler := handler.NewAuthHandler(userService)

    productRepo := repository.NewProductRepository(db)
    productService := service.NewProductService(productRepo)
    productHandler := handler.NewProductHandler(productService)

    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())

    // Auth routes
    e.POST("/register", authHandler.Register)
    e.POST("/login", authHandler.Login)

    // User routes
    e.GET("/users", userHandler.GetAllUsers)

    // Product routes (JWT protected)
    products := e.Group("/products")
    products.Use(custommw.JWTAuth)
    products.GET("", productHandler.GetAllProducts)
    products.POST("", productHandler.CreateProduct)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    e.Logger.Fatal(e.Start(":" + port))
}
```

---
## 9. Jalankan Project & AutoMigrate

- Jalankan aplikasi:
  ```sh
  go run main.go
  ```
- Ekspektasi: Akan muncul log
  ```
  PostgreSQL benar-benar connect ke: tokoku_abcde
  ```
  dan tabel akan otomatis dibuat di database tersebut tanpa error.

---

## 10. Troubleshooting Penting

- Jika error "column ... contains null values":
  - Kemungkinan besar masih connect ke database/postgres/schema yang salah.
  - Pastikan `.env` sudah benar, database benar-benar kosong, dan DSN tidak ada typo.
  - Cek log DSN dan nama database yang terkoneksi di output aplikasi.
- Jika tidak ada password, hapus/abaikan field password di DSN dan `.env`.

---

### Selesai!

Kamu sudah berhasil migrasi ke PostgreSQL + GORM dan project siap dikembangkan lebih lanjut.

---

**Tips lanjutan:**  
- Selalu jalankan query pengecekan tabel sebelum migrasi besar.
- Tambahkan test koneksi di awal startup untuk validasi environment.
