package database

import (
	"fmt"
	"log"

	"rag-app/backend/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=rag_db port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v", err)
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS vector;")

	err = DB.AutoMigrate(&models.Activity{})
	if err != nil {
		log.Fatalf("Falha ao migrar o banco de dados: %v", err)
	}
	fmt.Println("Conexão com o banco de dados estabelecida com sucesso.")
}
