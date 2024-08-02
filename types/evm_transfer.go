package types

import (
	"math/big"
	"time"
)

type EvmTransfer struct {
	ChainID      *big.Int
	TxID         string
	TokenAddress string
	FromAddress  string
	ToAddress    string
	Value        *big.Int
	BlockHeight  *big.Int
	BlockHash    string
	TxLogIndex   int
	LogIndex     int
	BlockTime    time.Time
}
