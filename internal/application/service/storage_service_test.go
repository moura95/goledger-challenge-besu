package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/moura95/goledger-challenge-besu/internal/domain/entity"
	"github.com/moura95/goledger-challenge-besu/internal/infrastructure/repository"
	"github.com/stretchr/testify/assert"
)

func NewReceiverServiceTest(repo repository.StorageRepository) *StorageService {
	return &StorageService{
		repository: repo,
	}
}

func TestSetValue(t *testing.T) {
	mockRepo := repository.NewStorageRepositoryMock()
	service := NewReceiverServiceTest(mockRepo)

	storageBound := entity.NewStorage(100, time.Now())
	err := storageBound.Validate()
	if err != nil {
		fmt.Println(err)
	}

	err = service.repository.Set(*storageBound)
	if err != nil {
		t.Error("Failed to created")
	}
	assert.NoError(t, err)

}

func TestGetValue(t *testing.T) {
	mockRepo := repository.NewStorageRepositoryMock()
	service := NewReceiverServiceTest(mockRepo)

	storageBound := entity.NewStorage(100, time.Now())
	err := storageBound.Validate()
	if err != nil {
		fmt.Println(err)
	}

	err = service.repository.Set(*storageBound)
	if err != nil {
		t.Error("Failed to created")
	}
	assert.NoError(t, err)

	storage, err := service.repository.Get()
	if err != nil {
		t.Error("Failed to get")
	}
	assert.Equal(t, int32(100), storage.Value)
}
