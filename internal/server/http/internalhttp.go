package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	app "github.com/inspectorvitya/fibonacci_service/internal/application"
	"go.uber.org/zap"
)

type Server struct {
	App        *app.App
	router     *mux.Router
	HTTPServer *http.Server
}

func New(app *app.App, port string) *Server {
	router := mux.NewRouter()

	server := &Server{
		HTTPServer: &http.Server{
			Addr:         net.JoinHostPort("", port),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			Handler:      router,
		},
		router: router,
		App:    app,
	}

	return server
}

func (s *Server) Start() error {
	s.router.HandleFunc("/fibonacci", s.GetFib).Methods(http.MethodGet)
	zap.L().Info("start http server...")
	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.HTTPServer.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
