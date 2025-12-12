package database

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func OpenDB(dsn string) (*sql.DB, error) {
	// Parse DSN to extract database name
	parsedDSN, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	dbName := strings.TrimPrefix(parsedDSN.Path, "/")
	if dbName == "" {
		return nil, fmt.Errorf("database name not found in DSN")
	}

	// Create DSN for connecting to default 'postgres' database
	defaultDSN := *parsedDSN
	defaultDSN.Path = "/postgres"

	// Connect to default database to create target database
	defaultDB, err := sql.Open("postgres", defaultDSN.String())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres database: %w", err)
	}
	defer defaultDB.Close()

	// Check if database exists, create if not
	var exists bool
	err = defaultDB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		log.Printf("Creating database '%s'...", dbName)
		_, err = defaultDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return nil, fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", dbName)
	}

	// Connect to target database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	migrationsDir := "internal/migrations"

	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		filePath := filepath.Join(migrationsDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		log.Printf("Running migration: %s", file.Name())
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
		}
	}

	log.Println("Migrations completed successfully")
	return nil
}
