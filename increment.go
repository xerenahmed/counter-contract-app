package main

import (
	"context"
	"counter_app/contract"
	"flag"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/joho/godotenv"
)

var times int64
var silent bool

func init() {
	flag.Int64Var(&times, "times", 1, "times to increment")
	flag.BoolVar(&silent, "q", false, "silent mode")
	flag.Parse()

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

	privateKey, err := crypto.HexToECDSA(os.Getenv("WALLET_PRIVATE_KEY"))
	if err != nil {
		log.Fatalf("failed to load private key: %v", err)
	}

	gas, err := conn.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("failed to suggest gas price: %v", err)
	}

	// convert gas to as ethereum
	gas = gas.Mul(gas, big.NewInt(100))
	gas = gas.Div(gas, big.NewInt(params.GWei))
	gasFloat := float64(gas.Int64()) / 100000000000 // TODO: test for Ethereum. This just tested with AVAX

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(43113))
	if err != nil {
		log.Fatalf("failed to create authorized transactor: %v", err)
	}

	var ethUsed float64
	for i := int64(0); i < times; i++ {
		tx, err := counter.Increment(auth)
		if err != nil {
			log.Printf("failed to create tx: %v", err)
			continue
		}

		if !silent {
			log.Printf("Increment tx sent\nTx pending: %s/%s\n", os.Getenv("EXPLORER_URL"), tx.Hash().Hex())
		}

		// wait for the transaction to be mined
		receipt, err := bind.WaitMined(context.Background(), conn, tx)
		if err != nil {
			log.Printf("error waiting for transaction to be mined: %v", err)
			continue
		}

		ethUsed += float64(receipt.GasUsed) * gasFloat

		if receipt.Status != 1 {
			log.Printf("transaction failed with status %d", receipt.Status)
			continue
		}
	}

	log.Println("All transactions mined")
	log.Printf("Total eth used: %f\n", ethUsed)
}
