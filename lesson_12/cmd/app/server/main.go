package main

import (
	"bufio"
	"log/slog"
	"net"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Error("failed to start server", "error", err)
		return
	}
	logger.Info("server listening", "port", 8080)
	Server := NewServer("TheServer", logger)
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("failed to accept connection", "error", err)
			continue
		}
		logger.Info("new client connected", "addr", conn.RemoteAddr())
		go Server.HandleConnection(conn)
	}
}

func _handleConnection(conn net.Conn, logger *slog.Logger) {
	defer func() {
		if err := conn.Close(); err != nil {
			logger.Error("failed to close connection", "error", err)
		} else {
			logger.Info("connection closed", "addr", conn.RemoteAddr())
		}
	}()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		line := scanner.Text()
		logger.Info("received message", "addr", conn.RemoteAddr(), "msg", line)

		_, err := writer.WriteString("pong: " + line + "\n")
		if err != nil {
			logger.Error("failed to write response", "error", err)
			return
		}
		if err := writer.Flush(); err != nil {
			logger.Error("failed to flush response", "error", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("scanner error", "error", err)
	}
}
