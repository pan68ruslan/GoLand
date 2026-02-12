package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
	slog.Info("write command: ")
	for rd.Scan() {
		line := rd.Text()
		ll := strings.Split(line, " ")
		if len(ll) <= 2 {
			// Prepare and send the command
			var command = cmd.NewCommand(conn)
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
			default:
				slog.Error("unknown command: "+line+"\n", "available commands", "add, get, put, del, list")
				return true
			}
			// Process the response
			if command.Value != "" {
				var response, er = command.Handle()
				rr := strings.Split(response, "|")
				if er != nil && len(rr) == 2 && rr[0] == cmd.ResponseCommandName {
					slog.Info("unknown response", "response", response)
					break
				} else if len(rr) == 2 && len(rr[1]) > 0 {
					switch command.Type {
					case cmd.AddCommandName:
						if id, err := strconv.Atoi(rr[1]); err == nil {
							doc := ds.NewDocument(c.name)
							doc.Fields["id"] = ds.DocumentField{Type: ds.DocumentFieldTypeNumber, Value: id}
							if e := c.documents.PutDocument(*doc); e == nil {
								slog.Info("new document added", "id", id)
							} else {
								slog.Error("can't add a new document", "error", e)
							}
						}
					case cmd.GetCommandName:
						var doc ds.Document
						var data = []byte(rr[1])
						if err := json.Unmarshal(data, &doc); err == nil {
							if e := c.documents.PutDocument(doc); e == nil {
								slog.Info("found document", "doc", rr[1])
							}
						}
					case cmd.PutCommandName:
						if id, err := strconv.Atoi(rr[1]); err == nil {
							slog.Info("the document was updated", "id", id)
						} else {
							slog.Error("can't put document", "error", err)
						}
					case cmd.ListCommandName:
						slog.Info("the server's list of", "documents", rr[1])
					case cmd.DeleteCommandName:
						if id, err := strconv.Atoi(rr[1]); err == nil {
							var res = ""
							if id == 0 {
								res = "wasn't"
							} else {
								res = "was"
							}
							var msg = fmt.Sprintf("the document %s deleted", res)
							slog.Info(msg, "id", id)
						} else {
							slog.Error("can't delete document", "error", err)
						}
					}
					slog.Info("the local list of", "documents:", c.documents.GetDocumentsList("3", "owner"))
				}
			} else {
				slog.Error("This command needs the second parameter.", "command", ll)
			}
		}
		slog.Info("write the next command: ")
	}
	return false
}
