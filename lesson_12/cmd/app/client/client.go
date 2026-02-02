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

func (c *Client) getDocuments() string {
	result := ""
	for _, d := range c.documents.Documents {
		if len(result) > 0 {
			result += ", "
		}
		id, _ := ds.ToInt(d.Fields["id"].Value)
		result += fmt.Sprintf("%d", id)
	}
	return result
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
			switch ll[0] {
			case cmd.AddCommandName:
				var doc = ds.NewDoc(c.name)
				if msg, ew := json.Marshal(doc); ew == nil {
					command.Value = string(msg)
				}
			case cmd.GetCommandName:
				command.Value = fmt.Sprintf("%s", ll[1])
			case cmd.PutCommandName:
				if id, ei := strconv.Atoi(ll[1]); ei == nil {
					if doc, ok := c.documents.GetDocument(id); ok {
						if err := doc.UpdateContent(c.name); err == nil {
							if msg, ew := json.Marshal(doc); ew == nil {
								command.Value = string(msg)
							}
						}
					}
				}
			case cmd.DeleteCommandName:
				command.Value = fmt.Sprintf("%s", ll[1])
			default:
				slog.Error("unknown command: " + line + "\n")
				return true
			}
			// Process the response
			if command.Value != "" {
				var response, er = command.Handle()
				rr := strings.Split(response, "|")
				if er != nil && len(rr) == 2 && rr[0] == cmd.ResponseCommandName && (len(rr[1]) > 0) {
					slog.Info("unknown response, quit", "response", response)
					break
				} else {
					switch command.Type {
					case cmd.AddCommandName:
						if id, err := strconv.Atoi(rr[1]); err == nil {
							doc := ds.NewDoc(c.name)
							doc.Fields["id"] = ds.DocumentField{Type: ds.DocumentFieldTypeNumber, Value: id}
							if e := c.documents.PutDocument(*doc); e == nil {
								slog.Info("new document added", "id", id)
							} else {
								slog.Error("can't put document", "error", e)
							}
						}
					case cmd.GetCommandName:
						var doc ds.Document
						var data = []byte(rr[1])
						if err := json.Unmarshal(data, &doc); err == nil {
							if e := c.documents.PutDocument(doc); e == nil {
								slog.Info("found document, quit", "doc", rr[1])
							}
						}
					case cmd.PutCommandName:
						if id, err := strconv.Atoi(rr[1]); err == nil {
							slog.Info("the document was updated", "id", id)
						} else {
							slog.Error("can't put document", "error", err)
						}
					case cmd.DeleteCommandName:
						if id, err := strconv.Atoi(rr[1]); err == nil {
							slog.Info("the document was deleted", "id", id)
						} else {
							slog.Error("can't delete document", "error", err)
						}
					}
					slog.Info("List of ", "documents:", c.getDocuments())
				}
			}
		}
		slog.Info("write the next command: ")
	}
	return false
}
