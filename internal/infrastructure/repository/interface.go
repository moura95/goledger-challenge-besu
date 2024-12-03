package repository

import (
	"github.com/moura95/goledger-challenge-besu/internal/domain/entity"
)

type StorageRepository interface {
	Get() (response *entity.Storage, err error)
	Set(value entity.Storage) error
}
