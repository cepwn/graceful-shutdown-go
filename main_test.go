package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"syscall"
	"testing"
	"time"
)

func TestServerGracefulShutdown(t *testing.T) {
	server := &http.Server{
		Addr: ":54332",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.Write([]byte("completed"))
		}),
	}

	serverErrorCh := make(chan error)
	go func() {
		serverErrorCh <- runServer(context.Background(), server, 5*time.Second)

	}()

	resp, err := http.Get("http://localhost" + server.Addr)

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	if err != nil {
		t.Fatalf("unable to send request to server: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 StatusOK, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unable to read response body: %v", err)
	}

	if string(body) != "completed" {
		t.Errorf("expected response body to be 'completed', got %s", string(body))
	}

	serverErr := <-serverErrorCh
	if serverErr != nil {
		t.Fatalf("expected no server error, got %v", serverErr)
	}
}

func TestServerTimeoutDuringShutdown(t *testing.T) {
	server := &http.Server{
		Addr: ":54331",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Second)
			w.Write([]byte("completed"))
		}),
	}

	serverErrorCh := make(chan error)
	go func() {
		serverErrorCh <- runServer(context.Background(), server, 5*time.Millisecond)
	}()

	requestErrorCh := make(chan error)
	go func() {

		_, err := http.Get("http://localhost" + server.Addr)
		requestErrorCh <- err
	}()

	time.Sleep(1 * time.Second)

	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	if <-requestErrorCh == nil {
		t.Errorf("expected client request to fail, but it succeeded")
	}

	serverErr := <-serverErrorCh

	if !errors.Is(serverErr, context.DeadlineExceeded) {
		t.Errorf("expected 'context.DeadlineExceeded' error, got %v", serverErr)
	}
}
