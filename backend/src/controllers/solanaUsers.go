package controllers

import (
	"Delingo/src/models"
	"Delingo/src/utils"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterSolanaUser registers a user using a Solana wallet
func RegisterSolanaUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the user has a valid Solana wallet address
	if user.SolanaWalletAddr == "" {
		http.Error(w, "Solana wallet address required", http.StatusBadRequest)
		return
	}

	// Clear other fields not needed for wallet-based sign-up
	user.Email = ""
	user.Password = ""

	// Save user to the database
	query := `INSERT INTO users (username, email, wallet_address, password, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	err = utils.SQLDB.QueryRow(query, user.Username, user.Email, user.SolanaWalletAddr, user.Password).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	// Respond with user details
	json.NewEncoder(w).Encode(user)
}

// GetSolanaUser retrieves a user by their Solana wallet address
func GetSolanaUser(w http.ResponseWriter, r *http.Request) {
	// Extract wallet address from URL params
	walletAddr := mux.Vars(r)["walletAddr"]

	// Ensure the wallet address is provided
	if walletAddr == "" {
		http.Error(w, "Solana wallet address required", http.StatusBadRequest)
		return
	}

	// Query the database for a user with the provided wallet address
	var user models.User
	query := `SELECT id, username, email, wallet_address, created_at FROM users WHERE wallet_address = $1`
	err := utils.SQLDB.QueryRow(query, walletAddr).Scan(&user.ID, &user.Username, &user.Email, &user.SolanaWalletAddr, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		}
		return
	}

	// Respond with user details
	json.NewEncoder(w).Encode(user)
}
