package websocket

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
