// config/config.go
package config

import (
	"log"
	"os"
)

var HeklaRPCURL = os.Getenv("HEKLA_RPC_URL")

func LoadConfig() {
	HeklaRPCURL = os.Getenv("HEKLA_RPC_URL")
	if HeklaRPCURL == "" {
		log.Fatal("HEKLA_RPC_URL is required")
	}
}
