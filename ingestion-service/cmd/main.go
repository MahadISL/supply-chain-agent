package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ingestive-service/config"
	"ingestive-service/internal/ai"
	"ingestive-service/internal/processor"
	"ingestive-service/internal/store"
	"io"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	pdfProc := processor.NewPDFProcessor()
	aiClient := ai.NewClient(cfg.HFToken)
	dbClient := store.NewPineconeClient(cfg.PineconeAPIKey, cfg.PineconeHost)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ingestion Service Status: ONLINE")
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Receive File
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("Processing file: %s", header.Filename)

		// Extract Text
		fileBytes, _ := io.ReadAll(file)
		chunks, err := pdfProc.Process(bytes.NewReader(fileBytes), int64(len(fileBytes)))
		if err != nil {
			http.Error(w, "PDF Processing failed", http.StatusInternalServerError)
			return
		}

		// Generate Embeddings
		var texts []string
		for _, c := range chunks {
			texts = append(texts, c.Text)
		}

		embeddings, err := aiClient.GenerateEmbeddings(texts)
		if err != nil {
			log.Printf("Embedding failed: %v", err)
			http.Error(w, "AI Embedding failed", http.StatusInternalServerError)
			return
		}

		// Store in Pinecone
		err = dbClient.Upsert(chunks, embeddings, header.Filename)
		if err != nil {
			log.Printf("Pinecone Upsert failed: %v", err)
			http.Error(w, "Database storage failed", http.StatusInternalServerError)
			return
		}

		log.Printf("Successfully ingested %s with %d chunks", header.Filename, len(chunks))

		// Success Response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":           "success",
			"filename":         header.Filename,
			"chunks_processed": len(chunks),
		})
	})

	log.Printf("Ingestion Service (Go) starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
