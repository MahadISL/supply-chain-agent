package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"ingestive-service/internal/processor"
	"net/http"
)

type PineconeClient struct {
	APIKey string
	Host   string
}

func NewPineconeClient(apiKey, host string) *PineconeClient {
	return &PineconeClient{
		APIKey: apiKey,
		Host:   host,
	}
}

type UpsertRequest struct {
	Vectors []Vector `json:"vectors"`
}

type Vector struct {
	ID       string                 `json:"id"`
	Values   []float32              `json:"values"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (p *PineconeClient) Upsert(chunks []processor.DocumentChunk, embeddings [][]float32, filename string) error {
	if len(chunks) != len(embeddings) {
		return fmt.Errorf("mismatch between chunks and embeddings count")
	}

	var vectors []Vector

	// Convert chunks + embeddings into Pinecone Vectors
	for i, chunk := range chunks {
		vectorID := fmt.Sprintf("%s_chunk_%d", filename, chunk.ChunkIdx)

		vectors = append(vectors, Vector{
			ID:     vectorID,
			Values: embeddings[i],
			Metadata: map[string]interface{}{
				"text":     chunk.Text,
				"filename": filename,
				"page_num": chunk.PageNum,
			},
		})
	}
	
	reqBody, err := json.Marshal(UpsertRequest{Vectors: vectors})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/vectors/upsert", p.Host)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Api-Key", p.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Pinecone API Error: %s", resp.Status)
	}

	return nil
}
