package storage

import (
	"database/sql"
	"fmt"

	"github.com/andrii-stp/users-crud/config"
	_ "github.com/lib/pq"
)

func Connect(cfg *config.Database) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	db, err := sql.Open(cfg.Driver, connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func InitDB(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		user_id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
		user_name VARCHAR(50) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		user_status VARCHAR(1) NOT NULL,
		department VARCHAR(255)
	  );
	`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	return nil
}
