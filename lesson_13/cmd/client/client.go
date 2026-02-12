package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"strings"

	cmd "lesson_13/internal/command"
	ds "lesson_13/internal/documentStore"
)

type Client struct {
	conn      net.Conn
	name      string
	documents ds.Collection
	logger    *slog.Logger
}

func NewClient(logger *slog.Logger) *Client {
	docs := ds.NewCollection("Documents", logger)
	return &Client{
		name:      "noname",
		documents: docs,
		logger:    logger,
	}
}

func (c *Client) Connect(address string, reader *bufio.Scanner) error {
	if c.conn != nil {
		c.logger.Info("Already connected", "address", c.conn.RemoteAddr())
		return nil
	}
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Set the client name: ")
	writer.Flush()
	reader.Scan()
	c.name = reader.Text()
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	c.conn = conn
	c.logger.Info("Connected to...", "address", address)
	return nil
}

func (c *Client) Disconnect() error {
	if c.conn == nil {
		return fmt.Errorf("no active connection")
	}
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to disconnect: %w", err)
	}
	c.logger.Info("Disconnected")
	c.conn = nil
	return nil
}

func (c *Client) ProcessResponse(command *cmd.Command) {
	var response, er = command.Handle()
	rr := strings.Split(response, "|")
	if er != nil && rr[0] != cmd.ResponseCommandName {
		c.logger.Info("unknown response, quit", "response", response)
	} else {
		if len(rr) < 2 || len(rr[1]) == 0 {
			c.logger.Info("empty response was received")
		}
	}
	if len(rr) == 2 && len(rr[1]) > 0 {
		switch command.Type {
		case cmd.AddCommandName:
			if id, err := strconv.Atoi(rr[1]); err == nil {
				doc := ds.NewDocument(c.name)
				doc.Fields["id"] = ds.DocumentField{Type: ds.DocumentFieldTypeNumber, Value: id}
				if e := c.documents.PutDocument(*doc); e == nil {
					c.logger.Info("new document added", "id", id)
				} else {
					c.logger.Error("can't add a new document", "error", e)
				}
			}
		case cmd.GetCommandName:
			var doc ds.Document
			var data = []byte(rr[1])
			if err := json.Unmarshal(data, &doc); err == nil {
				if e := c.documents.PutDocument(doc); e == nil {
					c.logger.Info("found document", "doc", rr[1])
				}
			}
		case cmd.PutCommandName:
			if id, err := strconv.Atoi(rr[1]); err == nil {
				c.logger.Info("the document was updated", "id", id)
			} else {
				c.logger.Error("can't put document", "error", err)
			}
		case cmd.ListCommandName:
			c.logger.Info("the server's list of", "documents", rr[1])
		case cmd.DeleteCommandName:
			if id, err := strconv.Atoi(rr[1]); err == nil {
				var res = ""
				if id == 0 {
					res = "wasn't"
				} else {
					res = "was"
				}
				var msg = fmt.Sprintf("the document %s deleted", res)
				c.logger.Info(msg, "id", id)
			} else {
				c.logger.Error("can't delete document", "error", err)
			}
		}
		c.logger.Info("the local list of", "documents:", c.documents.GetDocumentsList("owner"))
	}
}

func (c *Client) Start() {
	reader := bufio.NewScanner(os.Stdin)
	if c.conn == nil {
		err := c.Connect(cmd.Address, reader)
		if err == nil {
			c.logger.Info("client connected successfully", "address", c.conn.RemoteAddr())
		} else {
			c.logger.Error("fail at connecting", "error", err)
		}
	}
	c.logger.Info("write command: ")
	for reader.Scan() {
		line := reader.Text()
		ll := strings.Split(line, " ")
		if len(ll) <= 2 {
			// Prepare and send the command
			var command = cmd.NewCommand(c.conn)
			command.Type = ll[0]
			if len(ll) == 2 {
				command.Value = fmt.Sprintf("%s", ll[1])
			}
			switch ll[0] {
			case cmd.AddCommandName:
				var doc = ds.NewDocument(c.name)
				if msg, e := json.Marshal(doc); e == nil {
					command.Value = string(msg)
				}
			case cmd.GetCommandName:
				break
			case cmd.PutCommandName:
				if len(ll) == 2 {
					if id, err := strconv.Atoi(ll[1]); err == nil {
						if doc, ok := c.documents.GetDocument(id); ok {
							if err = doc.UpdateContent(c.name); err == nil {
								if msg, e := json.Marshal(doc); e == nil {
									command.Value = string(msg)
								}
							}
						}
					}
				}
			case cmd.ListCommandName:
				if len(ll) == 1 {
					command.Value = "0"
				}
			case cmd.DeleteCommandName:
				break
			case cmd.ConnectCommandName:
				if len(ll) == 1 {
					_ = c.Connect(cmd.Address, reader)
				} else {
					_ = c.Connect(ll[1], reader)
				}
			case cmd.DisconnectCommandName:
				_ = c.Disconnect()
			default:
				c.logger.Error("unknown command: "+line+"\n", "available commands", "add, get, put, del, list, connect, disconnect")
				continue
			}
			// Process the response
			if command.Value != "" {
				c.ProcessResponse(command)
			} else {
				c.logger.Error("This command needs the second parameter.", "command", ll)
			}
		}
		c.logger.Info("write the next command: ")
	}
}
