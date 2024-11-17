package routes

import (
	"Delingo/src/controllers"

	"github.com/gin-gonic/gin"
)

func ForumRoutes(r *gin.Engine) {
	// Grouping the forum-related routes
	forumGroup := r.Group("/forum")
	{
		// Thread Routes
		forumGroup.POST("/thread", controllers.CreateThread)       // Create a new thread
		forumGroup.GET("/threads", controllers.GetAllThreads)      // Get all threads
		forumGroup.GET("/thread/:id", controllers.GetThread)       // Get a specific thread by ID
		forumGroup.PUT("/thread/:id", controllers.UpdateThread)    // Update an existing thread
		forumGroup.DELETE("/thread/:id", controllers.DeleteThread) // Delete a specific thread

		// Post Routes
		forumGroup.POST("/post", controllers.CreatePost)             // Create a new post
		forumGroup.GET("/post/:id", controllers.GetPost)             // Get a post by ID
		forumGroup.PUT("/post/:id", controllers.UpdatePost)          // Update a specific post
		forumGroup.DELETE("/post/:id", controllers.DeletePost)       // Delete a specific post
		forumGroup.GET("/posts/:thread_id", controllers.GetAllPosts) // Get all posts in a thread

		// Comment Routes
		forumGroup.POST("/comment/:post_id", controllers.CreateComment)     // Create a new comment on a post
		forumGroup.PUT("/comment/:id", controllers.UpdateComment)           // Update a comment
		forumGroup.GET("/comments/:post_id", controllers.GetCommentsByPost) // Get all comments for a post
		forumGroup.GET("/comment/:id", controllers.GetComment)              // Get a comment by ID
		forumGroup.DELETE("/comment/:id", controllers.DeleteComment)        // Delete a comment

		// Vote Routes
		forumGroup.POST("/vote/thread/:thread_id", controllers.VoteOnThread) // Vote on a thread
		forumGroup.POST("/vote/post/:post_id", controllers.VoteOnPost)       // Vote on a post

		// User Votes
		forumGroup.GET("/votes/user/:user_id", controllers.GetUserVotes) // Get all votes by a user

		// Search Forum
		forumGroup.GET("/search", controllers.SearchForum) // Search threads and posts in the forum
	}
}
