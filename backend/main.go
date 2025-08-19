package main

import (
	"rag-app/backend/database"
	"rag-app/backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.ConnectDB()

	app := fiber.New()
	app.Use(cors.New())

	app.Post("/api/activities", handlers.CreateActivity)
	app.Get("/api/activities", handlers.GetActivities)
	app.Post("/api/chat", handlers.ChatWithRAG)

	app.Listen(":8080")
}
