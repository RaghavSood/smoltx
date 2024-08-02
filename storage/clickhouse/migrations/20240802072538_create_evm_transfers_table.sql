-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS evm_transfers (
  chain_id INTEGER NOT NULL,
  tx_id TEXT NOT NULL,
  token_address TEXT NOT NULL,
  from_address TEXT NOT NULL,
  to_address TEXT NOT NULL,
  value TEXT NOT NULL,
  block_height INTEGER NOT NULL,
  block_hash TEXT NOT NULL,
  tx_log_index INTEGER NOT NULL,
  log_index INTEGER NOT NULL,
  block_timestamp DateTime NOT NULL,
  PRIMARY KEY (chain_id, tx_id, token_address, tx_log_index)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS evm_transfers;
-- +goose StatementEnd
