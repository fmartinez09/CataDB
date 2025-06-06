package tx

import (
	"CataDB/internal/engine"
	"sync"
)

type Transaction struct {
	store     *engine.MVCCStore
	writeSet  map[string][]byte
	snapshot  uint64
	committed bool
	mu        sync.Mutex
}

func NewTransaction(store *engine.MVCCStore, snapshotTs uint64) *Transaction {
	return &Transaction{
		store:    store,
		snapshot: snapshotTs,
		writeSet: make(map[string][]byte),
	}
}

func (t *Transaction) Read(key string) ([]byte, error) {
	return t.store.Read(key, t.snapshot)
}

func (t *Transaction) Write(key string, value []byte) {
	t.writeSet[key] = value
}

func (t *Transaction) Commit(commitTs uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.committed {
		return nil
	}

	for k, v := range t.writeSet {
		t.store.Write(k, v, commitTs)
	}

	t.committed = true
	return nil
}
