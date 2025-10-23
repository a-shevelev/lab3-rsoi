package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ConnStr string `envconfig:"CONN_STRING" required:"true"`
}

type Client interface {
	Ping(ctx context.Context) error
	Conn() Connection
	Close()
}

func (c *dbClient) Conn() Connection {
	return c.conn
}

type dbClient struct {
	conn   *pgxpool.Pool
	cancel context.CancelFunc
	ctx    context.Context
}

func (c *dbClient) Ping(ctx context.Context) error {
	return c.conn.Ping(ctx)
}

func Connect(ctx context.Context, cfg Config) (*dbClient, error) {
	pool, err := pgxpool.New(ctx, cfg.ConnStr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	log.Info("Connecting to Postgres...")

	return &dbClient{
		conn:   pool,
		cancel: cancel,
		ctx:    ctx,
	}, nil
}

func (c *dbClient) Close() {
	if c.cancel != nil {
		c.cancel()
	}
	c.conn.Close()
}

type Connection interface {
	Query(ctx context.Context, query string, values ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, values ...any) pgx.Row
	Exec(ctx context.Context, query string, values ...any) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
}
