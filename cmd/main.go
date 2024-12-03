package main

import (
	"fmt"
	"log"

	"github.com/moura95/goledger-challenge-besu/config"
	server "github.com/moura95/goledger-challenge-besu/internal"
	"github.com/moura95/goledger-challenge-besu/scripts/db"
	"go.uber.org/zap"
)

func main() {
	// Configs
	loadConfig, _ := config.LoadConfig(".")

	// instance Db

	conn, err := db.ConnectPostgres(loadConfig.DBSource)
	store := conn.DB()
	if err != nil {
		fmt.Println("Failed to Connected Database")
		panic(err)
	}
	log.Print("connection is repository establish")

	// Zap Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Run Gin
	server.RunGinServer(loadConfig, store, sugar)
}
