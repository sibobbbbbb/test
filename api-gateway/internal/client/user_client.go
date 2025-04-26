package client

import (
	"context"
	"log"
	"time"

	pb "github.com/sibobbbbbb/backend-engineer-challenge/proto/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserClient adalah interface untuk berinteraksi dengan User Service
type UserClient interface {
	Register(ctx context.Context, username, email, password, fullname string) (*pb.User, error)
	Login(ctx context.Context, username, password string) (*pb.User, string, error)
	GetUser(ctx context.Context, id string) (*pb.User, error)
	UpdateUser(ctx context.Context, id, username, email, fullname string) (*pb.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	ValidateToken(ctx context.Context, token string) (bool, string, string, string, error)
}

type userClient struct {
	client pb.UserServiceClient
}

// NewUserClient membuat instance baru UserClient
func NewUserClient(userServiceAddr string) (UserClient, error) {
	conn, err := grpc.Dial(
		userServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Printf("Failed to connect to user service: %v", err)
		return nil, err
	}

	client := pb.NewUserServiceClient(conn)
	return &userClient{client: client}, nil
}

// Register mendaftarkan user baru
func (c *userClient) Register(ctx context.Context, username, email, password, fullname string) (*pb.User, error) {
	resp, err := c.client.Register(ctx, &pb.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
		Fullname: fullname,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

// Login melakukan login user
func (c *userClient) Login(ctx context.Context, username, password string) (*pb.User, string, error) {
	resp, err := c.client.Login(ctx, &pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, "", err
	}
	return resp.User, resp.Token, nil
}

// GetUser mendapatkan detail user
func (c *userClient) GetUser(ctx context.Context, id string) (*pb.User, error) {
	resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

// UpdateUser memperbarui data user
func (c *userClient) UpdateUser(ctx context.Context, id, username, email, fullname string) (*pb.User, error) {
	resp, err := c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
		Id:       id,
		Username: username,
		Email:    email,
		Fullname: fullname,
	})
	if err != nil {
		return nil, err
	}
	return resp.User, nil
}

// DeleteUser menghapus user
func (c *userClient) DeleteUser(ctx context.Context, id string) (bool, error) {
	resp, err := c.client.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: id,
	})
	if err != nil {
		return false, err
	}
	return resp.Success, nil
}

// ValidateToken memvalidasi token JWT
func (c *userClient) ValidateToken(ctx context.Context, token string) (bool, string, string, string, error) {
	resp, err := c.client.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, "", "", "", err
	}
	return resp.Valid, resp.UserId, resp.Username, resp.Role, nil
}