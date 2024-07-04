package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danilsgit/gym-stats-backend/db"
	"github.com/danilsgit/gym-stats-backend/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	// Carga las variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	jwtKey = []byte(os.Getenv("JWT_KEY"))
}

var jwtKey []byte // Actualizado para ser inicializado en init()

func GenerateJWT(userID string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // El token expira en 24 horas

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Error al decodificar las credenciales", http.StatusBadRequest)
		return
	}

	// Validar datos vacíos
	if credentials.Email == "" || credentials.Password == "" {
		http.Error(w, "Todos los campos son obligatorios", http.StatusBadRequest)
		return
	}

	// Buscar usuario en la base de datos
	var userExist models.User
	result := db.DB.Where("email = ?", credentials.Email).First(&userExist)
	if result.Error != nil {
		http.Error(w, "Usuario sin cuenta", http.StatusNotFound)
		return
	}

	// Comparar contraseñas
	errPassword := bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(credentials.Password))
	if errPassword != nil {
		// Si hay un error, la comparación falló, lo que significa que las contraseñas no coinciden
		http.Error(w, "Contraseña incorrecta", http.StatusUnauthorized)
		return
	}

	tokenString, err := GenerateJWT(userExist.ID)
	if err != nil {
		http.Error(w, "Error al generar el token JWT", http.StatusInternalServerError)
		return
	}

	// Información del usuario
	userInfo := map[string]interface{}{
		"id":       userExist.ID,
		"username": userExist.Username,
		"email":    userExist.Email,
		"role":     userExist.Role,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user":  userInfo,
	})
}

func LoginSocialHandler(w http.ResponseWriter, r *http.Request) {
	var credentials models.User
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Error al decodificar el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}

	// Buscar usuario en la base de datos
	var userExist models.User
	result := db.DB.Where("email = ?", credentials.Email).First(&userExist)
	if result.Error != nil {
		// Crear usuario si no existe, primero valida si el username ya existe, si existe se le concatenará la id
		resultUsername := db.DB.Where("username = ?", credentials.Username).First(&userExist)
		if resultUsername.Error == nil {
			// Agarrar los últimos 4 caracteres del ID para concatenarlos al username
			credentials.Username = credentials.Username + "-" + credentials.ID[len(credentials.ID)-4:]
		}
		createdUser := db.DB.Create(&credentials)
		if createdUser.Error != nil {
			http.Error(w, "Error al crear el usuario", http.StatusInternalServerError)
			return
		}
		userExist = credentials
	}

	// Si tiene contraseña, quiere decir que no se ha registrado con un proveedor social
	if userExist.Password != "" {
		http.Error(w, "Utiliza el inicio de sesión tradicional", http.StatusBadRequest)
		return
	}

	tokenString, err := GenerateJWT(userExist.ID)
	if err != nil {
		http.Error(w, "Error al generar el token JWT", http.StatusInternalServerError)
		return
	}

	// Información del usuario
	userInfo := map[string]interface{}{
		"id":       userExist.ID,
		"username": userExist.Username,
		"email":    userExist.Email,
		"role":     userExist.Role,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": tokenString,
		"user":  userInfo,
	})
}
