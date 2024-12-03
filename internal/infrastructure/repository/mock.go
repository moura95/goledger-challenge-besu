package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/moura95/goledger-challenge-besu/internal/domain/entity"
)

type StorageRepositoryMock struct {
	mu       sync.Mutex
	storage  *entity.Storage
	lastSync time.Time
}

func NewStorageRepositoryMock() *StorageRepositoryMock {
	return &StorageRepositoryMock{
		storage: &entity.Storage{
			Value:    0,
			LastSync: time.Now(),
		},
	}
}

func (r *StorageRepositoryMock) Set(storage entity.Storage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage.Value = storage.Value
	r.storage.LastSync = time.Now()
	return nil
}

func (r *StorageRepositoryMock) Get() (*entity.Storage, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.storage == nil {
		return nil, errors.New("no storage found")
	}
	return r.storage, nil
}
