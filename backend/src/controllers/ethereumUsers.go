// etherereum user controller
package controllers

import (
	"Delingo/src/models"
	"context"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Infura URL and Token Address
const infuraURL = "https://mainnet.infura.io/v3/YOUR_INFURA_KEY"
const tokenAddress = "0xYourERC20TokenAddress"

// ERC-20 ABI for balanceOf function
var erc20ABI = `[{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

// DepositERC20 checks the balance of an ERC20 token for a given wallet address
func DepositERC20(walletAddress string) (*models.Token, error) {
	// Create a new Ethereum client with Infura
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return nil, err
	}

	tokenContractAddress := common.HexToAddress(tokenAddress)
	tokenABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		log.Fatalf("Failed to parse ERC20 ABI: %v", err)
		return nil, err
	}

	balance, err := getERC20Balance(client, tokenContractAddress, tokenABI, walletAddress)
	if err != nil {
		log.Fatalf("Failed to get ERC-20 balance: %v", err)
		return nil, err
	}

	return &models.Token{
		Address: walletAddress,
		Balance: balance,
	}, nil
}

// getERC20Balance interacts with the Ethereum network to get the ERC20 token balance for a given address
func getERC20Balance(client *ethclient.Client, contractAddress common.Address, tokenABI abi.ABI, address string) (*big.Int, error) {
	addressHex := common.HexToAddress(address)
	data, err := tokenABI.Pack("balanceOf", addressHex)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var balance *big.Int
	err = tokenABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
