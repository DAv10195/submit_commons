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
	Labels			map[string]interface{}	`json:"labels"`
}

type TaskResponse struct {
	Payload		string						`json:"payload"`
	Handler		string						`json:"handler"`
	Task		string						`json:"task"`
	Status		int							`json:"status"`
	Labels		map[string]interface{}		`json:"labels"`
}

type TaskResponses struct {
	Responses []*TaskResponse	`json:"task_responses"`
}

type MossPair struct {
	Percentage1 	int 		`json:"percentage1"`
	Percentage2		int			`json:"percentage2"`
	Name1			string		`json:"name1"`
	Name2			string		`json:"name2"`
}

type MossOutput struct {
	Pairs	[]*MossPair		`json:"matches"`
	Link	string			`json:"link"`
}
