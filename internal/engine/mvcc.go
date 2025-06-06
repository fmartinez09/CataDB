package engine

import (
	"errors"
	"sync"
)

// Versión de un valor asociado a un blob
type VersionedValue struct {
	BlobID  string
	Value   []byte
	WriteTS uint64
	Deleted bool
}

// Almacén con control de versiones por clave
type MVCCStore struct {
	mu       sync.RWMutex
	versions map[string][]VersionedValue // blob-id → lista de versiones
}

// Constructor
func NewMVCCStore() *MVCCStore {
	return &MVCCStore{
		versions: make(map[string][]VersionedValue),
	}
}

// Escribe una nueva versión
func (m *MVCCStore) Write(blobID string, value []byte, writeTS uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	vv := VersionedValue{
		BlobID:  blobID,
		Value:   value,
		WriteTS: writeTS,
		Deleted: false,
	}

	m.versions[blobID] = append(m.versions[blobID], vv)
}

// Lee la última versión válida antes del snapshotTS
func (m *MVCCStore) Read(blobID string, snapshotTS uint64) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	history := m.versions[blobID]
	var result []byte

	for _, version := range history {
		if version.WriteTS <= snapshotTS && !version.Deleted {
			result = version.Value
		}
	}

	if result == nil {
		return nil, errors.New("no version found")
	}

	return result, nil
}
