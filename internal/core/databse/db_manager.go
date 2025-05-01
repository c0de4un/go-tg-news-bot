package database

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"net/url"
	"strings"
	"sync"
)

type ConnectionProvider struct {
	dbConnection *sql.DB
}

var (
	dbManagerInstance *ConnectionProvider
	dbManagerOnce     sync.Once
	initError         error // Track initialization error
)

func (dbm *ConnectionProvider) GetDBConnection() *sql.DB {
	return dbm.dbConnection
}

func InitializeDB() error {
	dbManagerOnce.Do(func() {
		var dbCreds *dbConnectionConfig
		dbCreds, initError = loadDBConfig("config/db.xml")
		if initError != nil {
			initError = fmt.Errorf("failed to load DB config: %w", initError)
			return
		}

		dbInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			dbCreds.User, dbCreds.Password, dbCreds.Host, dbCreds.DBName)

		err := ensureDBExists(dbInfo)
		if err != nil {
			initError = fmt.Errorf("failed to create DB: %w", err)
			return
		}

		dbCon, err := sql.Open("postgres", dbInfo)
		if err != nil {
			initError = fmt.Errorf("failed to open DB connection: %w", err)
			return
		}

		if err = dbCon.Ping(); err != nil {
			_ = dbCon.Close()
			initError = fmt.Errorf("failed to ping DB: %w", err)
			return
		}

		err = goose.Up(dbCon, "./db/migrations")
		if err != nil {
			_ = dbCon.Close()
			initError = fmt.Errorf("failed to run migrations: %w", err)
			return
		}

		dbManagerInstance = &ConnectionProvider{dbConnection: dbCon}
		fmt.Println("DB connection initialized successfully")
	})

	return initError
}

func GetDBManager() (*ConnectionProvider, error) {
	if dbManagerInstance == nil || initError != nil {
		return nil, fmt.Errorf("database not initialized: %w", initError)
	}
	return dbManagerInstance, nil
}

func TerminateDBManager() error {
	if dbManagerInstance == nil {
		return nil
	}

	if err := dbManagerInstance.dbConnection.Close(); err != nil {
		return fmt.Errorf("failed to close DB connection: %w", err)
	}

	dbManagerInstance = nil
	// Reset the sync.Once to allow reinitialization
	dbManagerOnce = sync.Once{}
	return nil
}

func ensureDBExists(connStr string) error {
	parsedURL, err := url.Parse(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Extract database name
	dbName := strings.TrimPrefix(parsedURL.Path, "/")
	if dbName == "" {
		return fmt.Errorf("database name is required")
	}

	// Connect to default 'postgres' database
	parsedURL.Path = "/postgres"
	defaultConnStr := parsedURL.String()
	db, err := sql.Open("postgres", defaultConnStr)
	if err != nil {
		return fmt.Errorf("failed to connect to default database: %w", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	// Check database existence
	var exists bool
	err = db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`,
		dbName,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("database check failed: %w", err)
	}

	// Create database if needed
	if !exists {
		// Safely quote database name to prevent SQL injection
		quotedDBName := pq.QuoteIdentifier(dbName)
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", quotedDBName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}
