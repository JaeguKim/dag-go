package config

type (
	RunningStatus int

	PrintStatus struct {
		RStatus RunningStatus
		NodeId  string
	}
)

// 생성 시키면 0, 이미 존재하면 1, 에러면 2
type CreateEdgeErrorType int

const (
	Create CreateEdgeErrorType = iota
	Exist
	Fault
)

// 각 노드에서 runner 를 실행시킬때 나태니는 status
const (
	PreFlight RunningStatus = iota
	PreFlightEnd
	PreFlightFailed
	InFlight
	InFlightEnd
	InFlightFailed
	PostFlight
	PostFlightFailed
	PostFlightEnd
	FlightEnd
)

const (
	Nodes = "nodes"
	Node = "node"
	Id   = "id"
	From    = "from"
	To      = "to"
	Command = "command"
)

const (
	StartNode = "start_node"
	EndNode   = "end_node"
)
