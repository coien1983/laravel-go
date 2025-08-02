package main

import (
	"context"
	"log"
	"net"

	pb "microservice-demo/proto/user"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// UserServer 用户服务实现
type UserServer struct {
	pb.UnimplementedUserServiceServer
}

// GetUser 获取用户信息
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user := &pb.User{
		Id:        req.Id,
		Name:      "测试用户",
		Email:     "test@example.com",
		Phone:     "13800138000",
		Status:    "active",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	return &pb.GetUserResponse{
		User:    user,
		Message: "获取用户成功",
		Code:    200,
	}, nil
}

func main() {
	port := ":9090"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServer{})
	reflection.Register(s)

	log.Printf("🚀 gRPC Server starting on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
