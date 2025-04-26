package grpc

import (
    "context"

    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    "github.com/sibobbbbbb/backend-engineer-challenge/user-service/internal/domain"
    pb "github.com/sibobbbbbb/backend-engineer-challenge/proto/user"
)

type UserHandler struct {
    pb.UnimplementedUserServiceServer
    userUsecase domain.UserUsecase
}

// NewUserHandler membuat instance baru UserHandler
func NewUserHandler(uu domain.UserUsecase) *UserHandler {
    return &UserHandler{
        userUsecase: uu,
    }
}

// Register menangani permintaan pendaftaran pengguna baru
func (h *UserHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
    user, err := h.userUsecase.Register(ctx, req.Username, req.Email, req.Password, req.Fullname)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to register: %v", err)
    }

    return &pb.UserResponse{
        User: convertDomainToProto(user),
    }, nil
}

// Login menangani permintaan login
func (h *UserHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
    user, token, err := h.userUsecase.Login(ctx, req.Username, req.Password)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
    }

    return &pb.LoginResponse{
        User:  convertDomainToProto(user),
        Token: token,
    }, nil
}

// GetUser menangani permintaan untuk mendapatkan detail pengguna
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
    user, err := h.userUsecase.GetByID(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
    }

    return &pb.UserResponse{
        User: convertDomainToProto(user),
    }, nil
}

// UpdateUser menangani permintaan pembaruan profil pengguna
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
    user, err := h.userUsecase.Update(ctx, req.Id, req.Username, req.Email, req.Fullname)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
    }

    return &pb.UserResponse{
        User: convertDomainToProto(user),
    }, nil
}

// DeleteUser menangani permintaan penghapusan pengguna
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    err := h.userUsecase.Delete(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
    }

    return &pb.DeleteUserResponse{
        Success: true,
    }, nil
}

// ValidateToken memvalidasi token JWT
func (h *UserHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
    valid, userID, username, role, err := h.userUsecase.ValidateToken(ctx, req.Token)
    if err != nil {
        return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
    }

    return &pb.ValidateTokenResponse{
        Valid:    valid,
        UserId:   userID,
        Username: username,
        Role:     role,
    }, nil
}

// convertDomainToProto mengonversi domain.User ke pb.User
func convertDomainToProto(user *domain.User) *pb.User {
    return &pb.User{
        Id:        user.ID,
        Username:  user.Username,
        Email:     user.Email,
        Fullname:  user.Fullname,
        Role:      user.Role,
        CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }
}