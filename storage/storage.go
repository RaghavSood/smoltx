package storage

import (
	"github.com/RaghavSood/smoltx/types"
)

type Storage interface {
	RecordEvmTransfers(transfers []types.EvmTransfer) error
}
