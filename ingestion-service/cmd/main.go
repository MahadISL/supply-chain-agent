package cmd

import (
	"fmt"
	"ingestive-service/config"
	"log"
	"net/http"
)

func main() {

	cfg := config.LoadConfig()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		log.Printf("Health check received from %s", r.RemoteAddr)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ingestion Service Status: ONLINE")
	})

	log.Printf("Ingestion Service (Go) starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
