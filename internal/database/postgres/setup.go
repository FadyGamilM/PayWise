package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"paywise/config"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
)

func Setup() (*sql.DB, error) {
	// configs, err := config.LoadPostgresConfig("./config")
	configs, err := config.LoadPostgresConfig_v2()
	if err != nil {
		fmt.Println("error trying to load config variables", err)
		return nil, err
	}
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

	log.Println(dsn.String())
	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		fmt.Println("error trying to open a postgres connection", err)
		return nil, err
	}

	log.Println("connecting to a database successfully ..")
	return db, nil
}

func SetupTest() (*sql.DB, error) {
	configs, err := config.LoadPostgresTestConfig()
	if err != nil {
		fmt.Println("error trying to load config variables", err)
		return nil, err
	}
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
