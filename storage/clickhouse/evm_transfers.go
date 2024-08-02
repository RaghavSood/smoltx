package clickhouse

import (
	"context"
	"fmt"

	"github.com/RaghavSood/smoltx/types"
)

func (d *ClickhouseBackend) RecordEvmTransfers(transfers []types.EvmTransfer) error {
	batch, err := d.db.PrepareBatch(context.Background(), "INSERT INTO evm_transfers")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, transfer := range transfers {
		err := batch.Append(
			transfer.ChainID.Int64(),
			transfer.TxID,
			transfer.TokenAddress,
			transfer.FromAddress,
			transfer.ToAddress,
			transfer.Value,
			transfer.BlockHeight.Int64(),
			transfer.BlockHash,
			transfer.TxLogIndex,
			transfer.LogIndex,
			transfer.BlockTime,
		)
		if err != nil {
			return fmt.Errorf("failed to append to batch: %w", err)
		}
	}

	return batch.Send()

}
