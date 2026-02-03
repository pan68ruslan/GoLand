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
	name      string
	documents ds.Collection
	logger    *slog.Logger
}

func NewServer(name string, logger *slog.Logger) *Server {
	docs := ds.NewCollection("Documents", logger)
	return &Server{
		name:      name,
		documents: docs,
		logger:    logger,
	}
}

func (s *Server) HandleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			s.logger.Error("[Server]failed to close connection", "error", err)
		} else {
			s.logger.Info("[Server]connection closed", "addr", conn.RemoteAddr())
		}
	}()
	netScanner := bufio.NewScanner(conn)
	netWriter := bufio.NewWriter(conn)
	for netScanner.Scan() {
		line := netScanner.Text()
		s.logger.Info("[Server]received message", "addr", conn.RemoteAddr(), "msg", line)
		ll := strings.Split(line, "|")
		response := ""
		if len(ll) == 2 {
			switch ll[0] {
			case cmd.AddCommandName: // create new document
				s.logger.Info("[Server]add command message", "name", ll)
				var doc ds.Document
				if err := json.Unmarshal([]byte(ll[1]), &doc); err == nil {
					s.logger.Info("[Server]unmarshal document", "doc", ll[1])
					id := s.documents.MaxId()
					s.logger.Info("[Server]max document's id", "id", id)
					doc.Fields["id"] = ds.DocumentField{Type: ds.DocumentFieldTypeNumber, Value: id + 1}
					if e := s.documents.PutDocument(doc); e == nil {
						response = fmt.Sprintf("%d", s.documents.MaxId())
					}
				} else {
					s.logger.Error("[Server]failed to unmarshal document", "error", err)
				}
			case cmd.GetCommandName: // get existing document
				if id, e := strconv.Atoi(ll[1]); e == nil {
					if dc, ok := s.documents.GetDocument(id); ok == true {
						if d, err := json.Marshal(dc); err == nil {
							response = string(d)
							s.logger.Info("[Server]marshal document", "doc", response)
						} else {
							s.logger.Error("[Server]failed to marshal document", "error", err)
						}
					}
				}
			case cmd.PutCommandName: // update the existing document's content
				s.logger.Info("[Server]add command message", "name", ll)
				var doc ds.Document
				if err := json.Unmarshal([]byte(ll[1]), &doc); err == nil {
					s.logger.Info("[Server]unmarshal document", "doc", ll[1])
					if e := s.documents.PutDocument(doc); e == nil {
						response = fmt.Sprintf("%d", s.documents.MaxId())
					}
				} else {
					s.logger.Error("[Server]failed to unmarshal document", "error", err)
				}
			case cmd.ListCommandName: // get document's list (0 - all documents, N - first N documents)
				response = s.documents.GetDocumentsList(ll[1], "owner")
			case cmd.DeleteCommandName: // delete existing doc
				if id, e := strconv.Atoi(ll[1]); e == nil {
					if ok := s.documents.DeleteDocument(id); ok == true {
						s.logger.Info("[Server]document ", "id", id)
						response = fmt.Sprintf("%d", id)
					} else {
						s.logger.Error("[Server]failed to find document", "id", id)
						response = "0"
					}
				}
			default:
				s.logger.Error(fmt.Sprintln("[Server]unknown command: ", line))
			}
		}
		_, err := netWriter.WriteString(fmt.Sprintf("%s|%s\n", cmd.ResponseCommandName, response))
		if err != nil {
			s.logger.Error("[Server]failed to write response", "error", err)
			return
		}
		if err := netWriter.Flush(); err != nil {
			s.logger.Error("[Server]failed to flush response", "error", err)
			return
		}
	}
	if err := netScanner.Err(); err != nil {
		s.logger.Error("[Server]scanner error", "error", err)
	}
}
