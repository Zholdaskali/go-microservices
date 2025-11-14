package service

import (
	pb "auth-service/pkg/api/service"
)

type AuthService interface {
	pb.AuthServiceServer
}
