package main

import (
	"net/http"

	"github.com/danilsgit/gym-stats-backend/db"
	"github.com/danilsgit/gym-stats-backend/models"
	"github.com/danilsgit/gym-stats-backend/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	db.DBConnection()

	db.DB.AutoMigrate(models.User{})
	db.DB.AutoMigrate(models.Routine{})
	db.DB.AutoMigrate(models.Exercise{})
	db.DB.AutoMigrate(models.Set{})

	r := mux.NewRouter()

	r.HandleFunc("/", routes.HomeHandler)

	// Usuario y autenticación
	r.HandleFunc("/users", routes.PostUserHandler).Methods("POST")
	r.HandleFunc("/login", routes.LoginHandler).Methods("POST")
	r.HandleFunc("/login/social", routes.LoginSocialHandler).Methods("POST")
	// Rutinas generales
	r.HandleFunc("/routines", routes.GetRoutinesHandler).Methods("GET")
	r.Handle("/routines/copy", routes.JwtAuthentication(http.HandlerFunc(routes.CopyRoutineHandler))).Methods("POST")
	r.HandleFunc("/routines/{id}", routes.GetRoutineHandler).Methods("GET")
	// Rutinas del usuario
	// Sin autenticación
	r.HandleFunc("/users/routines/{userId}", routes.GetRoutineByUserIdHandler).Methods("GET")
	// Con autenticación
	r.Handle("/users/routines", routes.JwtAuthentication(http.HandlerFunc(routes.GetUserRoutinesHandler))).Methods("GET")
	r.Handle("/users/routines", routes.JwtAuthentication(http.HandlerFunc(routes.CreateUserRoutineHandler))).Methods("POST")
	r.Handle("/users/routines/name", routes.JwtAuthentication(http.HandlerFunc(routes.UpdateNameUserRoutineHandler))).Methods("PUT")
	r.Handle("/users/routines/description", routes.JwtAuthentication(http.HandlerFunc(routes.UpdateDescriptionUserRoutineHandler))).Methods("PUT")
	r.Handle("/users/routines/{id}", routes.JwtAuthentication(http.HandlerFunc(routes.DeleteUserRoutineHandler))).Methods("DELETE")
	// Ejercicios del usuario
	r.Handle("/users/routines/exercises/name", routes.JwtAuthentication(http.HandlerFunc(routes.UpdateNameUserExerciseHandler))).Methods("PUT")
	r.Handle("/users/routines/exercises", routes.JwtAuthentication(http.HandlerFunc(routes.CreateUserExerciseHandler))).Methods("POST")
	r.Handle("/users/routines/exercises/{id}", routes.JwtAuthentication(http.HandlerFunc(routes.DeleteUserExerciseHandler))).Methods("DELETE")
	r.Handle("/users/routines/exercises/sets", routes.JwtAuthentication(http.HandlerFunc(routes.UpdateUserExerciseHandler))).Methods("PUT")
	// Configuración del usuario
	r.Handle("/users/config/username", routes.JwtAuthentication(http.HandlerFunc(routes.PutUserInUsernameHandler))).Methods("PUT")

	corsOpts := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173", "https://gymstats.netlify.app"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	http.ListenAndServe(":8080", corsOpts(r))
}
