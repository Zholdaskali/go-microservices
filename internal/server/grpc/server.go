package grpc

import (
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	pb "auth-service/pkg/api/service"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	log    logger.Logger
}

func NewServer(authHandler handler.AuthHandler, log logger.Logger) *Server {
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	return &Server{
		server: grpcServer,
		log:    log,
	}
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	s.log.Info("gRPC server starting", logger.F("port", port))

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.log.Fatal("gRPC server failed", logger.F("error", err))
		}
	}()

	return nil
}

func (s *Server) Stop() {
	s.server.Stop()
}
