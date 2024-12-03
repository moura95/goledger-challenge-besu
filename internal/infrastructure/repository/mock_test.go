package repository

import (
	"testing"
	"time"

	"github.com/moura95/goledger-challenge-besu/internal/domain/entity"
)

func TestStorageRepositoryMock(t *testing.T) {
	repo := NewStorageRepositoryMock()

	err := repo.Set(entity.Storage{
		Value:    42,
		LastSync: time.Now(),
	})
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	storage, err := repo.Get()
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if storage.Value != 42 {
		t.Errorf("Expected Value 42, got %d", storage.Value)
	}
}
