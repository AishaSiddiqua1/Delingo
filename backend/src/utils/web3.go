// utils/web3.go
package utils

import (
	"Delingo/src/config"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client

// InitWeb3 connects to the Ethereum network via Infura
func InitWeb3() {
	var err error
	client, err = ethclient.Dial(config.HeklaRPCURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
}

// GetClient returns the Ethereum client
func GetClient() *ethclient.Client {
	return client
}

// LoadContractABI loads the ABI for the contract
func LoadContractABI(abiFile string) abi.ABI {
	// Load ABI from a file or a string
	contractABI, err := abi.JSON(strings.NewReader(abiFile))
	if err != nil {
		log.Fatalf("Failed to load contract ABI: %v", err)
	}
	return contractABI
}

// DeployContract interacts with the deployed contract
func DeployContract(contractABI abi.ABI, contractAddress string) *common.Address {
	address := common.HexToAddress(contractAddress)
	return &address
}
