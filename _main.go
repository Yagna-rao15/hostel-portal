package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

func main() {
	db = initializeDatabase()

	router := mux.NewRouter()
	router.HandleFunc("/login", checkEmailHandler).Methods("POST")
	router.HandleFunc("/password", checkPasswordHandler).Methods("POST")
	router.HandleFunc("/verify", verifyEmailHandler).Methods("POST")
	router.HandleFunc("/update-password", updatePasswordHandler).Methods("POST")
	router.HandleFunc("/test-db", testDBHandler).Methods("GET")
	router.HandleFunc("/form", submitComplaint).Methods("POST")

	// CORS configuration
	c := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}),
		handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
	)

	// Create the HTTP server
	srv := http.Server{
		Addr:    ":8080",
		Handler: c(router), // Apply CORS middleware
	}

	log.Fatal(srv.ListenAndServe())
}

func initializeDatabase() *sql.DB {
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
		"hashed_password" TEXT NOT NULL,
		"session_token" TEXT,
		"csrf_token" TEXT
	);`

	anotherSQLTable := `CREATE TABLE IF NOT EXISTS complaints (
    "id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "email" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "room" TEXT NOT NULL,
    "mobile" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "complain_type" TEXT NOT NULL,
    "hostel" TEXT NOT NULL,
    "file_path" TEXT NOT NULL
);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	_, err = db.Exec(anotherSQLTable)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	return db
}

func checkEmailHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
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

func checkPasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := requestData.Email
	password := requestData.Password

	if err := validateEmail(email); err != nil {
		json.NewEncoder(w).Encode(map[string]bool{
			"valid":        false,
			"isRegistered": false,
			"isVerified":   false,
		})
		return
	}

	// Check if the user exists in the database
	var hashedPassword string
	if err := db.QueryRow("SELECT hashed_password FROM users WHERE email = ?", email).Scan(&hashedPassword); err != nil {
		json.NewEncoder(w).Encode(map[string]bool{
			"valid":        false,
			"isRegistered": false,
			"isVerified":   false,
		})
		return
	}

	// Verify the password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Password does not match
		json.NewEncoder(w).Encode(map[string]bool{
			"valid":        false,
			"isRegistered": true,
			"isVerified":   false,
		})
		return
	}

	// If everything is correct, return valid, registered, and verified as true
	json.NewEncoder(w).Encode(map[string]bool{
		"valid":        true,
		"isRegistered": true,
		"isVerified":   true,
	})
}

func updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

	var requestData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := requestData.Email
	password := requestData.Password

	// Hash the password (make sure you have a function to hash passwords)
	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Update the password in the database
	result, err := db.Exec("UPDATE users SET hashed_password = ? WHERE email = ?", hashedPassword, email)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected() // Get the number of rows affected
	if rowsAffected == 0 {
		// If no rows were affected, create a new user
		_, err = db.Exec("INSERT INTO users (email, hashed_password) VALUES (?, ?)", email, hashedPassword)
		if err != nil {
			http.Error(w, "Failed to create new user", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"updated": true,
			"message": "User created successfully.",
		})
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"updated": true,
			"message": "Password updated successfully!",
		})
	}
}

type Complaint struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	Room         string `json:"room"`
	Mobile       string `json:"mobile"`
	Description  string `json:"description"`
	ComplainType string `json:"complainType"`
	Hostel       string `json:"hostel"`
}

func submitComplaint(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

	// Parse the multipart form (for file uploads)
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Extract the JSON fields from the form data
	email := r.FormValue("email")
	name := r.FormValue("name")
	room := r.FormValue("room")
	mobile := r.FormValue("mobile")
	description := r.FormValue("description")
	complainType := r.FormValue("complainType")
	hostel := r.FormValue("hostel")

	// Create a complaint object
	complaint := Complaint{
		Email:        email,
		Name:         name,
		Room:         room,
		Mobile:       mobile,
		Description:  description,
		ComplainType: complainType,
		Hostel:       hostel,
	}

	// Handle the file upload
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the file to the server
	filePath := "./uploads/" + header.Filename
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Could not save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Use io.Copy to save the file
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Insert complaint details into the database
	insertSQL := `INSERT INTO complaints (email, name, room, mobile, description, complain_type, hostel, file_path) 
				  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(insertSQL, complaint.Email, complaint.Name, complaint.Room, complaint.Mobile, complaint.Description, complaint.ComplainType, complaint.Hostel, filePath)
	if err != nil {
		http.Error(w, "Error saving complaint to the database", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{
		"valid":        true,
		"isRegistered": false,
	})
}

func verifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s request for %s", r.Method, r.URL.Path)

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
