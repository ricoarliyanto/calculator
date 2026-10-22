package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

var (
	ErrDatabaseNotConnected = errors.New("database connection not established")
	ErrMigrationFailed      = errors.New("database migration failed")
)

// Database represents the database connection manager
type Database struct {
	conn     *sql.DB
	dataDir  string
	dbPath   string
	logger   *logrus.Logger
	once     sync.Once
	initErr  error
}

// NewDatabase creates a new database manager instance
func NewDatabase(dataDir string) *Database {
	db := &Database{
		dataDir: dataDir,
		logger: logrus.New(),
	}
	return db
}

// Connect establishes database connection
func (db *Database) Connect() error {
	db.once.Do(func() {
		// Create data directory if it doesn't exist
		if err := os.MkdirAll(db.dataDir, 0755); err != nil {
			db.initErr = fmt.Errorf("failed to create data directory: %w", err)
			return
		}

		// Set database path
		db.dbPath = filepath.Join(db.dataDir, "calculator.db")

		// Open database connection
		conn, err := sql.Open("sqlite3", db.dbPath)
		if err != nil {
			db.initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		// Configure connection pool
		conn.SetMaxOpenConns(10)
		conn.SetMaxIdleConns(5)
		conn.SetConnMaxLifetime(30 * 60) // 30 minutes

		// Verify connection
		if err := conn.Ping(); err != nil {
			db.initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		db.conn = conn
		db.logger.Info("Database connected successfully")
	})

	return db.initErr
}

// GetDB returns the underlying *sql.DB connection
func (db *Database) GetDB() (*sql.DB, error) {
	if db.conn == nil {
		if err := db.Connect(); err != nil {
			return nil, ErrDatabaseNotConnected
		}
	}
	return db.conn, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.conn != nil {
		if err := db.conn.Close(); err != nil {
			return fmt.Errorf("failed to close database: %w", err)
		}
		db.logger.Info("Database connection closed")
	}
	return nil
}

// GetDBPath returns the database file path
func (db *Database) GetDBPath() string {
	return db.dbPath
}

// InitializeDatabase initializes the database with required tables
// This is now handled by RunMigrations which uses the migration system
func (db *Database) InitializeDatabase() error {
	return db.RunMigrations()
}

// RunMigrations runs database migrations
func (db *Database) RunMigrations() error {
	migrationManager := NewMigrationManager(db)

	// Register initial schema migrations
	migrationManager.AddMigration(
		"001_initial_schema",
		"Initial database schema",
		`CREATE TABLE IF NOT EXISTS calculations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			expression TEXT NOT NULL,
			result TEXT NOT NULL,
			mode TEXT NOT NULL DEFAULT 'standard',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS settings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key TEXT NOT NULL UNIQUE,
			value TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_calculations_mode ON calculations(mode);
		CREATE INDEX IF NOT EXISTS idx_calculations_created ON calculations(created_at);
		CREATE INDEX IF NOT EXISTS idx_settings_key ON settings(key);`,
		1,
	)

	// Run migrations
	if err := migrationManager.ApplyMigrations(); err != nil {
		db.logger.WithError(err).Error("Failed to apply migrations")
		return fmt.Errorf("%w: %w", ErrMigrationFailed, err)
	}

	// Insert default settings if they don't exist
	conn, err := db.GetDB()
	if err != nil {
		return err
	}

	defaultSettings := []struct {
		key   string
		value string
	}{
		{"default_mode", "standard"},
		{"history_limit", "1000"},
		{"decimal_precision", "10"},
		{"theme", "light"},
	}

	for _, setting := range defaultSettings {
		var count int
		if err := conn.QueryRow("SELECT COUNT(*) FROM settings WHERE key = ?", setting.key).Scan(&count); err != nil {
			db.logger.WithError(err).Warn("Failed to check setting existence")
			continue
		}

		if count == 0 {
			if _, err := conn.Exec("INSERT INTO settings (key, value) VALUES (?, ?)", setting.key, setting.value); err != nil {
				db.logger.WithError(err).Warn("Failed to insert default setting")
			}
		}
	}

	db.logger.Info("Database initialized successfully with migrations")
	return nil
}
