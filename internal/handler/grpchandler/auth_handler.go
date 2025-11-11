package grpchandler

import (
	"context"
	"log"

	"auth-service/internal/logger"
	"auth-service/internal/service"
	pb "auth-service/pkg/api/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
	log         logger.Logger
}

func NewAuthHandler(authService *service.AuthService, log logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		log:         log,
	}
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := h.authService.Login(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// TODO: реализация регистрации
	res, err := h.authService.Register(ctx, req)

	if err != nil {
		log.Fatal("Ошибка при сохранении")
		return nil, err
	}
	return res, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	// TODO: реализация валидации токена
	return nil, status.Error(codes.Unimplemented, "method ValidateToken not implemented")
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {
	// TODO: реализация обновления токена
	return nil, status.Error(codes.Unimplemented, "method RefreshToken not implemented")
}
