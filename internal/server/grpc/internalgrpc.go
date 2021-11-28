package internalgrpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/inspectorvitya/fibonacci_service/api/fibanacci"
	app "github.com/inspectorvitya/fibonacci_service/internal/application"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	app        *app.App
	GRPCServer *grpc.Server
	pb.UnimplementedFibonacciServer
	port string
}

func New(fibSvc *app.App, port string) *Server {
	return &Server{
		app:        fibSvc,
		GRPCServer: grpc.NewServer(),
		port:       port,
	}
}

func (s *Server) GetFibSlice(ctx context.Context, in *pb.FibonacciRequest) (*pb.FibonacciResponse, error) {
	res, err := s.app.GetFibSlice(ctx, int(in.GetX()), int(in.GetY()))
	if err != nil {
		zap.L().Error("Ger fib err: ", zap.Error(err))
		return nil, err
	}
	return &pb.FibonacciResponse{Message: res}, nil
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("start grpc server failed: %w", err)
	}
	pb.RegisterFibonacciServer(s.GRPCServer, s)
	return s.GRPCServer.Serve(listener)
}

func (s *Server) Stop() {
	s.GRPCServer.GracefulStop()
}
