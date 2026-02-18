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

type Server struct {
	name   string
	store  ds.Store
	logger *slog.Logger
}

func NewServer(name string, logger *slog.Logger) *Server {
	store := ds.NewStore(name, logger)
	return &Server{
		name:   name,
		store:  *store,
		logger: logger,
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Error("[Server]failed to close connection", "error", err)
		} else {
			s.logger.Info("[Server]connection closed", "address", conn.RemoteAddr())
		}
	}()
	netScanner := bufio.NewScanner(conn)
	netWriter := bufio.NewWriter(conn)
	for netScanner.Scan() {
		line := netScanner.Text()
		s.logger.Info("[Server]received message", "address", conn.RemoteAddr(), "msg", line)
		ll := strings.Split(line, "|")
		response := ""
		collectionName := "ERROR"
		if len(ll) == 3 {
			var currCollection *ds.Collection
			if coll, success := s.store.GetCollection(ll[2]); success {
				currCollection = coll
				collectionName = coll.Name
				s.logger.Info("[Server]try to process command", "name", ll)
				switch ll[0] {
				case cmd.AddDocumentCmd: // create new document
					var doc ds.Document
					if err := json.Unmarshal([]byte(ll[1]), &doc); err == nil {
						s.logger.Info("[Server]unmarshal document", "doc", ll[1])
						id := currCollection.MaxID()
						s.logger.Info("[Server]max document's id", "id", id)
						doc.Fields["id"] = ds.DocumentField{Type: ds.DocumentFieldTypeNumber, Value: id + 1}
						if er := currCollection.PutDocument(doc); er == nil {
							if d, e := json.Marshal(doc); e == nil {
								response = string(d)
								s.logger.Info("[Server]marshal document after update", "doc", response)
							} else {
								response = e.Error()
								collectionName = "ERROR"
								s.logger.Error("[Server]failed to marshal document", "collection", coll.Name, "error", response)
							}
						} else {
							response = er.Error()
							collectionName = "ERROR"
							s.logger.Error("[Server]failed to put document", "collection", coll.Name, "error", response)
						}
					} else {
						response = err.Error()
						collectionName = "ERROR"
						s.logger.Error("[Server]failed to unmarshal document", "collection", coll.Name, "error", response)
					}
				case cmd.GetDocumentCmd: // get existing document
					if id, err := strconv.Atoi(ll[1]); err == nil {
						if dc, ok := currCollection.GetDocument(id); ok == true {
							if d, e := json.Marshal(dc); e == nil {
								response = string(d)
								s.logger.Info("[Server]marshal document", "doc", response)
							} else {
								response = e.Error()
								collectionName = "ERROR"
								s.logger.Error("[Server]failed to marshal document", "collection", coll.Name, "error", response)
							}
						} else {
							response = "wrong document id"
							collectionName = "ERROR"
							s.logger.Error("[Server]failed to get document")
						}
					} else {
						response = err.Error()
						collectionName = "ERROR"
						s.logger.Error("[Server]failed to parse document id", "collection", coll.Name, "error", response)
					}
				case cmd.PutDocumentCmd: // update the existing document's content
					var doc ds.Document
					if err := json.Unmarshal([]byte(ll[1]), &doc); err == nil {
						s.logger.Info("[Server]unmarshal document", "doc", ll[1])
						if e := currCollection.PutDocument(doc); e == nil {
							response = fmt.Sprintf("%d", currCollection.MaxID())
						} else {
							response = e.Error()
							collectionName = "ERROR"
							s.logger.Error("[Server]failed to put document", "collection", coll.Name, "error", response)
						}
					} else {
						response = err.Error()
						collectionName = "ERROR"
						s.logger.Error("[Server]failed to unmarshal document", "error", response)
					}
				case cmd.DelDocumentCmd: // delete existing doc
					if id, e := strconv.Atoi(ll[1]); e == nil {
						if ok := currCollection.DeleteDocument(id); ok == true {
							s.logger.Info("[Server]document ", "id", id)
							response = fmt.Sprintf("%d", id)
						} else {
							response = "failed to delete document"
							collectionName = "ERROR"
							s.logger.Error("[Server]failed to delete document", "id", id)
						}
					} else {
						response = e.Error()
						collectionName = "ERROR"
						s.logger.Error("[Server]failed to parse document id", "collection", coll.Name, "error", response)
					}
				case cmd.DocumentsListCmd: // get document's list (0 - all documents, N - first N documents)
					response = currCollection.GetDocumentsList(ll[1], "owner")
					collectionName = coll.Name
				case cmd.DelCollectionCmd:
					if ok := s.store.DeleteCollection(ll[1]); ok == true {
						response = fmt.Sprintf("%s", ll[1])
					} else {
						response = fmt.Sprintf("Collection %s wasn't deleted", ll[1])
						collectionName = "ERROR"
					}
				default:
					s.logger.Error(fmt.Sprintln("[Server]unknown command: ", line))
				}
			} else {
				if ll[0] == cmd.AddCollectionCmd {
					_, c := s.store.CreateCollection(ll[1], s.logger) // add new collection
					response = c.Name
					collectionName = c.Name
				} else if ll[0] == cmd.CollectionsListCmd { // get collection's list (0 - all documents, N - first N documents)
					response = s.store.GetCollectionList(ll[1])
					collectionName = ""
				}
			}
		}
		_, err := netWriter.WriteString(fmt.Sprintf("%s|%s|%s\n", cmd.ResponseCommandName, response, collectionName))
		if err != nil {
			s.logger.Error("[Server]failed to write response", "errorMsg", err)
			return
		}
		if err := netWriter.Flush(); err != nil {
			s.logger.Error("[Server]failed to flush response", "errorMsg", err)
			return
		}
	}
	if err := netScanner.Err(); err != nil {
		s.logger.Error("[Server]scanner error", "error", err)
	}
}
