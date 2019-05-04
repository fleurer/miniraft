package raft

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	leveldbutil "github.com/syndtr/goleveldb/leveldb/util"
)

const (
	kCurrentTerm  = "m:current-term"
	kVotedForPeer = "m:vote-for-peer"
	kCommitIndex  = "m:commit-index"
	kLastApplied  = "m:last-applied"
	kLogEntries   = "l:log-entries"
)

type RaftStorage struct {
	db        *leveldb.DB
	keyPrefix string
	path      string
}

func NewRaftStorage(path string, keyPrefix string) (*RaftStorage, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	s := &RaftStorage{
		db:        db,
		path:      path,
		keyPrefix: keyPrefix,
	}
	return s, nil
}

func (s *RaftStorage) Reset() {
	s.PutCommitIndex(0)
	s.PutCurrentTerm(0)
	s.PutLastApplied(0)
}

func (s *RaftStorage) Close() {
	s.db.Close()
}

func (s *RaftStorage) GetCurrentTerm() (uint64, error) {
	return s.dbGetUint64([]byte(kCurrentTerm))
}

func (s *RaftStorage) PutCurrentTerm(v uint64) error {
	return s.dbPutUint64([]byte(kCurrentTerm), v)
}

func (s *RaftStorage) GetCommitIndex() (uint64, error) {
	return s.dbGetUint64([]byte(kCommitIndex))
}

func (s *RaftStorage) PutCommitIndex(v uint64) error {
	return s.dbPutUint64([]byte(kCommitIndex), v)
}

func (s *RaftStorage) GetLastApplied() (uint64, error) {
	return s.dbGetUint64([]byte(kLastApplied))
}

func (s *RaftStorage) PutLastApplied(v uint64) error {
	return s.dbPutUint64([]byte(kLastApplied), v)
}

func (s *RaftStorage) AppendLogEntries(entries []RaftLogEntry) error {
	batch := new(leveldb.Batch)
	for _, le := range entries {
		k := []byte(fmt.Sprintf("%s:%s", kLogEntries, uint64ToBytes(le.Index)))
		v, _ := json.Marshal(le)
		batch.Put(k, v)
	}
	err := s.db.Write(batch, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *RaftStorage) GetLastLogEntry() (*RaftLogEntry, error) {
	prefix := leveldbutil.BytesPrefix([]byte(kLogEntries))
	iter := s.db.NewIterator(prefix, nil)
	defer iter.Release()
	exists := iter.Last()
	if !exists {
		return nil, nil
	}
	buf := iter.Value()
	le := RaftLogEntry{}
	err := json.Unmarshal(buf, &le)
	if err != nil {
		return nil, err
	}
	return &le, nil
}

func (s *RaftStorage) dbGetUint64(k []byte) (uint64, error) {
	key := []byte(fmt.Sprintf("%s:%s", s.keyPrefix, k))
	buf, err := s.db.Get(key, nil)
	if err != nil {
		return 0, err
	}
	n := binary.LittleEndian.Uint64(buf)
	return n, nil
}

func (s *RaftStorage) dbPutUint64(k []byte, v uint64) error {
	key := []byte(fmt.Sprintf("%s:%s", s.keyPrefix, k))
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	return s.db.Put(key, buf, nil)
}

// uint64ToBytes converts an uint64 number to a lexicographically order bytes
func uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}
