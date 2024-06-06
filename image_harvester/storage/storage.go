package storage

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/image-downloader/image_harvester/config"
	_ "github.com/lib/pq"
)

func ConnectDB(cfg config.Config) (*sql.DB, error) {
	dbHost := cfg.PostgresHost
	dbPort, _ := strconv.Atoi(cfg.PostgresPort)
	dbUser := cfg.PostgresUser
	dbPassword := cfg.PostgresPassword
	dbName := cfg.PostgresDB

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SaveImageToDB(db *sql.DB, filename string, query string, data []byte) error {
	queryStr := `
    INSERT INTO images (file_name, query, image_data)
    VALUES ($1, $2, $3)
    `
	_, err := db.Exec(queryStr, filename, query, data)
	if err != nil {
		return err
	}
	return nil
}

func CheckAndCreateTable(db *sql.DB) error {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS images (
        id SERIAL PRIMARY KEY,
        filename TEXT,
        query TEXT,
        data BYTEA
    )
    `
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return err
	}
	return nil
}
