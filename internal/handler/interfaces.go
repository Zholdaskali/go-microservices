package handler

import (
	pb "auth-service/pkg/api/service"
)

type AuthHandler interface {
	pb.AuthServiceServer
}
