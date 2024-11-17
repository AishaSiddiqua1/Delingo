// models/models.go
package models

import (
	"math/big"
)

// LingTokenInfo represents basic information about the Ling ERC-20 token
type LingTokenInfo struct {
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Decimals    uint8    `json:"decimals"`
	TotalSupply *big.Int `json:"total_supply"`
}

// NFTInfo represents information about an NFT token, like LingNFT
type NFTInfo struct {
	TokenID  *big.Int `json:"token_id"`
	Owner    string   `json:"owner"`
	Metadata string   `json:"metadata"`
}

// WalletBalance represents a user's balance for both ERC20 and ERC721 tokens
type WalletBalance struct {
	Address      string   `json:"address"`
	TokenBalance *big.Int `json:"token_balance"` // ERC20 balance
	NFTCount     int      `json:"nft_count"`     // Number of NFTs owned
}

// Transaction represents a basic structure to store transaction data
type Transaction struct {
	From        string   `json:"from"`
	To          string   `json:"to"`
	Value       *big.Int `json:"value"`
	TxHash      string   `json:"tx_hash"`
	BlockNumber uint64   `json:"block_number"`
}
