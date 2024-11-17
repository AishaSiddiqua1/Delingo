package routes

import (
	"Delingo/src/controllers"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	// Email-based User Routes
	router.HandleFunc("/api/users", controllers.RegisterUser).Methods("POST")      // Email user registration
	router.HandleFunc("/api/users/{id}", controllers.GetUser).Methods("GET")       // Get email user profile
	router.HandleFunc("/api/users/{id}", controllers.UpdateUser).Methods("PUT")    // Update email user profile
	router.HandleFunc("/api/users/{id}", controllers.DeleteUser).Methods("DELETE") // Delete email user

	// Solana-specific routes
	router.HandleFunc("/api/solana/register", controllers.RegisterSolanaUser).Methods("POST")
	router.HandleFunc("/api/solana/{id}", controllers.GetSolanaUser).Methods("GET")

	// Ethereum Wallet User Routes
	router.HandleFunc("/api/wallet/register", controllers.RegisterWalletUser).Methods("POST") // Register wallet user
	router.HandleFunc("/api/wallet/{id}", controllers.UpdateWalletAddress).Methods("PUT")     // Update wallet address for user
	router.HandleFunc("/api/wallet/{id}", controllers.GetUser).Methods("GET")                 // Get user profile by wallet address

}
func SetupRoutes(router *gin.Engine) {
	// Profile Routes
	router.GET("/api/profile", controllers.GetUserProfile)    // Get profile
	router.PUT("/api/profile", controllers.UpdateUserProfile) // Update profile
}
