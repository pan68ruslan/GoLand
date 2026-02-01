package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
)

const nameCl string = "alpha"

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
	//client.Conn = conn
	slog.Info("connect to server", "address", conn.RemoteAddr())

	for {
		if r := client.Connect(conn, reader); r {
			break
		}
	}

	/*	netReader := bufio.NewReader(conn)
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
	*/
}
