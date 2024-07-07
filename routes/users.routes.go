package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danilsgit/gym-stats-backend/db"
	"github.com/danilsgit/gym-stats-backend/models"
)

// Usuario
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /users"))
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /users"))
}

func PostUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	errDecoder := json.NewDecoder(r.Body).Decode(&user)

	if errDecoder != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	// Error si los campos no son válidos / están vacíos
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Comprobar si el correo ya existe en la BD
	var existingUser models.User
	result := db.DB.Where("email = ?", user.Email).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "El correo ya está en uso", http.StatusBadRequest)
		return
	}

	// Comprobar si el username ya existe en la BD
	result = db.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "El nombre de usuario ya está en uso", http.StatusBadRequest)
		return
	}

	createdUser := db.DB.Create(&user)

	err := createdUser.Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	json.NewEncoder(w).Encode(&user)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /users"))
}
