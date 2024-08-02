package main

import (
	"flag"
	"os"

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
}
