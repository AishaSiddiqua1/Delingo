package controllers

import (
	"Delingo/src/models"
	"Delingo/src/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterWalletUser handles wallet-based user registration
func RegisterWalletUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the user provides a wallet address (Ethereum or Solana)
	if user.EthereumWalletAddr == "" && user.SolanaWalletAddr == "" {
		http.Error(w, "Ethereum or Solana wallet address required", http.StatusBadRequest)
		return
	}

	// Clear email and password for wallet-based registration
	user.Email = ""
	user.Password = ""

	// Set registration method to wallet
	user.RegistrationMethod = "wallet"

	// Save user to the database, choosing the correct wallet address
	query := `INSERT INTO users (username, email, ethereum_wallet_address, solana_wallet_address, registration_method, created_at) 
              VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`
	err = utils.SQLDB.QueryRow(query, user.Username, user.Email, user.EthereumWalletAddr, user.SolanaWalletAddr, user.RegistrationMethod).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	// Respond with user details
	json.NewEncoder(w).Encode(user)
}

// UpdateWalletAddress updates a user's wallet address
func UpdateWalletAddress(w http.ResponseWriter, r *http.Request) {
	// Get user ID from the URL
	vars := mux.Vars(r)
	userID := vars["id"]

	// Parse the request body into a user object
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Determine which wallet address to update
	var query string
	var walletAddress string

	// If Ethereum wallet address is provided, update it
	if user.EthereumWalletAddr != "" {
		query = `UPDATE users SET ethereum_wallet_address=$1 WHERE id=$2`
		walletAddress = user.EthereumWalletAddr
	} else if user.SolanaWalletAddr != "" {
		// If Solana wallet address is provided, update it
		query = `UPDATE users SET solana_wallet_address=$1 WHERE id=$2`
		walletAddress = user.SolanaWalletAddr
	} else {
		// If neither address is provided, return an error
		http.Error(w, "No wallet address provided", http.StatusBadRequest)
		return
	}

	// Execute the update query
	_, err = utils.SQLDB.Exec(query, walletAddress, userID)
	if err != nil {
		http.Error(w, "Could not update wallet address", http.StatusInternalServerError)
		return
	}

	// Send the response back
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
