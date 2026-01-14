package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func createServer() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Slow request started...")
		time.Sleep(8 * time.Second)
		fmt.Fprintf(w, "Slow request completed at %v\n", time.Now())
	})

	return &http.Server{Addr: ":8080", Handler: mux}
}

func runServer(
	ctx context.Context,
	server *http.Server,
	shutdownTimeout time.Duration,
) error {
	serverErr := make(chan error, 1)

	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); !errors.Is(
			err, http.ErrServerClosed) {
			serverErr <- err
		}
		close(serverErr)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-serverErr:
		return err
	case <-stop:
		log.Println("Shutting down server...")
	case <-ctx.Done():
		log.Println("Context cancelled")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			return errors.Join(err, closeErr)
		}
		return err
	}
	log.Println("Server exited gracefully")
	return nil
}

func main() {
	server := createServer()

	if err := runServer(context.Background(), server, 3*time.Second); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
