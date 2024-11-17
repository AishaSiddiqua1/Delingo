package models

type Badge struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"not null"`
	Description string
	Criteria    string `gorm:"type:jsonb"` // Using jsonb for storing criteria, if needed
	TokenID     uint64 `gorm:"not null"`   // Storing the NFT's tokenID
	TokenURI    string // URI for badge metadata
}
