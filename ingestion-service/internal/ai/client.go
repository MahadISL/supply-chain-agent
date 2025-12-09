package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const ModelURL = "https://api-inference.huggingface.co/models/sentence-transformers/all-MiniLM-L6-v2"

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{Token: token}
}

func (c *Client) GenerateEmbeddings(texts []string) ([][]float32, error) {

	jsonData, err := json.Marshal(texts)
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
	
	var embeddings [][]float32
	if err := json.NewDecoder(resp.Body).Decode(&embeddings); err != nil {
		return nil, err
	}

	return embeddings, nil
}
