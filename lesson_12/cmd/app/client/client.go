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

	cmd "lesson_12/internal/command"
	ds "lesson_12/internal/documentStore"
)

type Client struct {
	name           string               `json:"Name"`
	store          ds.Store             `json:"Store"`
	cfg            *ds.CollectionConfig `json:"Config"`
	currCollection *ds.Collection
	conn           net.Conn
	logger         *slog.Logger
}

func NewClient(logger *slog.Logger) *Client {
	store := ds.NewStore("ClientStore", logger)
	cfg := ds.NewConfig()
	return &Client{
		name:           "noname",
		store:          *store,
		cfg:            cfg,
		currCollection: nil,
		conn:           nil,
		logger:         logger,
	}
}

func (c *Client) Connect(address string, reader *bufio.Scanner) {
	if c.conn != nil {
		c.logger.Info("Already connected", "address", c.conn.RemoteAddr())
	}
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Set the client name: ")
	writer.Flush()
	reader.Scan()
	c.name = reader.Text()
	conn, err := net.Dial("tcp", address)
	if err == nil {
		c.conn = conn
		c.logger.Info("client connected successfully", "address", c.conn.RemoteAddr())
	} else {
		c.logger.Error("fail at connecting", "error", err)
	}
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

func (c *Client) Start() {
	reader := bufio.NewScanner(os.Stdin)
	c.Connect(cmd.Address, reader)
	c.logger.Info("write command: ")
	for reader.Scan() {
		line := reader.Text()
		ll := strings.Split(line, " ")
		if len(ll) <= 2 {
			// Prepare the command parameters
			var command = cmd.NewCommand(c.conn)
			command.Type = ll[0]
			if len(ll) == 2 && ll[1] != "" {
				command.Value = ll[1]
			}
			if c.currCollection != nil {
				command.ColName = c.currCollection.Name
			}
			c.logger.Info("command type", "type", command.Type)
			switch command.Type {
			case cmd.DelCollectionCmd:
			case cmd.DelDocumentCmd:
			case cmd.GetDocumentCmd:
			case cmd.DocumentsListCmd:
				break
			case cmd.AddCollectionCmd:
				command.ColName = ""
			case cmd.CollectionsListCmd:
				command.ColName = ""
			case cmd.SetCollectionCmd:
				if coll, ok := c.store.GetCollection(ll[1]); ok {
					c.currCollection = coll
					c.logger.Info("set the current", "collection", c.currCollection.Name)
				}
				command.Value = ""
			case cmd.AddDocumentCmd:
				if c.currCollection != nil {
					doc := ds.NewDocument(c.name)
					if msg, e := json.Marshal(doc); e == nil {
						command.Value = string(msg)
					} else {
						c.logger.Error("add marshal document error", "error", e)
					}
				}
			case cmd.PutDocumentCmd:
				if c.currCollection != nil && len(ll) == 2 {
					if id, err := strconv.Atoi(ll[1]); err == nil {
						if doc, ok := c.currCollection.GetDocument(id); ok {
							if err = doc.UpdateContent(c.name); err == nil {
								if msg, e := json.Marshal(doc); e == nil {
									command.Value = string(msg)
								} else {
									c.logger.Error("put marshal document error", "error", e)
								}
							} else {
								c.logger.Error("fail at update document", "error", err)
							}
						}
					} else {
						c.logger.Error("fail at convert the document id", "error", err)
					}
				}
			case cmd.ConnectCommandName:
				var address string
				if len(ll) == 1 {
					address = cmd.Address
				} else {
					address = ll[1]
				}
				c.Connect(address, reader)
				command.Value = ""
			case cmd.DisconnectCommandName:
				if e := c.Disconnect(); e == nil {
					c.logger.Info("client disconnected successfully")
				} else {
					c.logger.Error("fail at disconnecting", "error", e)
				}
				command.Value = ""
			default:
				c.logger.Error("unknown command", "command", ll)
			}
			if command.Value == "0" && !(command.Type == cmd.AddDocumentCmd || command.Type == cmd.DocumentsListCmd || command.Type == cmd.AddCollectionCmd || command.Type == cmd.CollectionsListCmd) {
				c.logger.Error("This command needs the second parameter.", "current command:", ll)
			} else {
				if c.conn == nil {
					c.logger.Error("There are no active connection")
				} else if command.Value != "" && command.Type != cmd.UnknownCommandName {
					// Execute command and process the response
					c.ProcessResponse(command)
				}
			}
		} else {
			c.logger.Error("entered command is too long", "command", ll)
		}
		c.logger.Info("write the next command: ")
	}
}

func (c *Client) ProcessResponse(command *cmd.Command) {
	if response, err := command.Handle(); err == nil {
		rr := strings.Split(response, "|")
		if len(rr) == 3 && rr[0] == cmd.ResponseCommandName {
			if command.Type == cmd.AddCollectionCmd {
				if ok, cl := c.store.CreateCollection(rr[2], c.logger); ok {
					c.currCollection = cl
					c.logger.Info("Collection was created", "name", rr[2])
				}
			}
			if len(rr[2]) > 0 {
				if rr[2] == "ERROR" {
					c.logger.Error("last response get the error", "error", rr[1])
				} else {
					if col, success := c.store.GetCollection(rr[2]); success {
						c.currCollection = col
					} else {
						c.currCollection = nil
						c.logger.Error("failed to find collection", "name", rr[2])
					}
				}
			}
			if rr[2] == "" {
				switch command.Type {
				case cmd.AddCollectionCmd:
					break
				case cmd.DelCollectionCmd:
					if ok := c.store.DeleteCollection(rr[1]); ok {
						c.logger.Info("the collection was deleted", "name", rr[1])
					} else {
						c.logger.Error("the collection was not deleted", "name", rr[1])
					}
				case cmd.CollectionsListCmd:
					c.logger.Info("the server's list of", "collections", rr[1])
				}
			} else if c.currCollection != nil {
				switch command.Type {
				case cmd.DelDocumentCmd:
					if id, e := strconv.Atoi(rr[1]); e == nil {
						var res = ""
						if ok := c.currCollection.DeleteDocument(id); ok {
							res = "was"
						} else {
							res = "wasn't"
						}
						var msg = fmt.Sprintf("the document %s deleted in collection %s", res, c.currCollection.Name)
						c.logger.Info(msg, "id", id)
					} else {
						c.logger.Error("can't parse document id", "error", e)
					}
				case cmd.AddDocumentCmd:
					doc := ds.NewDocument(c.name)
					var data = []byte(rr[1])
					e := json.Unmarshal(data, &doc)
					if e == nil {
						e = c.currCollection.PutDocument(*doc)
						if e == nil {
							c.logger.Info("new document added")
						} else {
							c.logger.Error("failed to add document")
						}
					} else {
						c.logger.Error("failed to marshal document", "error", e)
					}
				case cmd.GetDocumentCmd:
					var doc ds.Document
					var data = []byte(rr[1])
					e := json.Unmarshal(data, &doc)
					if e == nil {
						if id, ei := ds.ToInt(doc.Fields["id"]); ei == nil {
							e = c.currCollection.PutDocument(doc)
							if e == nil {
								c.logger.Info("found document", "doc", doc)
							} else {
								c.logger.Error("failed to get document", "id", id)
							}
						} else {
							c.logger.Error("failed to get document's id", "error", ei)
						}
					} else {
						c.logger.Error("failed to marshal document", "error", e)
					}
				case cmd.PutDocumentCmd:
					if id, e := strconv.Atoi(rr[1]); e == nil {
						c.logger.Info("the document was updated", "id", id)
					} else {
						c.logger.Error("can't put document", "error", e)
					}
				case cmd.DocumentsListCmd:
					c.logger.Info("the server's list of", "documents", rr[1])
				}
				c.logger.Info("the local list of", "documents:", c.currCollection.GetDocumentsList("owner"), "current connection", c.currCollection.Name)
			}
			c.logger.Info("the local list of", "collections:", c.store.GetCollectionList())
		} else {
			c.logger.Error("the wrong response command", "command", rr[0])
		}
	} else {
		c.logger.Error("cannot parse the response", "response", response, "error", err)
	}
}
