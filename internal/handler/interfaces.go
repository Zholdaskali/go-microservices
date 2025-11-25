package handler

import (
	pb "github.com/Zholdaskali/go-microservices-proto/pkg/api/auth-service"
)

type AuthHandler interface {
	pb.AuthServiceServer
}
