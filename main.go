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

	r.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	r.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")
	r.HandleFunc("/users", routes.PostUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", routes.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/login", routes.LoginHandler).Methods("POST")
	r.HandleFunc("/login/social", routes.LoginSocialHandler).Methods("POST")

	corsOpts := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173", "https://gymstats.netlify.app"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
	)

	http.ListenAndServe(":8080", corsOpts(r))
}
