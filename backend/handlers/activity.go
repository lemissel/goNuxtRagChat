package handlers

import (
	"log"
	"rag-app/backend/database"
	"rag-app/backend/models"
	"rag-app/backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pgvector/pgvector-go"
)

// CreateActivity lida com a criação de uma nova atividade
func CreateActivity(c *fiber.Ctx) error {
	var activity models.Activity

	if err := c.BodyParser(&activity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	activity.Date = time.Now()
	log.Printf("Atividade recebida: %+v", activity)

	// Obtenha o embedding do texto
	embeddingFloat64, err := utils.GetEmbedding(activity.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Converte []float64 para []float32
	embeddingFloat32 := make([]float32, len(embeddingFloat64))
	for i, v := range embeddingFloat64 {
		embeddingFloat32[i] = float32(v)
	}

	// Converte o slice de float32 para o tipo pgvector.Vector
	activity.Embedding = pgvector.NewVector(embeddingFloat32)

	// Salva a atividade no banco de dados
	if err := database.DB.Create(&activity).Error; err != nil {
		log.Printf("ERROR: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Falha ao criar atividade no banco de dados"})
	}

	return c.Status(fiber.StatusCreated).JSON(activity)
}

func GetActivities(c *fiber.Ctx) error {
	var activities []models.Activity
	database.DB.Order("created_at desc").Find(&activities)
	return c.JSON(activities)
}
