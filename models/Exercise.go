package models

import (
	"time"

	"gorm.io/gorm"
)

// Exercise representa un ejercicio en la base de datos
type Exercise struct {
	gorm.Model
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Name      string         `gorm:"not null" json:"name"`
	Sets      []Set          `gorm:"foreignKey:ExerciseID" json:"sets"`
	Routines  []Routine      `gorm:"many2many:routine_work_exercise;" json:"routines"`
}
