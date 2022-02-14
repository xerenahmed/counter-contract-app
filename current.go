package main

import (
	"counter_app/contract"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to load .env file: %v", err)
	}
}

func main() {
	conn, err := ethclient.Dial(os.Getenv("RPC_URL"))
	if err != nil {
		log.Fatalf("failed to connect to Ethereum client: %v", err)
	}

	counter, err := contract.NewCounter(common.HexToAddress(os.Getenv("CONTRACT_ADDRESS")), conn)
	if err != nil {
		log.Fatalf("failed to instantiate contract: %v", err)
	}

	val, err := counter.Current(nil)
	if err != nil {
		log.Fatalf("failed to retrieve current value: %v", err)
	}

	fmt.Printf("Current value is %d\n", val.Int64())
}
