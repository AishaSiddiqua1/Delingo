package models

import (
	"math/big"
	"time"
)

type User struct {
	ID                 int       `json:"id"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	Password           string    `json:"password"`
	EthereumWalletAddr string    `json:"ethereum_wallet_address"`
	SolanaWalletAddr   string    `json:"solana_wallet_address"`
	RegistrationMethod string    `json:"registration_method"`
	CreatedAt          time.Time `json:"created_at"`
}

type Progress struct {
	ID       int     `json:"id"`
	UserID   int     `json:"user_id"`
	QuizID   int     `json:"quiz_id"`
	Status   string  `json:"status"`   // e.g., 'completed', 'in-progress'
	Progress float32 `json:"progress"` // percentage
}

// Token represents a token deposit object.
type Token struct {
	Address string   `json:"address"`
	Balance *big.Int `json:"balance"`
}
type Profile struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Bio         string    `json:"bio"`
	Level       int       `json:"level"`
	AvatarURL   string    `json:"avatar_url"`
	Location    string    `json:"location"`
	Website     string    `json:"website"`
	Preferences string    `json:"preferences"`
	SocialLinks string    `json:"social_links"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
