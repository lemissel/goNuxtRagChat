package handlers

import (
	"context"
	"fmt"
	"log"
	"rag-app/backend/database"
	"rag-app/backend/models"
	"rag-app/backend/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/ollama/ollama/api"
	"github.com/pgvector/pgvector-go"
)

type ChatRequest struct {
	Question string `json:"question"`
}

func ChatWithRAG(c *fiber.Ctx) error {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Não foi possível parsear a requisição"})
	}

	questionEmbedding, err := utils.GetEmbedding(req.Question)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Não foi possível gerar o embedding da pergunta"})
	}

	var similarActivities []models.Activity
	// Converte []float64 para []float32 para pgvector.NewVector
	questionEmbedding32 := make([]float32, len(questionEmbedding))
	for i, v := range questionEmbedding {
		questionEmbedding32[i] = float32(v)
	}
	vector := pgvector.NewVector(questionEmbedding32)

	// A busca de similaridade continua a mesma
	database.DB.Order(fmt.Sprintf("embedding <-> '%v'", vector)).Limit(5).Find(&similarActivities)

	contextText := ""
	for _, activity := range similarActivities {
		contextText += " - " + activity.Description + "\n"
	}

	// Criamos um prompt no formato de sistema para dar contexto
	// systemPrompt := fmt.Sprintf("Responda em português do Brasil. Use as seguintes atividades para responder à pergunta. Não use conhecimento externo. Se a resposta não puder ser encontrada, diga 'Não consigo responder com base nas atividades fornecidas.'\n\nAtividades:\n%s\n\n", contextText)

	// llmClient, err := api.ClientFromEnvironment()
	// if err != nil {
	// 	log.Fatalf("Falha ao criar cliente Ollama: %v", err)
	// }

	// var llmRes string
	// stream := false

	// // Agora, usamos o método Chat para interagir com o modelo
	// err = llmClient.Chat(context.Background(), &api.ChatRequest{
	// 	Model: "llama2",
	// 	Messages: []api.Message{
	// 		{Role: "system", Content: systemPrompt},
	// 		{Role: "user", Content: req.Question},
	// 	},
	// 	Stream: &stream,
	// }, func(resp api.ChatResponse) error {
	// 	llmRes = resp.Message.Content
	// 	return nil
	// })

	// if err != nil {
	// 	log.Printf("Erro na requisição Chat: %v", err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao gerar resposta do LLM"})
	// }

	// return c.JSON(fiber.Map{"answer": llmRes})
	// systemPrompt := fmt.Sprintf("Você é um assistente de busca e recuperação de dados. Responda em Português. Use **apenas** o contexto das atividades fornecidas para responder à pergunta do usuário. Não use conhecimento externo. Responda de forma concisa e direta. Se a resposta não puder ser encontrada **somente** com base nas atividades, responda **exatamente** com a frase: 'Não consigo responder com base nas atividades fornecidas.'\n\nAtividades:\n%s", contextText)
	systemPrompt := fmt.Sprintf("Responda em português. Use as seguintes atividades para responder à pergunta. Não use conhecimento externo. Se nenhuma resposta puder ser encontrada, diga 'Não consigo responder com base nas atividades fornecidas.'\n\nAtividades:\n%s\n\n", contextText)

	llmClient, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Falha ao criar cliente Ollama: %v", err)
	}

	var llmRes string
	stream := false

	// O corpo da requisição agora é mais explícito
	err = llmClient.Chat(context.Background(), &api.ChatRequest{
		Model: "llama2",
		Messages: []api.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: req.Question},
		},
		Stream: &stream,
	}, func(resp api.ChatResponse) error {
		llmRes = resp.Message.Content
		return nil
	})

	if err != nil {
		log.Printf("Erro na requisição Chat: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao gerar resposta do LLM"})
	}

	return c.JSON(fiber.Map{"answer": llmRes})
}
