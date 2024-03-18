package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	C "github.com/gabriel-tama/projectsprint-socmed/common/config"
)

type DB struct {
	Pool *pgxpool.Pool
}

func GetPostgresURL() string {
	env, err := C.Get()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := env.DBHost
	dbPort := env.DBPort
	dbUser := env.DBUser
	dbPass := env.DBPassword
	dbName := env.DBName

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPass,
		dbHost, dbPort, dbName)
}

func Config(DATABASE_URL string) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Your own Database URL
	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection pool to the database!!")
	}

	return dbConfig
}

func Init(ctx context.Context) (*DB, error) {
	PgPool, err := pgxpool.NewWithConfig(context.Background(), Config(GetPostgresURL()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		return nil, err
	}

	return &DB{Pool: PgPool}, nil
}

func (db *DB) DB() *DB {

	return db
}

func (db *DB) StartTx(ctx context.Context, f func(pgx.Tx) error) error {
	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return err
	}

	err = f(tx)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func (db *DB) Close(ctx context.Context) {
	db.Pool.Close()
}
