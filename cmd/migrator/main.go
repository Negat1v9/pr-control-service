package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Negat1v9/pr-review-service/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// load config
	cfg, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Printf("[error] load config: %v", err)
		os.Exit(1)
	}

	coreDSN := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PostgresConfig.DbUser,
		cfg.PostgresConfig.DbPassword,
		cfg.PostgresConfig.DbHost,
		cfg.PostgresConfig.DbPort,
		cfg.PostgresConfig.DbName,
		cfg.PostgresConfig.DbSslMode,
	)
	log.Printf("[debug] DSN: %s", maskedDSN(coreDSN))

	m, err := migrate.New("file://migrations", coreDSN)
	if err != nil {
		log.Printf("[error] migrate.New: %v", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("[warn] no migrations to apply")
		} else {
			log.Printf("[error] migrate.Up: %v", err)
			os.Exit(1)
		}
	}

	m.Close()
	log.Printf("[info] migrations applied successfully")

}

func maskedDSN(dsn string) string {
	if len(dsn) > 20 {
		return dsn[:20] + "***"
	}
	return dsn
}
