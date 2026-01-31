package command

const (
	AddCommandName        string = "add"
	GetCommandName        string = "get"
	PutCommandName        string = "put"
	DeleteCommandName     string = "del"
	ListCommandName       string = "list"
	ConnectCommandName    string = "connect"
	DisconnectCommandName string = "disconnect"
	PingCommandName       string = "ping"
)

const (
	ActiveClient string = "Connected"
)

type CommandRequest struct {
	Type   string `json:"tape"`
	Sender string `json:"sender"`
	Value  string `json:"value"`
}

type CommandResponse struct {
	Result bool   `json:"result"`
	Value  string `json:"value"`
}
