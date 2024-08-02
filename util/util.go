package util

import "github.com/ethereum/go-ethereum/common"

func HashToAddress(hash common.Hash) common.Address {
	return common.BytesToAddress(hash.Bytes())
}
