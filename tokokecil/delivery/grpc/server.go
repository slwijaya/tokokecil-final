package grpc

import (
	"log"
	"net"
	"tokokecil/pb"
	"tokokecil/service"

	"google.golang.org/grpc"
)

// Fungsi untuk inisialisasi dan start gRPC server asd
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
