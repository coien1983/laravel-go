package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "microservice-demo/proto/user"
)

// UserServer 用户服务实现
type UserServer struct {
	pb.UnimplementedUserServiceServer
}

// NewUserServer 创建用户服务实例
func NewUserServer() *UserServer {
	return &UserServer{}
}

// GetUser 获取用户信息
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// TODO: 实现获取用户逻辑
	user := &pb.User{
		Id:        req.Id,
		Name:      "示例用户",
		Email:     "user@example.com",
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

// CreateUser 创建用户
func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// TODO: 实现创建用户逻辑
	user := &pb.User{
		Id:        1,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    "active",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	return &pb.CreateUserResponse{
		User:    user,
		Message: "创建用户成功",
		Code:    201,
	}, nil
}

// UpdateUser 更新用户
func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO: 实现更新用户逻辑
	user := &pb.User{
		Id:        req.Id,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		Status:    "active",
		CreatedAt: "2024-01-01T00:00:00Z",
		UpdatedAt: "2024-01-01T00:00:00Z",
	}

	return &pb.UpdateUserResponse{
		User:    user,
		Message: "更新用户成功",
		Code:    200,
	}, nil
}

// DeleteUser 删除用户
func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// TODO: 实现删除用户逻辑
	return &pb.DeleteUserResponse{
		Message: "删除用户成功",
		Code:    200,
	}, nil
}

// ListUsers 获取用户列表
func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// TODO: 实现获取用户列表逻辑
	users := []*pb.User{
		{
			Id:        1,
			Name:      "用户1",
			Email:     "user1@example.com",
			Phone:     "13800138001",
			Status:    "active",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		},
		{
			Id:        2,
			Name:      "用户2",
			Email:     "user2@example.com",
			Phone:     "13800138002",
			Status:    "active",
			CreatedAt: "2024-01-01T00:00:00Z",
			UpdatedAt: "2024-01-01T00:00:00Z",
		},
	}

	return &pb.ListUsersResponse{
		Users:    users,
		Total:    2,
		Page:     req.Page,
		PageSize: req.PageSize,
		Message:  "获取用户列表成功",
		Code:     200,
	}, nil
}

// StartGRPCServer 启动gRPC服务器
func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServer{})
	
	// 启用反射服务（用于调试）
	reflection.Register(s)

	log.Printf("🚀 gRPC Server starting on %s", port)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}