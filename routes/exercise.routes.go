package routes

import "net/http"

func GetExercisesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /users"))
}

func GetExerciseHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /users"))
}

func PostExerciseHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /users"))
}

func DeleteExerciseHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GET /users"))
}
