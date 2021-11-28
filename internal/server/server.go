package server

import (
	internalgrpc "github.com/inspectorvitya/fibonacci_service/internal/server/grpc"
	internalhttp "github.com/inspectorvitya/fibonacci_service/internal/server/http"
)

type Server struct {
	HTTP *internalhttp.Server
	GRPC *internalgrpc.Server
}

func New(httpServer *internalhttp.Server, grpcServer *internalgrpc.Server) *Server {
	return &Server{
		HTTP: httpServer,
		GRPC: grpcServer,
	}
}
