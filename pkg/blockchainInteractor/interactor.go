package blockchainInteractor

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const contractABI = `[{"inputs":[],"name":"get","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"x","type":"uint256"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

type Interactor struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
	privateKey      string
}

func NewBlockchainInteractor(clientURL, contractAddr, privateKey string) (*Interactor, error) {
	client, err := ethclient.Dial(clientURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao cliente Ethereum: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatal(err)
	}

	return &Interactor{
		client:          client,
		contractAddress: common.HexToAddress(contractAddr),
		contractABI:     parsedABI,
		privateKey:      privateKey,
	}, nil
}

func (i *Interactor) SetValue(value uint64) (string, error) {
	privateKey, err := crypto.HexToECDSA(i.privateKey)
	if err != nil {
		return "", fmt.Errorf("falha ao carregar chave privada: %w", err)
	}

	chainID, err := i.client.ChainID(context.Background())
	if err != nil {
		return "", fmt.Errorf("falha ao obter Chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("falha ao criar transação assinada: %w", err)
	}

	contract := bind.NewBoundContract(i.contractAddress, i.contractABI, i.client, i.client, i.client)
	log.Printf("Contrato criado: %v", contract)

	bigValue := big.NewInt(int64(value))

	tx, err := contract.Transact(auth, "set", bigValue)
	if err != nil {
		return "", fmt.Errorf("falha ao enviar transação: %w", err)
	}

	fmt.Println("waiting until transaction is mined",
		"tx", tx.Hash().Hex(),
	)

	receipt, err := bind.WaitMined(
		context.Background(),
		i.client,
		tx,
	)
	if err != nil {
		log.Fatalf("error waiting for transaction to be mined: %v", err)
	}

	return receipt.TxHash.Hex(), nil
}

func (i *Interactor) GetValue() (interface{}, error) {
	var result interface{}

	opts := bind.CallOpts{
		Pending: false,
		Context: context.Background(),
	}

	contract := bind.NewBoundContract(i.contractAddress, i.contractABI, i.client, i.client, i.client)
	var output []interface{}

	err := contract.Call(&opts, &output, "get")
	if err != nil {
		log.Fatalf("error calling contract: %v", err)
	}

	result = output

	fmt.Println("Successfully called contract!", result)

	return result, err
}

func (i *Interactor) Close() {
	i.client.Close()
}
