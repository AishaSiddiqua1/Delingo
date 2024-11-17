package controllers

import (
	"Delingo/src/models"
	"Delingo/src/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateThread(c *gin.Context) {
	var thread models.Thread
	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Assuming userID is extracted from the JWT or session
	thread.UserID = 1 // Replace with actual userID extraction

	// Insert the thread into the database
	if err := utils.GormDB.Create(&thread).Error; err != nil {
		log.Println("Error creating thread:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create thread"})
		return
	}

	c.JSON(http.StatusOK, thread)
}

func GetThread(c *gin.Context) {
	threadID := c.Param("id")
	var thread models.Thread

	if err := utils.GormDB.Preload("Posts").First(&thread, threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
		return
	}

	c.JSON(http.StatusOK, thread)
}

// VoteOnPost allows a user to vote on a post. It ensures a user can only vote once per post.
// VoteOnPost allows users to vote on a post
func VoteOnPost(c *gin.Context) {
	// Retrieve user ID from JWT
	userID, _ := c.Get("userID")

	// Parse vote value (1 for upvote, -1 for downvote)
	var input struct {
		PostID    int `json:"post_id"`
		VoteValue int `json:"vote_value"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate vote value
	if input.VoteValue != 1 && input.VoteValue != -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote value"})
		return
	}

	// Check if the user has already voted on this post
	var existingVote models.Vote
	err := utils.GormDB.Where("user_id = ? AND post_id = ?", userID, input.PostID).First(&existingVote).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// If the user has already voted, update the vote
	if err == nil {
		// If the vote value is the same as the existing vote, cancel the vote (delete it)
		if existingVote.VoteValue == input.VoteValue {
			if err := utils.GormDB.Delete(&existingVote).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vote"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Vote removed"})
			return
		} else {
			// Otherwise, update the existing vote
			existingVote.VoteValue = input.VoteValue
			existingVote.UpdatedAt = time.Now()
			if err := utils.GormDB.Save(&existingVote).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Vote updated"})
			return
		}
	}

	// If no existing vote, create a new one
	newVote := models.Vote{
		UserID:    userID.(uint),
		PostID:    &input.PostID,
		VoteValue: input.VoteValue,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := utils.GormDB.Create(&newVote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote created"})
}

// Create a new post
func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	result := utils.GormDB.Create(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post created successfully", "post": post})
}

func GetPost(c *gin.Context) {
	id := c.Param("id")
	postID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	result := utils.GormDB.First(&post, postID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetAllThreads retrieves all threads from the database
func GetAllThreads(c *gin.Context) {
	var threads []models.Thread

	// Query all threads from the database
	if err := utils.GormDB.Find(&threads).Error; err != nil {
		log.Println("Error retrieving threads:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve threads"})
		return
	}

	c.JSON(http.StatusOK, threads)
}

// UpdateThread updates an existing thread
func UpdateThread(c *gin.Context) {
	threadID := c.Param("id")
	var thread models.Thread

	// Find the thread by ID
	if err := utils.GormDB.First(&thread, threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
		return
	}

	// Bind the updated values from the request body
	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the thread in the database
	if err := utils.GormDB.Save(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update thread"})
		return
	}

	c.JSON(http.StatusOK, thread)
}

// DeleteThread deletes a thread from the database
func DeleteThread(c *gin.Context) {
	threadID := c.Param("id")
	var thread models.Thread

	// Find the thread by ID
	if err := utils.GormDB.First(&thread, threadID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Thread not found"})
		return
	}

	// Delete the thread from the database
	if err := utils.GormDB.Delete(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete thread"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thread deleted successfully"})
}

// UpdatePost updates an existing post
func UpdatePost(c *gin.Context) {
	postID := c.Param("id")
	var post models.Post

	// Find the post by ID
	if err := utils.GormDB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Bind the updated values from the request body
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the post in the database
	if err := utils.GormDB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

// GetAllPosts retrieves all posts for a specific thread
func GetAllPosts(c *gin.Context) {
	threadID := c.Param("thread_id")
	var posts []models.Post

	// Query all posts related to the thread
	if err := utils.GormDB.Where("thread_id = ?", threadID).Find(&posts).Error; err != nil {
		log.Println("Error retrieving posts:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// DeletePost deletes a post from the database
func DeletePost(c *gin.Context) {
	postID := c.Param("id")
	var post models.Post

	// Find the post by ID
	if err := utils.GormDB.First(&post, postID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Delete the post from the database
	if err := utils.GormDB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// GetUserVotes retrieves all votes cast by a user
func GetUserVotes(c *gin.Context) {
	userID := c.Param("user_id")
	var votes []models.Vote

	// Query all votes by the user
	if err := utils.GormDB.Where("user_id = ?", userID).Find(&votes).Error; err != nil {
		log.Println("Error retrieving votes:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve votes"})
		return
	}

	c.JSON(http.StatusOK, votes)
}

// CreateComment adds a comment to a specific post
func CreateComment(c *gin.Context) {
	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get postID from the URL parameter
	postIDStr := c.Param("post_id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	comment.PostID = int(postID)

	// Save the comment to the database
	if err := utils.GormDB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// In controllers/commentController.go
func UpdateComment(c *gin.Context) {
	commentID := c.Param("id")
	var comment models.Comment

	// Find the comment by ID
	if err := utils.GormDB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Bind the updated values
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save updated comment
	if err := utils.GormDB.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetCommentsByPost retrieves all comments for a specific post
func GetCommentsByPost(c *gin.Context) {
	postID := c.Param("post_id")
	var comments []models.Comment

	if err := utils.GormDB.Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GetComment retrieves a comment by its ID
func GetComment(c *gin.Context) {
	commentID := c.Param("id")
	var comment models.Comment

	if err := utils.GormDB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment deletes a comment from the database
func DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	var comment models.Comment

	// Find the comment by ID
	if err := utils.GormDB.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Delete the comment from the database
	if err := utils.GormDB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// VoteOnThread allows users to vote on a thread
func VoteOnThread(c *gin.Context) {
	// Retrieve user ID from JWT
	userID, _ := c.Get("userID")

	// Parse vote value (1 for upvote, -1 for downvote)
	var input struct {
		ThreadID  int `json:"thread_id"`
		VoteValue int `json:"vote_value"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate vote value
	if input.VoteValue != 1 && input.VoteValue != -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote value"})
		return
	}

	// Check if the user has already voted on this thread
	var existingVote models.Vote
	err := utils.GormDB.Where("user_id = ? AND thread_id = ?", userID, input.ThreadID).First(&existingVote).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// If the user has already voted, update the vote
	if err == nil {
		// If the vote value is the same as the existing vote, cancel the vote (delete it)
		if existingVote.VoteValue == input.VoteValue {
			if err := utils.GormDB.Delete(&existingVote).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vote"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Vote removed"})
			return
		} else {
			// Otherwise, update the existing vote
			existingVote.VoteValue = input.VoteValue
			existingVote.UpdatedAt = time.Now()
			if err := utils.GormDB.Save(&existingVote).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Vote updated"})
			return
		}
	}

	// If no existing vote, create a new one
	newVote := models.Vote{
		UserID:    userID.(uint),
		ThreadID:  &input.ThreadID,
		VoteValue: input.VoteValue,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := utils.GormDB.Create(&newVote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote created"})
}

// SearchForum searches for threads or posts in the forum
func SearchForum(c *gin.Context) {
	query := c.DefaultQuery("query", "")
	var threads []models.Thread
	var posts []models.Post

	// Search threads
	if err := utils.GormDB.Where("title LIKE ?", "%"+query+"%").Find(&threads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search threads"})
		return
	}

	// Search posts
	if err := utils.GormDB.Where("content LIKE ?", "%"+query+"%").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"threads": threads, "posts": posts})
}
