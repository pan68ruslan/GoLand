package main

import (
	"bufio"
	"log/slog"
	"net"
)

type Server struct {
	name string
	//Conn   net.Conn
	logger *slog.Logger
}

func NewServer(name string, logger *slog.Logger) *Server {
	return &Server{
		name:   name,
		logger: logger,
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Error("failed to close connection", "error", err)
		} else {
			s.logger.Info("connection closed", "addr", conn.RemoteAddr())
		}
	}()

	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)

	for scanner.Scan() {
		line := scanner.Text()
		s.logger.Info("received message", "addr", conn.RemoteAddr(), "msg", line)

		_, err := writer.WriteString("pong: " + line + "\n")
		if err != nil {
			s.logger.Error("failed to write response", "error", err)
			return
		}
		if err := writer.Flush(); err != nil {
			s.logger.Error("failed to flush response", "error", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		s.logger.Error("scanner error", "error", err)
	}
}
