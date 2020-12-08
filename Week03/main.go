package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/sync/errgroup"
)

func createHttpServer() *http.Server {
	listenAddr := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "pong")
	})

	server := http.Server{
		Addr:         listenAddr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return &server
}

func main() {
	logger := log.New(os.Stdout, "week03: ", log.LstdFlags)
	ctx, cancel := context.WithCancel(context.Background())
	g := errgroup.WithContext(ctx)
	server := createHttpServer()

	// 监听信号
	g.Go(func(gCtx context.Context) error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			logger.Printf("Received signal: %s\n", sig)
			cancel()
		case <-gCtx.Done():
			logger.Println("closing signal goroutine")
			return gCtx.Err()
		}
		return nil
	})

	// 启动HTTP服务
	g.Go(func(context.Context) error {

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Println(err)
			cancel()
		}
		logger.Println("http server stopped")
		return nil
	})

	//关闭HTTP服务
	g.Go(func(gCtx context.Context) error {

		select {
		case <-gCtx.Done():
			break
		}

		timeoutCtx, timeoutCancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer timeoutCancel()

		logger.Println("shutting down http server, please wait...")

		return server.Shutdown(timeoutCtx)

	})

	// wait for shutdown
	if err := g.Wait(); err != nil {
		logger.Println(err.Error())
	}

	logger.Println("gracefully shoutdown server")
}
