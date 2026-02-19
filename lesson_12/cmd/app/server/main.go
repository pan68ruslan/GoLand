package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error("[Server]failed to start server", "error", err)
		return
	}
	logger.Info("[Server]server listening", "port", 8080)
	server := NewServer("TheServer", logger)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				logger.Error("[Server]connection failed", "error", err)
				return
			}
			logger.Info("[Server]connection accepted", "addr", conn.RemoteAddr())
			wg.Add(1)
			go func() {
				defer wg.Done()
				server.HandleConnection(conn)
			}()
		}
	}()
	<-stop
	logger.Info("[Server]shutdown signal received")
	if err := listener.Close(); err != nil {
		logger.Error("[Server]failed to close listener", "error", err)
	}
	wg.Wait()
	logger.Info("[Server]graceful shutdown complete")
}
