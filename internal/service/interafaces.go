package service

import (
	pb "github.com/Zholdaskali/go-microservices-proto/pkg/api/auth-service"
)

type AuthService interface {
	pb.AuthServiceServer
}
