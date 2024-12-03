package entity

import (
	"errors"
	"time"
)

type Storage struct {
	Value    int32
	LastSync time.Time
}

func NewStorage(Value int32, LastSync time.Time) *Storage {
	return &Storage{
		Value:    Value,
		LastSync: LastSync,
	}
}

func ToEntity(Value int32, LastSync time.Time) *Storage {
	return &Storage{
		Value:    Value,
		LastSync: LastSync,
	}
}

func (c *Storage) Validate() error {

	if c.Value < 0 {
		return errors.New("value invalid")
	}

	return nil
}
