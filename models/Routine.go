package models

import "gorm.io/gorm"

// Routine representa una rutina en la base de datos
type Routine struct {
	gorm.Model
	Name      string     `gorm:"unique;not null"`
	Users     []User     `gorm:"many2many:user_make_routine;"`
	Exercises []Exercise `gorm:"many2many:routine_work_exercise;"`
}
