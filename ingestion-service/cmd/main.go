package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ingestive-service/config"
	"ingestive-service/internal/processor"
	"io"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()
	pdfProc := processor.NewPDFProcessor()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ingestion Service Status: ONLINE")
	})

	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// A. Parse Multipart Form (10MB limit)
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("Received file: %s", header.Filename)

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		fileReader := bytes.NewReader(fileBytes)

		chunks, err := pdfProc.Process(fileReader, int64(len(fileBytes)))
		if err != nil {
			log.Printf("Processing failed: %v", err)
			http.Error(w, "Failed to process PDF", http.StatusInternalServerError)
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"filename": header.Filename,
			"chunks":   len(chunks),
			"data":     chunks, // Returning data so we can verify it works
		})
	})

	log.Printf("Ingestion Service (Go) starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
