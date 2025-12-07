package processor

import (
	"fmt"
	"io"
	"strings"

	"github.com/ledongthuc/pdf"
)

type DocumentChunk struct {
	Text     string `json:"text"`
	PageNum  int    `json:"page_num"`
	ChunkIdx int    `json:"chunk_idx"`
}

type PDFProcessor struct {
	ChunkSize int
	Overlap   int
}

func NewPDFProcessor() *PDFProcessor {
	return &PDFProcessor{
		ChunkSize: 1000,
		Overlap:   200,
	}
}

// Process reads a PDF stream, extracts text, and chunks it
func (p *PDFProcessor) Process(file io.ReaderAt, size int64) ([]DocumentChunk, error) {

	reader, err := pdf.NewReader(file, size)
	if err != nil {
		return nil, fmt.Errorf("failed to read PDF: %w", err)
	}

	var chunks []DocumentChunk
	totalContent := ""

	for pageIndex := 0; pageIndex < reader.NumPage(); pageIndex++ {
		page := reader.Page(pageIndex + 1)
		if page.V.IsNull() {
			continue
		}

		text, err := page.GetPlainText(nil)
		if err != nil {
			continue // Skip bad pages, don't crash
		}

		text = cleanText(text)
		totalContent += text + " "

		pageChunks := p.chunkText(text, pageIndex+1)
		chunks = append(chunks, pageChunks...)
	}

	return chunks, nil
}

func cleanText(raw string) string {
	// Replace multiple spaces/newlines with single space
	return strings.Join(strings.Fields(raw), " ")
}

func (p *PDFProcessor) chunkText(text string, pageNum int) []DocumentChunk {
	var chunks []DocumentChunk
	runes := []rune(text) // Convert to runes to handle UTF-8/Emoji correctly
	totalLen := len(runes)

	if totalLen == 0 {
		return chunks
	}

	for i := 0; i < totalLen; i += (p.ChunkSize - p.Overlap) {
		end := i + p.ChunkSize
		if end > totalLen {
			end = totalLen
		}

		chunkContent := string(runes[i:end])

		chunks = append(chunks, DocumentChunk{
			Text:     chunkContent,
			PageNum:  pageNum,
			ChunkIdx: len(chunks),
		})

		if end == totalLen {
			break
		}
	}
	return chunks
}
