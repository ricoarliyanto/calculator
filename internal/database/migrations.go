package database

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// Migration represents a database migration
type Migration struct {
	ID      string
	Name    string
	Query   string
	Version int
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db       *Database
	migrations []Migration
	logger    *logrus.Logger
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *Database) *MigrationManager {
	return &MigrationManager{
		db:    db,
		logger: logrus.New(),
	}
}

// AddMigration adds a migration to the manager
func (m *MigrationManager) AddMigration(id, name, query string, version int) {
	m.migrations = append(m.migrations, Migration{
		ID:      id,
		Name:    name,
		Query:   query,
		Version: version,
	})
	// Sort by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})
}

// GetCurrentVersion returns the current migration version from the database
func (m *MigrationManager) GetCurrentVersion() (int, error) {
	conn, err := m.db.GetDB()
	if err != nil {
		return 0, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Create migrations table if it doesn't exist
	if _, err := conn.Exec(`CREATE TABLE IF NOT EXISTS migrations (id TEXT PRIMARY KEY, name TEXT, version INTEGER, applied_at TIMESTAMP)`); err != nil {
		return 0, fmt.Errorf("failed to create migrations table: %w", err)
	}

	var version int
	if err := conn.QueryRow("SELECT COALESCE(MAX(version), 0) FROM migrations").Scan(&version); err != nil {
		return 0, fmt.Errorf("failed to query current version: %w", err)
	}

	return version, nil
}

// ApplyMigrations applies all pending migrations
func (m *MigrationManager) ApplyMigrations() error {
	conn, err := m.db.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Create migrations table
	if _, err := conn.Exec(`CREATE TABLE IF NOT EXISTS migrations (id TEXT PRIMARY KEY, name TEXT, version INTEGER, applied_at TIMESTAMP)`); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	currentVersion, err := m.GetCurrentVersion()
	if err != nil {
		return err
	}

	m.logger.WithField("current_version", currentVersion).Info("Applying migrations")

	for _, migration := range m.migrations {
		if migration.Version <= currentVersion {
			continue
		}

		m.logger.WithFields(logrus.Fields{
			"migration_id": migration.ID,
			"migration_name": migration.Name,
			"version": migration.Version,
		}).Info("Applying migration")

		tx, err := conn.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		// Execute migration query
		if _, err := tx.Exec(migration.Query); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration.ID, err)
		}

		// Record migration
		if _, err := tx.Exec(
			"INSERT INTO migrations (id, name, version, applied_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)",
			migration.ID,
			migration.Name,
			migration.Version,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.ID, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.ID, err)
		}

		m.logger.WithField("migration_id", migration.ID).Info("Migration applied successfully")
	}

	return nil
}

// GetAppliedMigrations returns a list of applied migrations
func (m *MigrationManager) GetAppliedMigrations() ([]Migration, error) {
	conn, err := m.db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	rows, err := conn.Query("SELECT id, name, version FROM migrations ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	var migrations []Migration
	for rows.Next() {
		migration := Migration{}
		if err := rows.Scan(&migration.ID, &migration.Name, &migration.Version); err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

// ValidateMigrations validates that all required migrations are applied
func (m *MigrationManager) ValidateMigrations() error {
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Create a map of applied migration IDs
	appliedMap := make(map[string]bool)
	for _, migration := range applied {
		appliedMap[migration.ID] = true
	}

	// Check if all registered migrations are applied
	var missing []string
	for _, migration := range m.migrations {
		if !appliedMap[migration.ID] {
			missing = append(missing, migration.ID)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing migrations: %s", strings.Join(missing, ", "))
	}

	return nil
}
