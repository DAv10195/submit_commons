package websocket

const (
	Agents 							= "agents"

	MessageTypeKeepalive			= "keepalive"
	MessageTypeKeepaliveResponse	= "keepalive_response"

	MessageTypeTask					= "task"
	MessageTypeTaskResponses		= "task_responses"

	AgentIdHeader					= "X-Submit-Agent-ID"
)
