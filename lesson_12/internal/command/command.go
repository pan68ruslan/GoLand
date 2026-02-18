package command

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const (
	AddCollectionCmd      string = "add-col"    // create and add collection : 								add-col name
	DelCollectionCmd      string = "del-col"    // delete the collection : 									del-col name
	SetCollectionCmd      string = "set-col"    // set the collection as current : 							set-coll name
	CollectionsListCmd    string = "lst-col"    // get the collection's list from server:					col-list number { default 0 - all connections }
	AddDocumentCmd        string = "add-doc"    // create new document and add to the current collection: 	add-doc
	GetDocumentCmd        string = "get-doc"    // get the document from the current collection at server:	get-doc id
	PutDocumentCmd        string = "put-doc"    // update and put local document to the current collection at server: put-doc id
	DelDocumentCmd        string = "del-doc"    // delete the document from the current collection: 		del-doc id
	DocumentsListCmd      string = "lst-doc"    // get the document's list from the current collection: 	doc-lst number { default 0 - all documents }
	ConnectCommandName    string = "connect"    // connect to the server:									connect address { default 8080 }
	DisconnectCommandName string = "disconnect" // disconnect from the server:								disconnect
	ResponseCommandName   string = "response"
	UnknownCommandName    string = "unknown"
	//PingCommandName       string = "ping"
)

const Address string = ":8080"

type Command struct {
	conn    net.Conn
	Type    string `json:"type"`
	Value   string `json:"value"`
	ColName string `json:"colname"`
}

func NewCommand(c net.Conn) *Command {
	return &Command{
		conn:    c,
		Type:    UnknownCommandName,
		Value:   "0",
		ColName: "",
	}
}

func (c *Command) Handle() (string, error) {
	//defer c.conn.Close()
	var reader = bufio.NewReader(c.conn)
	var writer = bufio.NewWriter(c.conn)
	_, _ = writer.WriteString(fmt.Sprintf("%s|%s|%s\n", c.Type, c.Value, c.ColName))
	_ = writer.Flush()
	resp, err := reader.ReadString('\n')
	return strings.TrimSpace(resp), err
}
