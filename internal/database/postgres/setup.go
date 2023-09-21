package postgres

import (
	"database/sql"
	"fmt"
	"net/url"
	"paywise/config"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

func Setup() (*sql.DB, error) {
	configs, err := config.LoadPostgresConfig()

	// construct the conn string
	dsn := url.URL{
		Scheme: "postgres",
		Host:   configs.Postgresdb.Host,
		User:   url.UserPassword(configs.Postgresdb.User, configs.Postgresdb.Password),
		Path:   configs.Postgresdb.Dbname,
	}

	q := dsn.Query()
	q.Add("sslmode", configs.Postgresdb.Sslmode)

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		fmt.Println("error trying to open a postgres connection", err)
		return nil, err
	}

	return db, nil
}

func SetupTest() (*sql.DB, error) {
	configs, err := config.LoadPostgresTestConfig()

	// construct the conn string
	dsn := url.URL{
		Scheme: "postgres",
		Host:   configs.Testpostgresdb.Host,
		User:   url.UserPassword(configs.Testpostgresdb.User, configs.Testpostgresdb.Password),
		Path:   configs.Testpostgresdb.Dbname,
	}

	q := dsn.Query()
	q.Add("sslmode", configs.Testpostgresdb.Sslmode)

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		fmt.Println("error trying to open a postgres connection", err)
		return nil, err
	}

	return db, nil
}
