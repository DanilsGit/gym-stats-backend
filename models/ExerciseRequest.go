package models

// Estructura para recibir los datos de la solicitud
type ExerciseRequest struct {
	// ID de la rutina a la que pertenece el ejercicio para agregarse
	// (No es necesario para actualizar)
	IDRoutine uint `json:"idRoutine"`
	// ID del ejercicio para actualizar (No es necesario para crear)
	IDExercise uint `json:"idExercise"`
	// Demás datos del ejercicio
	Name string `json:"name"`
	Sets []struct {
		// ID del set para actualizar (No es necesario para crear)
		ID     uint    `json:"id_set"`
		Reps   int     `json:"reps"`
		Weight float64 `json:"weight"`
		Rest   float64 `json:"rest"`
		Note   string  `json:"note"`
	} `json:"sets"`
}

// Estructura para cambiar el nombre del ejercicio
type UpdateNameExerciseRequest struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
