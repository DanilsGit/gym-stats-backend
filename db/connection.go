package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnection() {

	// Carga las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error cargando el archivo .env")
	}

	// Construye el DSN usando las variables de entorno
	DSN := "host=" + os.Getenv("POSTGRES_HOST") +
		" user=" + os.Getenv("POSTGRES_USER") +
		" password=" + os.Getenv("POSTGRES_PASSWORD") +
		" dbname=" + os.Getenv("POSTGRES_DATABASE") +
		" port=5432"

	var error error
	DB, error = gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if error != nil {
		log.Fatal("Error connecting to database")
	} else {
		log.Println("Connected to database")
	}
}
