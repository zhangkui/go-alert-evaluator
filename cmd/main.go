package main

import (
	"context"
	"errors"
	"github.com/zhangkui/go-alert-evaluator/internal/clock"
	"github.com/zhangkui/go-alert-evaluator/internal/httpapi"
	"github.com/zhangkui/go-alert-evaluator/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := service.New(clock.Real{}, nil, 24*time.Hour)
	server := &http.Server{Addr: ":8080", Handler: httpapi.New(app), ReadHeaderTimeout: 5 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()
	log.Printf("listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
