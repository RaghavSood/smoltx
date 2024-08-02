package evm

import (
	"context"
	"math/big"
	"time"

	"github.com/RaghavSood/smoltx/smollogger"
	"github.com/RaghavSood/smoltx/storage"
	"github.com/RaghavSood/smoltx/types"
	"github.com/RaghavSood/smoltx/util"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	gtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var log = smollogger.NewLogger("evm")

type EVM struct {
	db      storage.Storage
	rpc     *ethclient.Client
	chainID *big.Int
}

func NewEVM(rpcURL string, db storage.Storage) (*EVM, error) {
	rpc, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	chainID, err := rpc.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	return &EVM{
		db:      db,
		rpc:     rpc,
		chainID: chainID,
	}, nil
}

func (e *EVM) Close() {
	e.rpc.Close()
}

func (e *EVM) IndexHeaders() error {
	headers := make(chan *gtypes.Header)
	sub, err := e.rpc.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to subscribe to new headers")
	}

	for {
		select {
		case err := <-sub.Err():
			log.Error().
				Err(err).
				Msg("Subscription error")
		case header := <-headers:
			e.processEVMBlock(header.Number, header.Time)
		}
	}
}

func (e *EVM) processEVMBlock(height *big.Int, timestamp uint64) {
	blockTime := time.Unix(int64(timestamp), 0)
	topic := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	query := ethereum.FilterQuery{
		Topics:    [][]common.Hash{{topic}},
		FromBlock: height,
		ToBlock:   height,
	}

	log.Debug().
		Int64("height", height.Int64()).
		Msg("Processing block")

	logs, err := e.rpc.FilterLogs(context.Background(), query)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to filter logs")
	}

	var transfers []types.EvmTransfer
	for _, entry := range logs {
		if len(entry.Topics) != 3 {
			continue
		}

		log.Info().
			Str("tx", entry.TxHash.Hex()).
			Str("contract", entry.Address.Hex()).
			Str("from", util.HashToAddress(entry.Topics[1]).Hex()).
			Str("to", util.HashToAddress(entry.Topics[2]).Hex()).
			Msg("Transfer")

		transfer := types.EvmTransfer{
			ChainID:      e.chainID,
			TxID:         entry.TxHash.Hex(),
			TokenAddress: entry.Address.Hex(),
			FromAddress:  util.HashToAddress(entry.Topics[1]).Hex(),
			ToAddress:    util.HashToAddress(entry.Topics[2]).Hex(),
			Value:        new(big.Int).SetBytes(entry.Data),
			BlockHeight:  height,
			BlockHash:    entry.BlockHash.Hex(),
			BlockTime:    blockTime,
			TxLogIndex:   int(entry.TxIndex),
			LogIndex:     int(entry.Index),
		}

		transfers = append(transfers, transfer)
	}

	// Store the transfers
	err = e.db.RecordEvmTransfers(transfers)
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to record transfers")
	}
}
