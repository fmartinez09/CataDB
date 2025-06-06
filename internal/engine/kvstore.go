package engine

import (
	"sync"
)

type KVStore struct {
	mu    sync.RWMutex
	store map[string][]byte
}

func NewKVStore() *KVStore {
	return &KVStore{
		store: make(map[string][]byte),
	}
}

func (k *KVStore) Get(key string) ([]byte, bool) {
	k.mu.RLock()
	defer k.mu.RUnlock()

	val, ok := k.store[key]
	return val, ok
}

func (k *KVStore) Put(key string, value []byte) {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.store[key] = value
}
