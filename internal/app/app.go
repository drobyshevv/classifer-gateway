package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ssov1 "github.com/drobyshevv/protos/gen/go/proto/sso"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
}

func New(log *slog.Logger, grpcAddr, httpAddr string) *App {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// handlers for Auth service
	err := ssov1.RegisterAuthHandlerFromEndpoint(context.Background(), mux, grpcAddr, opts)
	if err != nil {
		panic(fmt.Sprintf("failed to register gateway: %v", err))
	}

	return &App{
		log: log,
		httpServer: &http.Server{
			Addr:    httpAddr,
			Handler: mux,
		},
	}
}

func (a *App) MustRun() {
	a.log.Info("starting HTTP gateway",
		slog.String("http_addr", a.httpServer.Addr),
	)

	if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("failed to start gateway: %v", err))
	}
}

func (a *App) Stop() {
	a.log.Info("stopping HTTP gateway")

	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		a.log.Error("failed to stop gateway", slog.String("error", err.Error()))
	}
}
