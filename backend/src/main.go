package main

import (
	"Delingo/src/routes"
	"Delingo/src/utils"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	err := utils.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Initialize logger
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// Initialize the router
	router := mux.NewRouter()

	// Create a new Gin engine
	ginEngine := gin.Default()

	// Register Gin routes
	routes.ForumRoutes(ginEngine)

	// Serve Gin on a specific path (e.g., "/api")
	router.Handle("/api/", http.StripPrefix("/api", ginEngine))

	// Register routes
	routes.RegisterRoutes(router)

	// Add basic routes
	router.HandleFunc("/signup", signupHandler).Methods(http.MethodPost)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodPost)

	// Add a health check route
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check triggered")
		w.Write([]byte("Hello, World!"))
	})

	// Start the server
	port := "8080"
	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// signupHandler handles user signup requests
func signupHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Signup handler triggered")
	// Process signup request (dummy response for now)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
	log.Println("Signup completed")
}

// loginHandler handles user login requests
func loginHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Login handler triggered")
	// Process login request (dummy response for now)
	w.Write([]byte("User logged in successfully"))
	log.Println("Login completed")
}
