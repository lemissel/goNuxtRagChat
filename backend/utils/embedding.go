package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
)

// GetEmbedding gera um embedding para uma determinada string
func GetEmbedding(text string) ([]float64, error) {
	ollamaURL, err := url.Parse("http://localhost:11434")
	if err != nil {
		return nil, fmt.Errorf("não foi possível parsear a URL do Ollama: %w", err)
	}

	client := api.NewClient(ollamaURL, http.DefaultClient)

	// O campo da requisição foi corrigido aqui para 'Input'
	req := &api.EmbedRequest{
		Model: "mxbai-embed-large",
		Input: text,
	}

	resp, err := client.Embed(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("não foi possível gerar o embedding: %w", err)
	}

	// Converte [][]float32 para [][]float64
	embeddingFloat64 := make([][]float64, len(resp.Embeddings))
	for i, v := range resp.Embeddings {
		vec := make([]float64, len(v))
		for j, val := range v {
			vec[j] = float64(val)
		}
		embeddingFloat64[i] = vec
	}

	// Se você espera apenas um embedding, retorne o primeiro vetor
	if len(embeddingFloat64) > 0 {
		return embeddingFloat64[0], nil
	}
	return nil, fmt.Errorf("nenhum embedding retornado")
}
