package raft

import "github.com/Fleurer/miniraft/pkg/storage"

const (
	SUCCESS        = 200
	NOT_FOUND      = 404
	BAD_REQUEST    = 400
	INTERNAL_ERROR = 500
)

const (
	kNop    = "nop"
	kPut    = "put"
	kGet    = "get"
	kDelete = "delete"
)

type AppendEntriesMessage struct {
	Term         uint64 `json:"term"`
	LeaderID     string `json:"leaderID"`
	CommitIndex  uint64 `json:"commitIndex"`
	PrevLogIndex uint64 `json:"prevLogIndex"`
	PrevLogTerm  uint64 `json:"prevLogTerm"`

	LogEntries []storage.RaftLogEntry `json:"logEntries,omitempty"`
}

type AppendEntriesReply struct {
	Term         uint64 `json:"term"`
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	LastLogIndex uint64 `json:"last_log_index"`
}

type RequestVoteMessage struct {
	Term         uint64 `json:"term"`
	CandidateID  string `json:"candidateID"`
	LastLogIndex uint64 `json:"lastLogIndex"`
	LastLogTerm  uint64 `json:"lasstLogTerm"`
}

type RequestVoteReply struct {
	Term        uint64 `json:"term"`
	VoteGranted bool   `json:"voteGranted"`
	Message     string `json:"message"`
}

type ShowStatusMessage struct {
}

type ShowStatusReply struct {
	Term        uint64               `json:"term"`
	CommitIndex uint64               `json:"commitIndex"`
	Peers       map[string]Peer      `json:"peers"`
	State       string               `json:"state"`
}

type ServerReply struct {
	Code    int    `code:"code"`
	Message string `json:"message"`
}

type CommandMessage struct {
	Command storage.RaftCommand `json:"command"`
}

type CommandReply struct {
	Message string `json:"message"`
	Value   []byte `json:"value,omitempty"`
}

func newRequestVoteReply(success bool, term uint64, message string) *RequestVoteReply {
	return &RequestVoteReply{VoteGranted: success, Term: term, Message: message}
}

func newAppendEntriesReply(success bool, term uint64, lastLogIndex uint64, message string) *AppendEntriesReply {
	return &AppendEntriesReply{Success: success, Term: term, LastLogIndex: lastLogIndex, Message: message}
}

func newServerReply(code int, message string) *ServerReply {
	return &ServerReply{Code: code, Message: message}
}
