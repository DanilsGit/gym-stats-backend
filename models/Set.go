package models

import (
	"time"

	"gorm.io/gorm"
)

// Set representa un set de un ejercicio en la base de datos
type Set struct {
	gorm.Model
	ID         uint           `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Reps       int            `gorm:"not null" json:"reps"`
	Weight     float64        `gorm:"not null" json:"weight"`
	Rest       float64        `gorm:"not null" json:"rest"`
	Note       string         `gorm:"size 255" json:"note"`
	ExerciseID uint           // Llave foránea que referencia a Exercise
}
