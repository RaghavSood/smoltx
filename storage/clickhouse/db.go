package sqlite

import (
	"context"
	"crypto/tls"
	"database/sql"
	"embed"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*
var embeddedMigrations embed.FS

type ClickhouseBackend struct {
	db    *clickhouse.Conn
	sqlDb *sql.DB
}

func NewClickhouseBackend(readonly bool) (*ClickhouseBackend, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: fmt.Sprintf("%s:%d", os.Getenv("CLICKHOUSE_HOST"), os.Getenv("CLICKHOUSE_PORT")),
		Auth: clickhouse.Auth{
			Database: os.Getenv("CLICKHOUSE_DATABASE"),
			Username: os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: true,
		Debugf: func(format string, v ...any) {
			fmt.Printf(format+"\n", v...)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Second * 30,
		MaxOpenConns:         5,
		MaxIdleConns:         5,
		ConnMaxLifetime:      time.Duration(10) * time.Minute,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "smoltx", Version: "beta"},
			},
		},
	})
	if err != nil {
		return err
	}

	connSql := clickhouse.OpenDB(&clickhouse.Options{
		Addr: fmt.Sprintf("%s:%d", os.Getenv("CLICKHOUSE_HOST"), os.Getenv("CLICKHOUSE_PORT")),
		Auth: clickhouse.Auth{
			Database: os.Getenv("CLICKHOUSE_DATABASE"),
			Username: os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		},
		TLS: &tls.Config{
			InsecureSkipVerify: true,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: time.Second * 30,
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		Debug:                true,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
		ClientInfo: clickhouse.ClientInfo{ // optional, please see Client info section in the README.md
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "smoltx", Version: "beta"},
			},
		},
	})
	connSql.SetMaxIdleConns(5)
	connSql.SetMaxOpenConns(10)
	connSql.SetConnMaxLifetime(time.Hour)

	log.Info().
		Msg("Database opened")

	backend := &ClickhouseBackend{
		db:    conn,
		sqlDb: connSql,
	}

	if !readonly {
		if err := backend.Migrate(); err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}
	}

	return backend, nil
}

func (d *SqliteBackend) Close() error {
	return d.db.Close()
}

func (d *SqliteBackend) Migrate() error {
	goose.SetBaseFS(embeddedMigrations)
	if err := goose.SetDialect("clickhouse"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(d.db, "migrations"); err != nil {
		return fmt.Errorf("failed to run goose up: %w", err)
	}
	return nil
}
