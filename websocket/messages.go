package websocket

type Keepalive struct {
	OsType		string	`json:"os_type"`
	IpAddress	string	`json:"ip_address"`
	Hostname	string	`json:"hostname"`
}

type KeepaliveResponse struct {
	Message		string 	`json:"message"`
}
