package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/drobyshevv/classifer-gateway/internal/app"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Разные адреса для разных сервисов
	ssoGRPCAddr := getEnv("SSO_GRPC_ADDR", "localhost:44044")       // SSO gRPC
	expertGRPCAddr := getEnv("EXPERT_GRPC_ADDR", "localhost:50051") // Expert-search gRPC
	httpAddr := getEnv("HTTP_ADDR", ":8080")                        // Gateway HTTP

	// Передаем оба адреса
	application := app.New(log, ssoGRPCAddr, expertGRPCAddr, httpAddr)

	go application.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Info("stopping gateway", slog.String("signal", sign.String()))
	application.Stop()
	log.Info("gateway stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
