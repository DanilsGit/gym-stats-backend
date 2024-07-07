package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danilsgit/gym-stats-backend/db"
	"github.com/danilsgit/gym-stats-backend/models"
	"github.com/gorilla/mux"
)

// Ejercicios del usuario

// Crear un ejercicio
func CreateUserExerciseHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de la solicitud
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

	// Decodificar la solicitud en una estructura UpdateNameExerciseRequest
	var req models.ExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Buscar la rutina por ID
	idRoutine := req.IDRoutine
	var routine models.Routine
	if err := db.DB.First(&routine, "id = ?", idRoutine).Error; err != nil {
		http.Error(w, "Rutina no encontrada", http.StatusNotFound)
		return
	}

	// Crear un nuevo ejercicio
	exercise := models.Exercise{Name: req.Name}
	for _, setReq := range req.Sets {
		set := models.Set{
			Reps:   setReq.Reps,
			Weight: setReq.Weight,
			Rest:   setReq.Rest,
			Note:   setReq.Note,
		}
		exercise.Sets = append(exercise.Sets, set)
	}

	// Guardar el ejercicio en la base de datos
	if err := db.DB.Create(&exercise).Error; err != nil {
		http.Error(w, "Error al guardar el ejercicio", http.StatusInternalServerError)
		return
	}

	// Asociar el ejercicio a la rutina
	if err := db.DB.Model(&routine).Association("Exercises").Append(&exercise); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Actualizar un ejercicio
func UpdateUserExerciseHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de la solicitud
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

	// Decodificar la solicitud en una estructura ExerciseRequest
	var req models.ExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Buscar el ejercicio por ID
	var exercise models.Exercise
	if err := db.DB.First(&exercise, "id = ?", req.IDExercise).Error; err != nil {
		http.Error(w, "Ejercicio no encontrado", http.StatusNotFound)
		return
	}

	// Obtener todos los sets actuales del ejercicio
	var currentSets []models.Set
	if err := db.DB.Where("exercise_id = ?", exercise.ID).Find(&currentSets).Error; err != nil {
		http.Error(w, "Error al obtener sets actuales", http.StatusInternalServerError)
		return
	}

	// Crear un slice para almacenar los sets finales
	var finalSets []models.Set

	// Crear un mapa de sets actuales por ID
	currentSetsMap := make(map[uint]models.Set)
	for _, set := range currentSets {
		currentSetsMap[set.ID] = set
	}

	// Procesar sets de la solicitud
	for _, setReq := range req.Sets {
		if setReq.ID != 0 { // Si el set tiene un ID, actualizar
			if set, ok := currentSetsMap[setReq.ID]; ok {
				set.Reps = setReq.Reps
				set.Weight = setReq.Weight
				set.Rest = setReq.Rest
				set.Note = setReq.Note
				db.DB.Save(&set)
				finalSets = append(finalSets, set) // Añadir al slice de sets finales
				delete(currentSetsMap, set.ID)     // Eliminar de mapa para no considerarlo para eliminación
			}
		} else { // Si el set es nuevo, crear
			newSet := models.Set{
				ExerciseID: exercise.ID,
				Reps:       setReq.Reps,
				Weight:     setReq.Weight,
				Rest:       setReq.Rest,
				Note:       setReq.Note,
			}
			db.DB.Create(&newSet)
			finalSets = append(finalSets, newSet) // Añadir al slice de sets finales
		}
	}

	// Eliminar sets que no están en la solicitud
	for id := range currentSetsMap {
		db.DB.Delete(&models.Set{}, id)
	}

	// Devolver el arreglo de sets actualizado
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(finalSets)
}

// Cambiar el nombre del ejercicio
func UpdateNameUserExerciseHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de la solicitud
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

	// Decodificar la solicitud en una estructura UpdateNameExerciseRequest
	var req models.UpdateNameExerciseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	// Buscar el ejercicio por ID
	var exercise models.Exercise
	if err := db.DB.First(&exercise, "id = ?", req.ID).Error; err != nil {
		http.Error(w, "Ejercicio no encontrado", http.StatusNotFound)
		return
	}

	// Cambiar el nombre del ejercicio
	exercise.Name = req.Name
	if err := db.DB.Save(&exercise).Error; err != nil {
		http.Error(w, "Error al guardar el ejercicio", http.StatusInternalServerError)
		return
	}

	// Devolver el ejercicio actualizado
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(exercise)
}

// Eliminar un ejercicio
func DeleteUserExerciseHandler(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario de la solicitud
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

	// Obtener el id del ejercicio de los parámetros
	params := mux.Vars(r)
	exerciseId := params["id"]

	// Buscar el ejercicio por ID
	var exercise models.Exercise
	if err := db.DB.First(&exercise, "id = ?", exerciseId).Error; err != nil {
		http.Error(w, "Ejercicio no encontrado", http.StatusNotFound)
		return
	}

	// Eliminar los sets del ejercicio
	for _, set := range exercise.Sets {
		if err := db.DB.Delete(&set).Error; err != nil {
			http.Error(w, "Error al eliminar los sets", http.StatusInternalServerError)
			return
		}
	}

	// Eliminar el ejercicio
	if err := db.DB.Delete(&exercise).Error; err != nil {
		http.Error(w, "Error al eliminar el ejercicio", http.StatusInternalServerError)
		return
	}
}
