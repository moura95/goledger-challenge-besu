package service

import (
	"fmt"
	"time"

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
	interact, err := blockchainInteractor.NewBlockchainInteractor(
		s.config.NetworkUrl,
		s.config.ContractAddress,
		s.config.PrivateKey,
	)
	if err != nil {
		s.logger.Error("Failed to create interact: ", err)
		return nil, fmt.Errorf("failed to create interact: %w", err)
	}
	value, err := interact.GetValue()
	intValue, ok := value.(int32)
	if !ok {
		s.logger.Error("Value is not of type int32")
		return nil, fmt.Errorf("value is not of type int32")
	}

	storage := entity.ToEntity(intValue, time.Now())
	return storage, nil
}

func (s *StorageService) Check() (bool, error) {
	interact, err := blockchainInteractor.NewBlockchainInteractor(
		s.config.NetworkUrl,
		s.config.ContractAddress,
		s.config.PrivateKey,
	)
	if err != nil {
		s.logger.Error("Failed to create interact: ", err)
		return false, fmt.Errorf("failed to create interact: %w", err)
	}
	value, err := interact.GetValue()
	intValue, ok := value.(int32)
	if !ok {
		s.logger.Error("Value is not of type int32")
		return false, fmt.Errorf("value is not of type int32")
	}
	valueStorage, err := s.repository.Get()
	if err != nil {
		s.logger.Error("Failed to get storage: ", err)
		return false, fmt.Errorf("failed to get storage: %w", err)
	}
	if !ok {
		s.logger.Error("Value is not of type int32")
		return false, fmt.Errorf("value is not of type int32")
	}
	if valueStorage.Value != intValue {
		return false, nil
	}

	return true, nil
}

func (s *StorageService) Sync() (bool, error) {
	interact, err := blockchainInteractor.NewBlockchainInteractor(
		s.config.NetworkUrl,
		s.config.ContractAddress,
		s.config.PrivateKey,
	)
	if err != nil {
		s.logger.Error("Failed to create interact: ", err)
		return false, fmt.Errorf("failed to create interact: %w", err)
	}
	value, err := interact.GetValue()
	intValue, ok := value.(int32)
	if !ok {
		s.logger.Error("Value is not of type int32")
		return false, fmt.Errorf("value is not of type int32")
	}

	storage := entity.ToEntity(intValue, time.Now())
	err = s.repository.Set(*storage)
	if err != nil {
		s.logger.Error("Failed to set storage: ", err)
		return false, fmt.Errorf("failed to set storage: %w", err)
	}
	return true, nil
}
