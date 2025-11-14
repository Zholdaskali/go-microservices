package grpchandler

import (
	"context"

	"auth-service/internal/logger"
	"auth-service/internal/service"
	pb "auth-service/pkg/api/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
	log         logger.Logger
}

func NewAuthHandler(authService service.AuthService, log logger.Logger) *authHandler {
	return &authHandler{
		authService: authService,
		log:         log,
	}
}

func (h *authHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := h.authService.Login(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (h *authHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// TODO: реализация регистрации
	res, err := h.authService.Register(ctx, req)

	if err != nil {
		h.log.Error("Registration failed") // Только логируем
		return nil, err
	}
	return res, nil
}

func (h *authHandler) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	// TODO: реализация валидации токена
	return nil, status.Error(codes.Unimplemented, "method ValidateToken not implemented")
}

func (h *authHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {
	// TODO: реализация обновления токена
	return nil, status.Error(codes.Unimplemented, "method RefreshToken not implemented")
}
