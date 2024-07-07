package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User representa un usuario en la base de datos
type User struct {
	gorm.Model
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	Username  string         `gorm:"unique;not null" json:"username"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"size:100" json:"password"`
	Role      string         `gorm:"default:'user'" json:"role"`
	Routines  []Routine      `gorm:"many2many:user_make_routine;" json:"routines"`
}

// Las tablas intermedias user_make_routine y routine_work_exercise son manejadas automáticamente por GORM gracias a las anotaciones many2many.

// BeforeCreate se llama antes de crear un usuario en la base de datos
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	// Si el ID está vacío, generar un nuevo UUID
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	// Hashear la contraseña si no está vacía
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}
	return nil
}
