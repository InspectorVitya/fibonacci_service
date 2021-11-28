package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	app "github.com/inspectorvitya/fibonacci_service/internal/application"
	"github.com/inspectorvitya/fibonacci_service/internal/server"
	internalgrpc "github.com/inspectorvitya/fibonacci_service/internal/server/grpc"
	internalhttp "github.com/inspectorvitya/fibonacci_service/internal/server/http"
	redisdb "github.com/inspectorvitya/fibonacci_service/internal/storage/redis"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	PortHTTP  string `env:"PORT_HTTP" env-default:"8080"`
	PortGRPC  string `env:"PORT_GRPC" env-default:"8082"`
	PortRedis string `env:"REDIS_URL" env-default:":6379"`
}

func main() {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	zapLogger, err := config.Build()
	if err != nil {
		log.Fatal("failed to init logger: ", err)
	}
	zap.ReplaceGlobals(zapLogger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	ctx := context.Background()
	db := redisdb.New("fib", cfg.PortRedis)
	fibApp := app.New(db)
	if err := fibApp.Init(ctx); err != nil {
		zap.L().Fatal("err init app", zap.Error(err))
	}
	svr := server.New(internalhttp.New(fibApp, cfg.PortHTTP), internalgrpc.New(fibApp, cfg.PortGRPC))

	go func() {
		if err := svr.HTTP.Start(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				zap.L().Info("http server http stopped....")
			} else {
				zap.L().Fatal("failed to start http server: ", zap.Error(err))
			}
		}
	}()
	go func() {
		zap.L().Info("grpc server grpc starts....")
		if err := svr.GRPC.Start(); err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				zap.L().Info("http server grpc stopped....")
			} else {
				zap.L().Fatal("failed to start http server: ", zap.Error(err))
			}
		}
	}()
	<-stop
	svr.GRPC.Stop()
	ctxClose, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = svr.HTTP.Stop(ctxClose)
	if err != nil {
		zap.L().Fatal("failed server http stop: ", zap.Error(err))
	}
}
