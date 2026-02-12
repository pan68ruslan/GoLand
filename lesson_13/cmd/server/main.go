package main

import (
	"log/slog"
	"net"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error("[Server]failed to start server", "error", err)
		return
	}
	logger.Info("[Server]server listening", "port", 8080)
	Server := NewServer("TheServer", logger)
	for {
		conn, err := listener.Accept()
		if err == nil {
			logger.Info("[Server]connection accepted", "addr", conn.RemoteAddr())
			go Server.HandleConnection(conn)
		} else {
			logger.Error("[Server]connection failed", "error", err)
			continue
		}
	}
}
