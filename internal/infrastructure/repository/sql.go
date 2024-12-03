package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/moura95/goledger-challenge-besu/internal/domain/entity"
)

type StorageRepositorySql struct {
	db *sqlx.DB
}

func NewStorageRepository(db *sqlx.DB) StorageRepository {
	return &StorageRepositorySql{db: db}
}

type StorageModel struct {
	Value    int32     `db:"variable_value"`
	LastSync time.Time `db:"last_synced_at"`
}

func (r StorageRepositorySql) Set(storage entity.Storage) error {
	query := `UPDATE simple_storage SET variable_value = $1, last_synced_at = now()`
	args := []interface{}{
		storage.Value,
	}
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r StorageRepositorySql) Get() (*entity.Storage, error) {
	var result StorageModel

	query := `SELECT variable_value, last_synced_at FROM simple_storage limit 1`

	err := r.db.Get(&result, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return entity.NewStorage(result.Value, result.LastSync), nil
}
