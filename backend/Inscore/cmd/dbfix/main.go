package main

import (
	"context"
	"log"
	"os"

	"github.com/newage-saint/insuretech/backend/inscore/db"
	"github.com/newage-saint/insuretech/backend/inscore/pkg/logger"
	"github.com/newage-saint/insuretech/ops/env"
)

func main() {
	_ = logger.Initialize(logger.NoFileConfig())
	_ = env.Load()
	
	configPath := os.Getenv("INSCORE_DB_CONFIG")
	if configPath == "" {
		configPath = "../database.yaml"
	}
	
	if err := db.InitializeManagerForService(configPath); err != nil {
		log.Fatalf("Failed to initialize db: %v", err)
	}
	
	conn := db.GetDB()
	if conn == nil {
		log.Fatal("DB connection is nil")
	}
	
	err := conn.WithContext(context.Background()).Exec("ALTER TABLE insurance_schema.fraud_cases ALTER COLUMN outcome TYPE VARCHAR(50);").Error
	if err != nil {
		log.Fatalf("Failed to execute alter table: %v", err)
	}
	log.Println("Successfully altered table")
}
