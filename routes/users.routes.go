package routes

import (
	"encoding/json"
	"net/http"

	"github.com/danilsgit/gym-stats-backend/db"
	"github.com/danilsgit/gym-stats-backend/models"
)

// Usuario

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

// Editar username del usuario
func PutUserInUsernameHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de la solicitud
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "No se encontró el ID del usuario en la solicitud", http.StatusInternalServerError)
		return
	}

	// Obtener el nuevo username de los parámetros de la solicitud
	var updateInfo struct {
		Username string `json:"username"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateInfo)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Error si el campo está vacío
	if updateInfo.Username == "" {
		http.Error(w, "El campo username es obligatorio", http.StatusBadRequest)
		return
	}

	// Error si otro user tiene el mismo username
	var existingUser models.User
	result := db.DB.Where("username = ?", updateInfo.Username).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "El nombre de usuario ya está en uso", http.StatusBadRequest)
		return
	}

	// Buscar el usuario por ID
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Actualizar el username del usuario
	if err := db.DB.Model(&user).Update("username", updateInfo.Username).Error; err != nil {
		http.Error(w, "Error al actualizar el usuario", http.StatusInternalServerError)
		return
	}

	// Enviar respuesta de éxito
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Usuario actualizado con éxito")
}
