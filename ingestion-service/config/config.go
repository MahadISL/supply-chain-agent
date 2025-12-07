package config

import (
	"log"
	"os"
)

type Config struct {
	Port           string
	PineconeAPIKey string
	PineconeHost   string
}

func LoadConfig() *Config {

	apiKey := os.Getenv("PINECONE_API_KEY")
	if apiKey == "" {
		log.Println("Warning: PINECONE_API_KEY is missing (Ignore if running generic health check)")
	}

	host := os.Getenv("PINECONE_INDEX_HOST")

	port := os.Getenv("INGESTION_PORT")
	if port == "" {
		port = "8081" // Default
	}

	return &Config{
		Port:           port,
		PineconeAPIKey: apiKey,
		PineconeHost:   host,
	}
}
