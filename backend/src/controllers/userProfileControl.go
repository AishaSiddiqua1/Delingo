package controllers

import (
	"Delingo/src/models" // Import models for Profile
	"Delingo/src/utils"  // Database utils
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin" // Add the gin import
)

// GET /profile - Fetch the user's profile
func GetUserProfile(c *gin.Context) {
	// Get the user from the token
	userID, err := getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists before fetching profile
	if !userExists(c, userID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Query the DB for the user's profile
	profile, err := GetUserProfileFromDB(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user profile"})
		return
	}

	// Return the profile as a JSON response
	c.JSON(http.StatusOK, profile)
}

// PUT /profile - Update the user's profile
func UpdateUserProfile(c *gin.Context) {
	// Get the user from the token
	userID, err := getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists before updating
	if !userExists(c, userID) {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var newProfile models.Profile // Declare a variable to hold the updated profile

	// Decode the request body into the Profile model
	if err := json.NewDecoder(c.Request.Body).Decode(&newProfile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the function to update the user's profile
	err = UpdateUserProfileInDB(userID, newProfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// Fetch user profile from the DB
func GetUserProfileFromDB(userID uint) (models.Profile, error) {
	var profile models.Profile
	query := `SELECT level, avatar_url, bio, preferences, location, website, social_links 
			  FROM profiles WHERE user_id = $1`

	row := utils.SQLDB.QueryRow(query, userID)
	err := row.Scan(&profile.Level, &profile.AvatarURL, &profile.Bio, &profile.Preferences,
		&profile.Location, &profile.Website, &profile.SocialLinks)

	if err != nil {
		log.Println("Error fetching user profile:", err)
		return profile, err
	}

	return profile, nil
}

// Update user profile in the DB
func UpdateUserProfileInDB(userID uint, newProfile models.Profile) error {
	// Ensure the user exists before attempting to update the profile
	if !userExists(nil, userID) {
		return fmt.Errorf("user does not exist")
	}

	query := `UPDATE profiles 
			  SET level = $1, avatar_url = $2, bio = $3, preferences = $4, location = $5, website = $6, social_links = $7, updated_at = NOW() 
			  WHERE user_id = $8`

	_, err := utils.SQLDB.Exec(query,
		newProfile.Level,
		newProfile.AvatarURL,
		newProfile.Bio,
		newProfile.Preferences,
		newProfile.Location,
		newProfile.Website,
		newProfile.SocialLinks,
		userID,
	)

	if err != nil {
		log.Println("Error updating profile:", err)
		return err
	}

	return nil
}

// Check if the user exists in the users table
func userExists(c *gin.Context, userID uint) bool {
	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)`
	err := utils.SQLDB.QueryRow(query, userID).Scan(&exists)
	if err != nil || !exists {
		if c != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		}
		return false
	}
	return true
}

// getUserFromToken extracts the user ID from the JWT token in the context
func getUserFromToken(c *gin.Context) (uint, error) {
	// Retrieve the userID from the context
	userID, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("User not authenticated")
	}

	// Type assert userID to uint (since we stored it as uint in the context)
	return userID.(uint), nil
}
