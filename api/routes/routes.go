package routes

import (
	"net/http"

	"github.com/Yagna-rao15/hostel-portal/api/handlers" // Use full module path
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Authentication and User Routes
	router.HandleFunc("/login", handlers.CheckEmailHandler).Methods("POST")
	router.HandleFunc("/password", handlers.CheckPasswordHandler).Methods("POST")
	router.HandleFunc("/verify", handlers.VerifyEmailHandler).Methods("POST")
	router.HandleFunc("/update-password", handlers.UpdatePasswordHandler).Methods("POST")

	// Complaint Routes
	router.HandleFunc("/form", handlers.SubmitComplaintHandler).Methods("POST")

	// Health check route for the database
	router.HandleFunc("/test-db", handlers.TestDBHandler).Methods("GET")

	return router
}
