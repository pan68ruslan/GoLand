package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	reader := bufio.NewScanner(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Set the client name:")
	writer.Flush()
	reader.Scan()
	name := reader.Text()
	client := NewClient(name, logger)
	slog.Info(fmt.Sprintf("The client with %s was created", name))
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		slog.Error("failed to connect to server", "error", err)
		return
	}
	slog.Info("connect to server", "address", conn.RemoteAddr())
	for {
		if r := client.Connect(conn, reader); r {
			break
		}
	}
}
