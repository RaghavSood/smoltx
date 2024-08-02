package main

import (
	"flag"
	"os"

	"github.com/RaghavSood/smoltx/evm"
	"github.com/RaghavSood/smoltx/storage/clickhouse"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var noindex bool

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: log.Output(zerolog.ConsoleWriter{Out: os.Stderr})})

	flag.BoolVar(&noindex, "noindex", false, "Don't index the blockchain, run in read-only mode")
	flag.Parse()
}

func main() {
	db, err := clickhouse.NewClickhouseBackend(noindex)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize storage backend")
	}

	evmIndexer, err := evm.NewEVM(os.Getenv("ETH_RPC"), db)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize EVM backend")
	}

	evmIndexer.IndexHeaders()
}
