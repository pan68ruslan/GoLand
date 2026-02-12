package command

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	AddCommandName        string = "add"
	GetCommandName        string = "get"
	PutCommandName        string = "put"
	DeleteCommandName     string = "del"
	ListCommandName       string = "list"
	ConnectCommandName    string = "connect"
	DisconnectCommandName string = "disconnect"
	ResponseCommandName   string = "response"
	UnknownCommandName    string = "unknown"
	PingCommandName       string = "ping"
)

type Command struct {
	conn  net.Conn
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewCommand(c net.Conn) *Command {
	return &Command{
		conn:  c,
		Type:  UnknownCommandName,
		Value: "",
	}
}

func (c *Command) Handle() (string, error) {
	//defer c.conn.Close()
	var reader = bufio.NewReader(c.conn)
	var writer = bufio.NewWriter(c.conn)
	_, _ = writer.WriteString(fmt.Sprintf("%s|%s\n", c.Type, c.Value))
	_ = writer.Flush()
	resp, err := reader.ReadString('\n')
	return strings.TrimSpace(resp), err
}
