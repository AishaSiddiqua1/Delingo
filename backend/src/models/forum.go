// models/post.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	ThreadID uint
	UserID   uint
	Content  string
	Votes    []Vote
}

type Vote struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`
	PostID    *int           `json:"post_id,omitempty"`   // Optional for thread votes
	ThreadID  *int           `json:"thread_id,omitempty"` // Optional for post votes
	VoteValue int            `json:"vote_value"`          // 1 for upvote, -1 for downvote
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

type Thread struct {
	gorm.Model
	UserID uint
	Title  string
	Posts  []Post
}

type Comment struct {
	ID        int       `json:"id" gorm:"primary_key"`
	PostID    int       `json:"post_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
