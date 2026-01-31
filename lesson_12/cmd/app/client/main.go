package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	cmd "lesson_12/internal/command"
)

const name string = "alpha"

func addCommand(args []string) {

}

func getCommand(args []string) {

}

func putCommand(args []string) {

}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	consoleReader := bufio.NewScanner(os.Stdin)
	consoleWriter := bufio.NewWriter(os.Stdout)

	consoleWriter.WriteString("Set the client name:")
	consoleWriter.Flush()
	consoleReader.Scan()
	name := consoleReader.Text()
	client := NewClient(name, logger)
	consoleWriter.WriteString(fmt.Sprintf("The client with %s was created", name))

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		slog.Error("failed to connect to server", "error", err)
		return
	}
	client.Conn = conn
	slog.Info("connected to server", "address", conn.RemoteAddr())

	netReader := bufio.NewReader(conn)
	netWriter := bufio.NewWriter(conn)
	//server listener
	go func() {
		for {
			msg, err := netReader.ReadString('\n')
			if err != nil {
				slog.Error("connection closed", "error", err)
				return
			}
			fmt.Printf("Server: %s", msg)
		}
	}()

	for consoleReader.Scan() {
		line := consoleReader.Text()

		ll := strings.Split(line, " ")
		if len(ll) <= 2 {
			switch ll[0] {
			case cmd.AddCommandName:
				go addCommand(ll[1:])
			case cmd.GetCommandName:
				go getCommand(ll[1:])
			case cmd.PutCommandName:
				go putCommand(ll[1:])
			default:
				_, _ = consoleWriter.WriteString("unknown command: " + line + "\n")
				_ = consoleWriter.Flush()
			}
		}

		_, _ = netWriter.WriteString(line + "\n")
		_ = netWriter.Flush()

		resp, _ := netReader.ReadString('\n')
		_, _ = consoleWriter.WriteString(resp + "\n")
		_ = consoleWriter.Flush()
	}
}
