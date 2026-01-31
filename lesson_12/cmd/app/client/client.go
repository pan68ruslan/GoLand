package main

import (
	"log/slog"
	"net"
)

type Client struct {
	name   string
	Conn   net.Conn
	logger *slog.Logger
}

func NewClient(name string, logger *slog.Logger) *Client {
	return &Client{
		name:   name,
		logger: logger,
	}
}
