package models

import (
	"time"

	"gorm.io/gorm"
)

// Routine representa una rutina en la base de datos
type Routine struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	Name        string         `gorm:"unique;not_null" json:"name"`
	Description string         `json:"description"`
	Public      bool           `gorm:"default:true" json:"public"`
	Users       []User         `gorm:"many2many:user_make_routine;" json:"users"`
	Exercises   []Exercise     `gorm:"many2many:routine_work_exercise;" json:"exercises"`
}
