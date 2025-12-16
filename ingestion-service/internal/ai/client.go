package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const ModelURL = "https://api-inference.huggingface.co/models/BAAI/bge-small-en-v1.5"

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{Token: token}
}

// Payload structure to match HF API requirements strictly
type HFPayload struct {
	Inputs []string `json:"inputs"`
}

func (c *Client) GenerateEmbeddings(texts []string) ([][]float32, error) {
	// Wrap inputs in a struct
	payload := HFPayload{Inputs: texts}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", ModelURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HuggingFace API Error: %s", resp.Status)
	}

	// Parse Response (List of List of Floats)
	var embeddings [][]float32
	if err := json.NewDecoder(resp.Body).Decode(&embeddings); err != nil {
		return nil, err
	}

	return embeddings, nil
}
