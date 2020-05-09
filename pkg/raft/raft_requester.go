package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type RaftRequester interface {
	SendRequestVoteRequest(p Peer, req *RequestVoteRequest) (*RequestVoteResponse, error)
	SendAppendEntriesRequest(p Peer, req *AppendEntriesRequest) (*AppendEntriesResponse, error)
}

type raftRequester struct {
	logger *Logger
}

func NewRaftRequester(logger *Logger) RaftRequester {
	return &raftRequester{logger: logger}
}

func (rr *raftRequester) SendAppendEntriesRequest(p Peer, request *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	url := fmt.Sprintf("http://%s/_raft/append-entries", p.Addr)
	resp := AppendEntriesResponse{}
	err := rr.post(p, url, request, &resp)
	// rr.logger.Debugf("raft.request.send-append-entries to=%s req=%#v err=%s", p.ID, request, err)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (rr *raftRequester) SendRequestVoteRequest(p Peer, request *RequestVoteRequest) (*RequestVoteResponse, error) {
	url := fmt.Sprintf("http://%s/_raft/request-vote", p.Addr)
	resp := RequestVoteResponse{}
	err := rr.post(p, url, request, &resp)
	rr.logger.Debugf("raft.request.send-request vote to=%s req=%#v err=%s", p.ID, request, err)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (rr *raftRequester) post(p Peer, url string, request interface{}, response interface{}) error {
	buf, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(buf))
	if err != nil {
		return err
	}

	c := http.Client{Timeout: time.Duration(100 * time.Millisecond)}
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, response)
}

type mockRaftRequester struct {
	rafts map[string]*raft
}

func (r *mockRaftRequester) SendRequestVoteRequest(p Peer, req *RequestVoteRequest) (*RequestVoteResponse, error) {
	raft := r.rafts[p.ID]
	ev := newRaftEV(req)
	raft.eventc <- ev
	resp := (<-ev.respc)
	if r, ok := resp.(ServerResponse); ok {
		return nil, fmt.Errorf("bad result: %s", r.Message)
	} else if r, ok := resp.(RequestVoteResponse); ok {
		return &r, nil
	} else {
		return nil, fmt.Errorf("unknown type")
	}
}

func (r *mockRaftRequester) SendAppendEntriesRequest(p Peer, req *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	raft := r.rafts[p.ID]
	ev := newRaftEV(req)
	raft.eventc <- ev
	resp := <-ev.respc
	if r, ok := resp.(ServerResponse); ok {
		return nil, fmt.Errorf("bad result: %s", r.Message)
	} else if r, ok := resp.(AppendEntriesResponse); ok {
		return &r, nil
	} else {
		return nil, fmt.Errorf("unknown type")
	}
}
