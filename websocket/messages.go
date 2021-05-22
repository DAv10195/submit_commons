package websocket

import "github.com/DAv10195/submit_commons/containers"

const (
	TaskRespExecStatusOk = iota
	TaskRespExecStatusErr = iota
)

type Keepalive struct {
	OsType			string	`json:"os_type"`
	IpAddress		string	`json:"ip_address"`
	Hostname		string	`json:"hostname"`
	Architecture	string	`json:"architecture"`
	NumRunningTasks	int		`json:"num_running_tasks"`
}

type KeepaliveResponse struct {
	Message		string 	`json:"message"`
}

type Task struct {
	ID				string					`json:"id"`
	Command			string					`json:"command"`
	ResponseHandler	string					`json:"response_handler"`
	Timeout			int						`json:"timeout"`
	Dependencies	*containers.StringSet	`json:"dependencies"`
}

type TaskResponse struct {
	Payload		string	`json:"payload"`
	Handler		string	`json:"handler"`
	Task		string	`json:"task"`
	Status		int		`json:"status"`
}
