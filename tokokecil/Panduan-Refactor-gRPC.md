# Panduan Implementasi gRPC CRUD pada Project Tokokecil

---

## 1. Prasyarat & Tools yang Perlu Di-install

### **A. Install Go & Protobuf**
- [Go](https://go.dev/dl/), minimal versi 1.18
- [Protoc (Protocol Buffers Compiler)](https://github.com/protocolbuffers/protobuf/releases)

### **B. Install Plugin Protobuf untuk Go**
Jalankan di terminal:
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

> Pastikan `$GOPATH/bin` sudah ada di $PATH, agar plugin dapat dijalankan global.

### **C. (Opsional) Install grpcurl untuk testing CLI**
```sh
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

### **D. (Opsional) Postman versi terbaru**
- Download di https://www.postman.com/downloads/

---

## 2. Struktur Folder Project

```plaintext
tokokecil/
├── proto/                # File .proto
├── pb/                   # Hasil generate file Go dari .proto
├── delivery/
│   └── grpc/             # Handler & server gRPC
├── model/                # Struct data produk
├── repository/           # DB logic
├── service/              # Bisnis logic
├── main.go               # Entrypoint aplikasi
└── ...                   # File/folder lain
```

---

## 3. Membuat & Mengupdate File Proto

Buat file: `proto/product.proto`

```proto
syntax = "proto3";
package product;
option go_package = "./pb;pb";

service ProductService {
  rpc GetAllProducts (Empty) returns (ProductListResponse);
  rpc CreateProduct (CreateProductRequest) returns (ProductResponse);
  rpc UpdateProduct (UpdateProductRequest) returns (ProductResponse);
  rpc DeleteProduct (DeleteProductRequest) returns (DeleteProductResponse);
}

message Empty {}
message Product {
  uint32 id = 1;
  string name = 2;
  double price = 3;
  int32 stock = 4;
}
message ProductListResponse {
  repeated Product products = 1;
}
message CreateProductRequest {
  string name = 1;
  double price = 2;
}
message UpdateProductRequest {
  uint32 id = 1;
  string name = 2;
  double price = 3;
  int32 stock = 4;
}
message DeleteProductRequest {
  uint32 id = 1;
}
message ProductResponse {
  Product product = 1;
}
message DeleteProductResponse {
  string message = 1;
}
```

---

## 4. Generate File Go dari .proto

Jalankan dari root folder project:
```sh
protoc --go_out=. --go-grpc_out=. proto/product.proto
```
- File hasil generate otomatis masuk ke folder `pb/`.

---

## 5. Update Layer Repository dan Service
- Implementasikan fungsi CRUD (GetAll, Create, Update, Delete, FindByID) di `repository/product_repository.go` dan `service/product_service.go`.
- Contoh implementasi lihat langkah di atas.

---

## 6. Implementasi Handler gRPC

Buat file `delivery/grpc/product_grpc_handler.go`:
```go
package grpc

import (
	"context"
	"tokokecil/model"
	"tokokecil/pb"
	"tokokecil/service"
)

type ProductGRPCHandler struct {
	pb.UnimplementedProductServiceServer
	ProductService service.ProductService
}

func (h *ProductGRPCHandler) GetAllProducts(ctx context.Context, _ *pb.Empty) (*pb.ProductListResponse, error) {
	products, err := h.ProductService.GetAllProducts()
	if err != nil {
		return nil, err
	}
	var pbProducts []*pb.Product
	for _, p := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:    uint32(p.ID),
			Name:  p.Name,
			Price: p.Price,
			Stock: int32(p.Stock),
		})
	}
	return &pb.ProductListResponse{Products: pbProducts}, nil
}

func (h *ProductGRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	prod, err := h.ProductService.CreateProduct(model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: 100,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:    uint32(prod.ID),
			Name:  prod.Name,
			Price: prod.Price,
			Stock: int32(prod.Stock),
		},
	}, nil
}

func (h *ProductGRPCHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	prod, err := h.ProductService.UpdateProduct(model.Product{
		ID:    uint(req.Id),
		Name:  req.Name,
		Price: req.Price,
		Stock: int(req.Stock),
	})
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:    uint32(prod.ID),
			Name:  prod.Name,
			Price: prod.Price,
			Stock: int32(prod.Stock),
		},
	}, nil
}

func (h *ProductGRPCHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.ProductService.DeleteProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{Message: "Product deleted successfully"}, nil
}
```

---

## 7. Modularisasi Server gRPC

Buat file `delivery/grpc/server.go`:
```go
package grpc

import (
	"log"
	"net"
	"tokokecil/pb"
	"tokokecil/service"
	"google.golang.org/grpc"
)

func StartGRPCServer(productService service.ProductService) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, &ProductGRPCHandler{
		ProductService: productService,
	})
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
```

---

## 8. Update main.go

```go
package main

import (
	"log"
	"tokokecil/config"
	grpcdelivery "tokokecil/delivery/grpc"
	"tokokecil/model"
	"tokokecil/repository"
	"tokokecil/service"
)

func main() {
	config.LoadEnv()
	db := config.DBInit()
	if err := db.AutoMigrate(&model.Product{}); err != nil {
		panic("Gagal auto-migrate tabel: " + err.Error())
	}
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	log.Println("Starting gRPC Product Service on :50051 ...")
	grpcdelivery.StartGRPCServer(productService)
}
```

---

## 9. Testing dengan grpcurl

```sh
grpcurl -plaintext -d '{}' localhost:50051 product.ProductService/GetAllProducts
grpcurl -plaintext -d '{"name":"Kopi Gayo","price":25000}' localhost:50051 product.ProductService/CreateProduct
grpcurl -plaintext -d '{"id":1,"name":"Kopi Aceh","price":35000,"stock":120}' localhost:50051 product.ProductService/UpdateProduct
grpcurl -plaintext -d '{"id":1}' localhost:50051 product.ProductService/DeleteProduct
```

---

## 10. Testing dengan Postman (gRPC)

1. **New → gRPC Request**
2. Import file `product.proto`
3. Server: `localhost:50051`
4. Pilih Service & Method
5. Masukkan Message sesuai kebutuhan:
   - Update:
     ```json
     { "id":1, "name":"Kopi Aceh", "price":35000, "stock":120 }
     ```
   - Delete:
     ```json
     { "id":1 }
     ```
6. Klik Invoke dan lihat hasilnya

---

## 11. (Opsional) Hapus/Comment HTTP/Echo
- Jika sudah tidak pakai REST, comment handler & routes Echo di main.go.

---

## 12. Maintenance & Dokumentasi
- Simpan file .proto untuk dokumentasi & auto generate stub client jika dibutuhkan.
- Dokumentasikan langkah-langkah dan endpoint di README.md project Anda.

---

## 13. Selesai!
- Project tokokecil sudah berjalan dengan full gRPC CRUD dan clean architecture.

---

