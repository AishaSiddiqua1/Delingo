package controllers

import (
	"Delingo/src/models"
	"Delingo/src/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterUser handles email-based user registration
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Ensure the user provides an email and password
	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Clear wallet-related fields for email registration
	user.EthereumWalletAddr = ""
	user.SolanaWalletAddr = ""
	user.RegistrationMethod = "email"

	// Save user to the database
	query := `INSERT INTO users (username, email, password, registration_method, created_at) 
              VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	err = utils.SQLDB.QueryRow(query, user.Username, user.Email, user.Password, user.RegistrationMethod).Scan(&user.ID)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	// Respond with user details
	json.NewEncoder(w).Encode(user)
}

// GetUser retrieves a user profile by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	query := `SELECT id, username, email, wallet_address, created_at FROM users WHERE id=$1`
	row := utils.SQLDB.QueryRow(query, userID)
	// query := `SELECT id, username, email, ethereum_wallet_addr, solana_wallet_addr, created_at FROM users WHERE id=$1`
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.EthereumWalletAddr, &user.SolanaWalletAddr, &user.CreatedAt)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates a user's profile information
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE users SET username=$1, email=$2 WHERE id=$3`
	_, err = utils.SQLDB.Exec(query, user.Username, user.Email, userID)
	if err != nil {
		http.Error(w, "Could not update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser deletes a user profile
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	query := `DELETE FROM users WHERE id=$1`
	_, err := utils.SQLDB.Exec(query, userID)
	if err != nil {
		http.Error(w, "Could not delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted"))
}
