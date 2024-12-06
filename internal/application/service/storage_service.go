package service

import (
	"fmt"

	"github.com/moura95/goledger-challenge-besu/config"
	"github.com/moura95/goledger-challenge-besu/internal/domain/entity"
	"github.com/moura95/goledger-challenge-besu/internal/infrastructure/repository"
	"github.com/moura95/goledger-challenge-besu/pkg/blockchainInteractor"
	"go.uber.org/zap"
)

type StorageService struct {
	repository repository.StorageRepository
	config     config.Config
	logger     *zap.SugaredLogger
}

func NewStorageService(repo repository.StorageRepository, cfg config.Config, log *zap.SugaredLogger) *StorageService {
	return &StorageService{
		repository: repo,
		config:     cfg,
		logger:     log,
	}
}

func (s *StorageService) Set(value int32) (string, error) {
	storage := entity.Storage{
		Value: value,
	}
	err := storage.Validate()
	if err != nil {
		s.logger.Error("Failed to validate storage value: ", err)
		return "", fmt.Errorf("failed to validate storage value: %w", err)
	}

	interact, err := blockchainInteractor.NewBlockchainInteractor(
		s.config.NetworkUrl,
		s.config.ContractAddress,
		s.config.PrivateKey,
		s.config.ABI,
	)
	if err != nil {
		s.logger.Error("Failed to create interact: ", err)
		return "", fmt.Errorf("failed to create interact: %w", err)
	}

	defer interact.Close()

	txHash, err := interact.SetValue(uint64(value))
	if err != nil {
		s.logger.Error("Failed to set storage value: ", err)
		return "", fmt.Errorf("failed to set storage value: %w", err)
	}

	return txHash, nil
}

func (s *StorageService) Get() (*entity.Storage, error) {
	storage, err := s.repository.Get()
	if err != nil {
		s.logger.Error("Failed to get storage value: ", err)
		return nil, fmt.Errorf("failed to get storage value: %w", err)
	}
	if storage == nil {
		s.logger.Warn("Storage is empty")
		return nil, fmt.Errorf("no storage data found")
	}
	return storage, nil
}
