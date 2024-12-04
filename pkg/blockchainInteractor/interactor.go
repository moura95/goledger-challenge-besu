package blockchainInteractor

import (
	"context"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Interactor struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
	privateKey      string
}

func NewBlockchainInteractor(clientURL, contractAddr, privateKey string, abiJSON string) (*Interactor, error) {
	client, err := ethclient.Dial(clientURL)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao cliente Ethereum: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		return nil, fmt.Errorf("falha ao carregar ABI: %w", err)
	}

	return &Interactor{
		client:          client,
		contractAddress: common.HexToAddress(contractAddr),
		contractABI:     contractABI,
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
	tx, err := contract.Transact(auth, "set", value)
	if err != nil {
		return "", fmt.Errorf("falha ao enviar transação: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (i *Interactor) GetValue() (uint64, error) {
	var result []interface{}
	var value uint64
	contract := bind.NewBoundContract(i.contractAddress, i.contractABI, i.client, i.client, i.client)

	err := contract.Call(nil, &result, "get")
	if err != nil {
		return 0, fmt.Errorf("falha ao ler valor do contrato: %w", err)
	}
	if len(result) == 0 {
		return 0, fmt.Errorf("falha ao ler valor do contrato: %w", err)
	}
	value = result[0].(uint64)

	return value, nil
}

func (i *Interactor) Close() {
	i.client.Close()
}
