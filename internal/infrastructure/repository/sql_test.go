package repository

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/moura95/goledger-challenge-besu/internal/domain/entity"
	"github.com/moura95/goledger-challenge-besu/scripts/db"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	pgContainer testcontainers.Container
	connStr     string
)

func setupPostgresContainer() (func(), error) {
	ctx := context.Background()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("../../..", "scripts/db/migrations", "000001_init_schema.up.sql.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, err
	}

	connStr = "postgres://postgres:postgres@localhost:" + mappedPort.Port() + "/test-db?sslmode=disable"

	return func() {
		if err := container.Terminate(ctx); err != nil {
			fmt.Printf("failed to terminate pgContainer: %s", err)
		}
	}, nil
}

func TestMain(m *testing.M) {
	cleanup, err := setupPostgresContainer()
	if err != nil {
		panic(fmt.Sprintf("Failed to set up PostgreSQL container: %s", err))
	}
	defer cleanup()

	m.Run()
}

func TestStorageRepositorySql_Set(t *testing.T) {
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	repo := NewStorageRepository(conn.DB())

	storage := entity.Storage{
		Value:    42,
		LastSync: time.Now(),
	}

	err = repo.Set(storage)
	assert.Equal(t, nil, err)
}

func TestStorageRepositorySql_Get(t *testing.T) {
	conn, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	repo := NewStorageRepository(conn)

	storage := entity.Storage{
		Value:    100,
		LastSync: time.Now(),
	}
	err = repo.Set(storage)
	assert.NoError(t, err)

	result, err := repo.Get()
	assert.NoError(t, err)
	assert.NotNil(t, result)
	fmt.Println(result)
	assert.Equal(t, storage.Value, result.Value)
}
