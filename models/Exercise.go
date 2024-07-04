package models

import "gorm.io/gorm"

// Exercise representa un ejercicio en la base de datos
type Exercise struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string
	Sets        []Set     `gorm:"foreignKey:ExerciseID"`
	Routines    []Routine `gorm:"many2many:routine_work_exercise;"`
}
