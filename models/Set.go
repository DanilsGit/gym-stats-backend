package models

import "gorm.io/gorm"

// Set representa un set de un ejercicio en la base de datos
type Set struct {
	gorm.Model
	Reps       int     `gorm:"not null"`
	Weight     float64 `gorm:"not null"`
	Note       string  `gorm:"size 255"`
	ExerciseID uint    // Llave foránea que referencia a Exercise
}
