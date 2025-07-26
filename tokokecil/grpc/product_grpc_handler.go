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

// GetAllProducts implementasi dari proto
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

// CreateProduct implementasi dari proto
func (h *ProductGRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	prod, err := h.ProductService.CreateProduct(model.Product{
		Name:  req.Name,
		Price: req.Price,
		Stock: 100, // default
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

// UpdateProduct
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

// DeleteProduct
func (h *ProductGRPCHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := h.ProductService.DeleteProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteProductResponse{Message: "Product deleted successfully"}, nil
}
