package entity

import (
	"testing"
	"time"
)

func TestNewStorageValid(t *testing.T) {
	storage := NewStorage(10, time.Now())
	err := storage.Validate()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestNewStorageInvalidValue(t *testing.T) {
	storage := NewStorage(-5, time.Now())
	err := storage.Validate()
	if err == nil {
		t.Error("Expected an error for invalid value, got nil")
	}
	if err.Error() != "value invalid" {
		t.Errorf("Expected error 'value invalid', got: %v", err)
	}
}

func TestToEntityValid(t *testing.T) {
	lastSync := time.Now()
	storage := ToEntity(15, lastSync)
	if storage.Value != 15 {
		t.Errorf("Expected Value to be 15, got: %d", storage.Value)
	}
	if !storage.LastSync.Equal(lastSync) {
		t.Errorf("Expected LastSync to be %v, got: %v", lastSync, storage.LastSync)
	}
}

func TestStorageValidateNegativeValue(t *testing.T) {
	storage := &Storage{
		Value:    -10,
		LastSync: time.Now(),
	}
	err := storage.Validate()
	if err == nil {
		t.Error("Expected an error, got nil")
	}
	if err.Error() != "value invalid" {
		t.Errorf("Expected error 'value invalid', got: %v", err)
	}
}

func TestStorageValidateZeroValue(t *testing.T) {
	storage := &Storage{
		Value:    0,
		LastSync: time.Now(),
	}
	err := storage.Validate()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}
