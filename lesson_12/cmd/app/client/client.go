package main

import (
	"bufio"
	"encoding/json"
	"log/slog"
	"net"
	"strconv"
	"strings"

	cmd "lesson_12/internal/command"
	ds "lesson_12/internal/documentStore"
)

type Client struct {
	name      string
	documents ds.Collection
	logger    *slog.Logger
}

func NewClient(name string, logger *slog.Logger) *Client {
	docs := ds.NewCollection("Documents", logger)
	return &Client{
		name:      name,
		documents: docs,
		logger:    logger,
	}
}

func (c *Client) Connect(conn net.Conn, rd *bufio.Scanner) bool {
	//consoleReader := bufio.NewScanner(os.Stdin)
	//consoleWriter := bufio.NewWriter(os.Stdout)
	//server listener
	/*go func() {
		for {
			//msg, err := c.reader.ReadString('\n')
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				slog.Error("connection closed", "error", err)
				return
			}
			fmt.Printf("Server: %s", msg)
		}
	}()*/
	slog.Info("write command: ")
	for rd.Scan() {
		line := rd.Text()
		ll := strings.Split(line, " ")
		if len(ll) <= 2 {
			var command = cmd.NewCommand(conn)
			var doc *ds.Document = nil
			command.Type = ll[0]
			switch ll[0] {
			case cmd.AddCommandName:
				doc = ds.NewDoc(c.name)
				if msg, ew := json.Marshal(doc); ew == nil {
					//command.Value = base64.StdEncoding.EncodeToString(msg)
					command.Value = string(msg)
				}
			case cmd.GetCommandName:
				//getCommand(ll[1:])
			case cmd.PutCommandName:
				//putCommand(ll[1:])
			default:
				slog.Error("unknown command: " + line + "\n")
				return true
			}
			if command.Value != "" {
				var response, er = command.Handle()
				if er != nil {
					slog.Info("unknown response, quit", "response", response)
					break
				}
				switch command.Type {
				case cmd.AddCommandName:
					if doc != nil {
						rr := strings.Split(response, "|")
						if len(rr) == 2 && rr[0] == cmd.ResponseCommandName {
							if id, err := strconv.Atoi(rr[1]); err == nil {
								doc.Fields["id"] = ds.DocumentField{Type: ds.DocumentFieldTypeNumber, Value: id}
								c.documents.PutDocument(*doc)
							}
						}
					}
				case cmd.GetCommandName:
				case cmd.PutCommandName:

				}
			}
		}
		slog.Info("write the next command: ")
	}
	return false //command.conn != nil
}

/*func (c *Client) addCommand(args []string) {
	doc := ds.NewDoc(c.name)
	if msg, err := doc.MarshalJSON(); err != nil {
		c.writer.WriteString(fmt.Sprintf("%s %s\n", cmd.AddCommandName, base64.StdEncoding.EncodeToString([]byte(msg))))
	}
}

func (c *Client) getCommand(args []string) {

}

func (c *Client) putCommand(args []string) {

*/
