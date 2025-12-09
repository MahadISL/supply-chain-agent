package config

import (
	"log"
	"os"
)

type Config struct {
	Port           string
	PineconeAPIKey string
	PineconeHost   string
	HFToken        string
}

func LoadConfig() *Config {
	apiKey := os.Getenv("PINECONE_API_KEY")
	host := os.Getenv("PINECONE_INDEX_HOST")
	hfToken := os.Getenv("HUGGINGFACE_TOKEN")

	if apiKey == "" || host == "" || hfToken == "" {
		log.Println("Warning: Missing API Keys in Environment Variables")
	}

	port := os.Getenv("INGESTION_PORT")
	if port == "" {
		port = "8081"
	}

	return &Config{
		Port:           port,
		PineconeAPIKey: apiKey,
		PineconeHost:   host,
		HFToken:        hfToken,
	}
}
