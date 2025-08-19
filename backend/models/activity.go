package models

import (
	"time"

	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

// Activity representa a estrutura de uma atividade
type Activity struct {
	gorm.Model                  // Isso já inclui o campo ID (uint)
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Embedding   pgvector.Vector `gorm:"type:vector(1536)" json:"embedding"` // Ajuste o tamanho conforme necessário
	Date        time.Time       `json:"date"`
}
