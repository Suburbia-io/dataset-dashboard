package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/Suburbia-io/dashboard/pkg/application"
	"github.com/Suburbia-io/dashboard/pkg/database"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s DATASET_ID", os.Args[1])
	}

	datasetID := os.Args[1]

	configPath := os.Getenv("SUBURBIA_DASHBOARD_CONFIG")
	if configPath == "" {
		log.Fatalf("Environment variable SUBURBIA_DASHBOARD_CONFIG isn't set.")
	}

	config := application.Config{}
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatalf("Failed to open config file: %s", err)
	}

	db, err := database.Bootstrap(config.DB, false)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	db.UpdateConsensusForDatasetCmd(datasetID)

	log.Printf("Done!")
}
