package database

import (
	"database/sql"
	"fmt"
)

// SettingsRepository handles settings database operations
type SettingsRepository struct {
	db *Database
}

// NewSettingsRepository creates a new settings repository
func NewSettingsRepository(db *Database) *SettingsRepository {
	return &SettingsRepository{db: db}
}

// SaveSetting saves or updates a setting
func (r *SettingsRepository) SaveSetting(key, value string) error {
	conn, err := r.db.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Check if setting exists
	var existingID int64
	err = conn.QueryRow("SELECT id FROM settings WHERE key = ?", key).Scan(&existingID)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to check setting existence: %w", err)
	}

	if existingID > 0 {
		// Update existing setting
		if _, err := conn.Exec("UPDATE settings SET value = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", value, existingID); err != nil {
			return fmt.Errorf("failed to update setting: %w", err)
		}
	} else {
		// Insert new setting
		if _, err := conn.Exec("INSERT INTO settings (key, value) VALUES (?, ?)", key, value); err != nil {
			return fmt.Errorf("failed to insert setting: %w", err)
		}
	}

	return nil
}

// GetSetting retrieves a setting by key
func (r *SettingsRepository) GetSetting(key string) (*Setting, error) {
	conn, err := r.db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	setting := &Setting{}
	err = conn.QueryRow(
		"SELECT id, key, value, created_at, updated_at FROM settings WHERE key = ?",
		key,
	).Scan(&setting.ID, &setting.Key, &setting.Value, &setting.CreatedAt, &setting.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query setting: %w", err)
	}

	return setting, nil
}

// GetSettingValue retrieves just the value of a setting by key
func (r *SettingsRepository) GetSettingValue(key, defaultValue string) (string, error) {
	setting, err := r.GetSetting(key)
	if err != nil {
		return "", fmt.Errorf("failed to get setting: %w", err)
	}

	if setting == nil {
		return defaultValue, nil
	}

	return setting.Value, nil
}

// GetAllSettings retrieves all settings
func (r *SettingsRepository) GetAllSettings() ([]Setting, error) {
	conn, err := r.db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	rows, err := conn.Query("SELECT id, key, value, created_at, updated_at FROM settings ORDER BY key")
	if err != nil {
		return nil, fmt.Errorf("failed to query settings: %w", err)
	}
	defer rows.Close()

	var settings []Setting
	for rows.Next() {
		setting := Setting{}
		if err := rows.Scan(&setting.ID, &setting.Key, &setting.Value, &setting.CreatedAt, &setting.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan setting: %w", err)
		}
		settings = append(settings, setting)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating settings: %w", err)
	}

	return settings, nil
}

// DeleteSetting deletes a setting by key
func (r *SettingsRepository) DeleteSetting(key string) error {
	conn, err := r.db.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if _, err := conn.Exec("DELETE FROM settings WHERE key = ?", key); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	return nil
}

// GetIntSetting retrieves an integer setting with default value
func (r *SettingsRepository) GetIntSetting(key string, defaultValue int) (int, error) {
	value, err := r.GetSettingValue(key, fmt.Sprintf("%d", defaultValue))
	if err != nil {
		return 0, err
	}

	var result int
	_, err = fmt.Sscanf(value, "%d", &result)
	if err != nil {
		return defaultValue, nil
	}

	return result, nil
}

// GetBoolSetting retrieves a boolean setting with default value
func (r *SettingsRepository) GetBoolSetting(key string, defaultValue bool) (bool, error) {
	value, err := r.GetSettingValue(key, fmt.Sprintf("%t", defaultValue))
	if err != nil {
		return false, err
	}

	if value == "true" || value == "1" {
		return true, nil
	}
	if value == "false" || value == "0" {
		return false, nil
	}

	return defaultValue, nil
}
