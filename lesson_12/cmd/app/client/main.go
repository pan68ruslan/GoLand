package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	client := NewClient(logger)
	client.logger.Info(fmt.Sprintf("The client was created"))
	client.Start()
	client.logger.Info(fmt.Sprintf("The client was closed"))
}
