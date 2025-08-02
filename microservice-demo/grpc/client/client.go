package client

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "microservice-demo/proto/user"
)

// UserClient gRPC用户客户端
type UserClient struct {
	client pb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserClient 创建用户客户端
func NewUserClient(serverAddr string) (*UserClient, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}

	client := pb.NewUserServiceClient(conn)
	return &UserClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭连接
func (c *UserClient) Close() error {
	return c.conn.Close()
}

// GetUser 获取用户
func (c *UserClient) GetUser(id int64) (*pb.GetUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.GetUser(ctx, &pb.GetUserRequest{Id: id})
}

// CreateUser 创建用户
func (c *UserClient) CreateUser(name, email, phone, password string) (*pb.CreateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.CreateUser(ctx, &pb.CreateUserRequest{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: password,
	})
}

// UpdateUser 更新用户
func (c *UserClient) UpdateUser(id int64, name, email, phone, avatar string) (*pb.UpdateUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:     id,
		Name:   name,
		Email:  email,
		Phone:  phone,
		Avatar: avatar,
	})
}

// DeleteUser 删除用户
func (c *UserClient) DeleteUser(id int64) (*pb.DeleteUserResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: id})
}

// ListUsers 获取用户列表
func (c *UserClient) ListUsers(page, pageSize int32, search string) (*pb.ListUsersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return c.client.ListUsers(ctx, &pb.ListUsersRequest{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	})
}

// ExampleUsage 使用示例
func ExampleUsage() {
	client, err := NewUserClient("localhost:9090")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// 获取用户
	user, err := client.GetUser(1)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
	} else {
		log.Printf("User: %+v", user.User)
	}

	// 创建用户
	createResp, err := client.CreateUser("张三", "zhangsan@example.com", "13800138000", "password123")
	if err != nil {
		log.Printf("Failed to create user: %v", err)
	} else {
		log.Printf("Created user: %+v", createResp.User)
	}
}