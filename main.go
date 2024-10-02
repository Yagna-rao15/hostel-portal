package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB // Global variable for the database connection
var err error

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

func main() {
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}
	log.Println("Database successfully connected!")

	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"email" TEXT NOT NULL UNIQUE,
		"hashed_password" TEXT NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/login", checkEmailHandler).Methods("POST")
	router.HandleFunc("/verify-email", verifyEmailHandler).Methods("POST")
	router.HandleFunc("/test-db", testDBHandler).Methods("GET")

	c := handlers.CORS(handlers.AllowedOrigins([]string{"http://localhost:5173"}))
	srv := http.Server{
		Addr:    ":8080",
		Handler: c(router),
	}

	log.Fatal(srv.ListenAndServe())
}

func checkEmailHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)
	enableCors(&w)

	var requestData struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := requestData.Email
	if err := validateEmail(email); err != nil {
		json.NewEncoder(w).Encode(map[string]bool{
			"valid":        false,
			"isRegistered": false,
		})
		return
	}

	// Check if the user exists in the database
	var hashedPassword string
	err := db.QueryRow("SELECT hashed_password FROM users WHERE email = ?", email).Scan(&hashedPassword)
	isRegistered := err == nil

	// If user doesn't exist, send OTP
	if !isRegistered {
		err = sendOTP(email)
		if err != nil {
			log.Fatal("Error sending OTP:", err)
		}
	}

	json.NewEncoder(w).Encode(map[string]bool{
		"valid":        true,
		"isRegistered": isRegistered,
	})
}

func validateEmail(email string) error {
	regex := `^[a-zA-Z0-9._%+-]+@([a-zA-Z]+)\.svnit\.ac\.in$`
	valid, _ := regexp.MatchString(regex, email)
	if !valid {
		return errors.New("invalid email domain")
	}
	return nil
}

func enableCors(w *http.ResponseWriter) {
	header := (*w).Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	header.Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)
	enableCors(&w)

	var requestData struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	otp := requestData.OTP
	email := requestData.Email
	if err := validateEmail(email); err != nil {
		json.NewEncoder(w).Encode(map[string]bool{
			"valid":        false,
			"isRegistered": false,
			"login":        false,
		})
		return
	}

	if verifyOTP(email, otp) {
		json.NewEncoder(w).Encode(map[string]bool{
			"valid":        true,
			"isRegistered": true,
			"login":        true,
		})
		return
	}
}

func testDBHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database connection is nil", http.StatusInternalServerError)
		return
	}

	var testEmail string
	err := db.QueryRow("SELECT email FROM users LIMIT 1").Scan(&testEmail)
	if err != nil {
		log.Printf("Error reading from database: %v", err)
		http.Error(w, "Database query failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"email": testEmail,
	})
}
