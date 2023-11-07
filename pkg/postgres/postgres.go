package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type DB interface {
	BeginTx(ctx context.Context) (*pgxpool.Conn, pgx.Tx, error)
	EndTx(context.Context, pgx.Tx) error
}

type PostgresDB struct {
	*pgxpool.Pool
}

func (db *PostgresDB) BeginTx(ctx context.Context) (*pgxpool.Conn, pgx.Tx, error) {
	conn, err := db.Acquire(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot acquire connection from dbpool")
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "cannot begin transaction")
	}
	return conn, tx, nil
}

func (db *PostgresDB) EndTx(ctx context.Context, tx pgx.Tx) error {
	err := tx.Commit(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot commit transaction")
	}
	return nil
}

func NewPostgresDB(ctx context.Context, config string, tracer pgx.QueryTracer) (*PostgresDB, error) {
	pool, err := NewPool(ctx, config, tracer)
	if err != nil {
		return nil, err
	}
	return &PostgresDB{
		Pool: pool,
	}, nil
}

func NewPool(ctx context.Context, config string, tracer pgx.QueryTracer) (*pgxpool.Pool, error) {
	c, err := pgxpool.ParseConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse postgres pool config: %w", err)
	}
	c.ConnConfig.Tracer = tracer
	db, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("cannot create postgres pool: %w", err)
	}
	return db, nil
}

func NewConn(ctx context.Context, config string, tracer pgx.QueryTracer) (*pgx.Conn, error) {
	c, err := pgx.ParseConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse postgres config: %w", err)
	}
	c.Tracer = tracer
	db, err := pgx.ConnectConfig(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("cannot establish a connection with a PostgreSQL server: %w", err)
	}
	return db, nil
}
