package routes

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/danilsgit/gym-stats-backend/db"
	"github.com/danilsgit/gym-stats-backend/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Rutinas
func GetRoutinesHandler(w http.ResponseWriter, r *http.Request) {
	var routines []models.Routine
	db.DB.Find(&routines)
	json.NewEncoder(w).Encode(&routines)
}

// Rutinas del usuario
func GetUserRoutinesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "No se encontró el ID del usuario en la solicitud", http.StatusInternalServerError)
		return
	}

	// Buscar el usuario por ID y obtener sus rutinas
	var user models.User
	if err := db.DB.Preload("Routines.Exercises.Sets").First(&user, "id = ?", userID).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Devolver las rutinas del usuario como respuesta
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.Routines)
}

func CreateUserRoutineHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de la solicitud
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "No se encontró el ID del usuario en la solicitud", http.StatusInternalServerError)
		return
	}

	// Buscar el usuario por ID
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Decodificar la solicitud en una estructura RoutineRequest
	var req models.RoutineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Mensaje de error si la solicitud no se puede decodificar
		http.Error(w, "COD400"+err.Error(), http.StatusBadRequest)
		return
	}

	// Crear la rutina y sus relaciones con ejercicios y sets
	routine := models.Routine{Name: req.Name, Description: req.Description}
	// Recorrer los ejercicios de la solicitud y crearlos
	// for _ significa que no nos importa el índice, solo el valor
	// Se crea exReq en cada iteración
	for _, exReq := range req.ExerciseRequest {
		exercise := models.Exercise{Name: exReq.Name}
		for _, setReq := range exReq.Sets {
			set := models.Set{Reps: setReq.Reps, Weight: setReq.Weight, Rest: setReq.Rest, Note: setReq.Note}
			// El append agrega un elemento al final de un slice
			exercise.Sets = append(exercise.Sets, set)
		}
		routine.Exercises = append(routine.Exercises, exercise)
	}

	// Comprobar que el nombre sea único
	// si no lo es asignarle un nuevo nombre con los últimos 4 dígitos random hasta que sea único
	var existingRoutine models.Routine
	if err := db.DB.Unscoped().First(&existingRoutine, "name = ?", routine.Name).Error; err == nil {
		unique := false
		count := 0
		for !unique {
			routine.Name = routine.Name + generateFourRandomDigits()
			var stillExistingRoutine models.Routine
			if err := db.DB.Unscoped().First(&stillExistingRoutine, "name = ?", routine.Name).Error; err != nil {
				unique = true
			}
			count = count + 1
			if count > 1000 {
				http.Error(w, "Error: Intenta cambiar el nombre de la rutina", http.StatusBadRequest)
				return
			}
		}
	}

	// Guardar la rutina en la base de datos
	if err := db.DB.Create(&routine).Error; err != nil {
		http.Error(w, "ERROR AL CREAR LA RUTINA"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Asociar la rutina al usuario
	if err := db.DB.Model(&user).Association("Routines").Append(&routine); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(routine)
}

func UpdateNameUserRoutineHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "No se encontró el ID del usuario en la solicitud", http.StatusInternalServerError)
		return
	}

	// Buscar el usuario por ID
	var user models.User
	if err := db.DB.Preload("Routines.Exercises.Sets").First(&user, "id = ?", userID).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	// Decodificar la solicitud en un UpdateNameRoutineRequest
	var req models.UpdateNameRoutineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "COD400"+err.Error(), http.StatusBadRequest)
		return
	}

	// Buscar la rutina por ID
	var routine models.Routine
	if err := db.DB.First(&routine, "id = ?", req.ID).Error; err != nil {
		http.Error(w, "Rutina no encontrada", http.StatusNotFound)
		return
	}

	// Comprobar que el nombre sea único
	// si no lo es asignarle un nuevo nombre con los últimos 4 dígitos random hasta que sea único
	var existingRoutine models.Routine
	if err := db.DB.First(&existingRoutine, "name = ?", req.Name).Error; err == nil {
		unique := false
		count := 0
		for !unique {
			req.Name = req.Name + generateFourRandomDigits()
			var stillExistingRoutine models.Routine
			if err := db.DB.First(&stillExistingRoutine, "name = ?", req.Name).Error; err != nil {
				unique = true
			}
			count = count + 1
			if count > 1000 {
				http.Error(w, "Error: Intenta cambiar el nombre de la rutina", http.StatusBadRequest)
				return
			}
		}
	}

	// Actualizar el nombre de la rutina
	routine.Name = req.Name
	if err := db.DB.Save(&routine).Error; err != nil {
		http.Error(w, "ERR"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(routine)
}

func DeleteUserRoutineHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de la solicitud
	userID := r.Context().Value("userID")
	if userID == nil {
		http.Error(w, "No se encontró el ID del usuario en la solicitud", http.StatusInternalServerError)
		return
	}

	// Obtener el id de la rutina de los parámetros de la solicitud
	params := mux.Vars(r)
	routineId := params["id"]

	// Buscar la rutina por ID
	var routine models.Routine
	if err := db.DB.Preload("Exercises.Sets").First(&routine, "id = ?", routineId).Error; err != nil {
		http.Error(w, "Rutina no encontrada", http.StatusNotFound)
		return
	}

	// Antes de eliminar la rutina, cambiarle el nombre para que no sea igual a otro
	// Con un uuid de la librería
	uuidToAdd := uuid.NewString()
	routine.Name = routine.Name + uuidToAdd

	// Guardar la rutina con el nuevo nombre
	if err := db.DB.Save(&routine).Error; err != nil {
		http.Error(w, "Error al eliminar con nombre uuid", http.StatusInternalServerError)
		return
	}

	// Eliminar cada Exercise y sus Sets asociados
	for _, exercise := range routine.Exercises {
		for _, set := range exercise.Sets {
			if err := db.DB.Delete(&set).Error; err != nil {
				http.Error(w, "Error al eliminar los sets", http.StatusInternalServerError)
				return
			}
		}
		if err := db.DB.Delete(&exercise).Error; err != nil {
			http.Error(w, "Error al eliminar los ejercicios", http.StatusInternalServerError)
			return
		}
	}

	// Eliminar la rutina
	if err := db.DB.Delete(&routine).Error; err != nil {
		http.Error(w, "Error al eliminar la rutina", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func generateFourRandomDigits() string {
	// Crea un generador de números aleatorios local
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Genera un número aleatorio entre 0 y 9999 usando el generador local
	number := r.Intn(10000)

	// Convierte el número a una cadena
	strNumber := strconv.Itoa(number)

	// Añade ceros al principio si es necesario para asegurar 4 dígitos
	for len(strNumber) < 4 {
		strNumber = "0" + strNumber
	}

	return strNumber
}
