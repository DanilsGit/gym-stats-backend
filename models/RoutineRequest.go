package models

// Estructura para recibir los datos de la solicitud
type RoutineRequest struct {
	Name            string            `json:"name"`
	Description     string            `json:"description"`
	Public          bool              `json:"public"`
	ExerciseRequest []ExerciseRequest `json:"exercises"`
}

// Estructura para cambiar el nombre de la rutina
type UpdateNameRoutineRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Estructura para cambiar la descripción de la rutina
type UpdateDescriptionRoutineRequest struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
}
