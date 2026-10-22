package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Mode           string `yaml:"mode" json:"mode"`
	AngleUnit      string `yaml:"angle-unit" json:"angle_unit"`
	DecimalPlaces  int    `yaml:"decimal-places" json:"decimal_places"`
	HistoryLimit   int    `yaml:"history-limit" json:"history_limit"`
	AutoSave       bool   `yaml:"auto-save" json:"auto_save"`
	ShowTimestamps bool   `yaml:"show-timestamps" json:"show_timestamps"`
	Theme          string `yaml:"theme" json:"theme"`
}

// DefaultConfig returns a new configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Mode:           "standard",
		AngleUnit:      "rad",
		DecimalPlaces:  10,
		HistoryLimit:   1000,
		AutoSave:       true,
		ShowTimestamps: true,
		Theme:          "light",
	}
}

// ConfigManager manages configuration loading and saving
type ConfigManager struct {
	config     *Config
	configPath string
	mu         sync.RWMutex
}

// NewConfigManager creates a new config manager
func NewConfigManager(configDir string) *ConfigManager {
	return &ConfigManager{
		config:     DefaultConfig(),
		configPath: filepath.Join(configDir, "config.yaml"),
	}
}

// Load loads configuration from file
func (cm *ConfigManager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Try to read YAML first
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		// If YAML doesn't exist, try JSON
		jsonPath := filepath.Join(filepath.Dir(cm.configPath), "config.json")
		data, err = os.ReadFile(jsonPath)
		if err != nil {
			// No config file exists, use defaults
			cm.config = DefaultConfig()
			return nil
		}
		// Load from JSON
		if err := json.Unmarshal(data, &cm.config); err != nil {
			return fmt.Errorf("failed to parse JSON config: %w", err)
		}
		return nil
	}

	// Load from YAML
	if err := yaml.Unmarshal(data, &cm.config); err != nil {
		return fmt.Errorf("failed to parse YAML config: %w", err)
	}

	return nil
}

// Save saves configuration to file
func (cm *ConfigManager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Create directory if it doesn't exist
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(cm.config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (cm *ConfigManager) GetConfig() *Config {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config
}

// SetMode sets the calculator mode
func (cm *ConfigManager) SetMode(mode string) error {
	if mode != "standard" && mode != "scientific" {
		return fmt.Errorf("invalid mode: %s (must be 'standard' or 'scientific')", mode)
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.Mode = mode
	return nil
}

// GetMode returns the calculator mode
func (cm *ConfigManager) GetMode() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.Mode
}

// SetAngleUnit sets the angle unit
func (cm *ConfigManager) SetAngleUnit(unit string) error {
	if unit != "rad" && unit != "deg" {
		return fmt.Errorf("invalid angle unit: %s (must be 'rad' or 'deg')", unit)
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.AngleUnit = unit
	return nil
}

// GetAngleUnit returns the angle unit
func (cm *ConfigManager) GetAngleUnit() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.AngleUnit
}

// SetDecimalPlaces sets the number of decimal places
func (cm *ConfigManager) SetDecimalPlaces(places int) error {
	if places < 0 || places > 20 {
		return fmt.Errorf("invalid decimal places: %d (must be between 0 and 20)", places)
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.DecimalPlaces = places
	return nil
}

// GetDecimalPlaces returns the number of decimal places
func (cm *ConfigManager) GetDecimalPlaces() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.DecimalPlaces
}

// SetHistoryLimit sets the maximum number of history entries
func (cm *ConfigManager) SetHistoryLimit(limit int) error {
	if limit < 0 || limit > 100000 {
		return fmt.Errorf("invalid history limit: %d (must be between 0 and 100000)", limit)
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.HistoryLimit = limit
	return nil
}

// GetHistoryLimit returns the maximum number of history entries
func (cm *ConfigManager) GetHistoryLimit() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.HistoryLimit
}

// SetTheme sets the UI theme
func (cm *ConfigManager) SetTheme(theme string) error {
	if theme != "light" && theme != "dark" {
		return fmt.Errorf("invalid theme: %s (must be 'light' or 'dark')", theme)
	}
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.Theme = theme
	return nil
}

// GetTheme returns the UI theme
func (cm *ConfigManager) GetTheme() string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.Theme
}

// SetAutoSave sets whether to auto-save history
func (cm *ConfigManager) SetAutoSave(autoSave bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.AutoSave = autoSave
}

// GetAutoSave returns whether to auto-save history
func (cm *ConfigManager) GetAutoSave() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.AutoSave
}

// SetShowTimestamps sets whether to show timestamps in history
func (cm *ConfigManager) SetShowTimestamps(show bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config.ShowTimestamps = show
}

// GetShowTimestamps returns whether to show timestamps in history
func (cm *ConfigManager) GetShowTimestamps() bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.config.ShowTimestamps
}

// Validate validates the configuration
func (cm *ConfigManager) Validate() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	config := cm.config

	if config.Mode != "standard" && config.Mode != "scientific" {
		return fmt.Errorf("invalid mode: %s", config.Mode)
	}

	if config.AngleUnit != "rad" && config.AngleUnit != "deg" {
		return fmt.Errorf("invalid angle unit: %s", config.AngleUnit)
	}

	if config.DecimalPlaces < 0 || config.DecimalPlaces > 20 {
		return fmt.Errorf("invalid decimal places: %d", config.DecimalPlaces)
	}

	if config.HistoryLimit < 0 || config.HistoryLimit > 100000 {
		return fmt.Errorf("invalid history limit: %d", config.HistoryLimit)
	}

	if config.Theme != "light" && config.Theme != "dark" {
		return fmt.Errorf("invalid theme: %s", config.Theme)
	}

	return nil
}

// GetConfigPath returns the path to the config file
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}

// ConfigFileExists checks if the config file exists
func (cm *ConfigManager) ConfigFileExists() bool {
	_, err := os.Stat(cm.configPath)
	return err == nil
}

// EnsureConfigFile ensures the config file exists with default values
func (cm *ConfigManager) EnsureConfigFile() error {
	if cm.ConfigFileExists() {
		return nil
	}

	// Create directory
	dir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Save default config
	return cm.Save()
}

// ExportConfig exports configuration to a map for display
func (cm *ConfigManager) ExportConfig() map[string]interface{} {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return map[string]interface{}{
		"mode":            cm.config.Mode,
		"angle_unit":      cm.config.AngleUnit,
		"decimal_places":  cm.config.DecimalPlaces,
		"history_limit":   cm.config.HistoryLimit,
		"auto_save":       cm.config.AutoSave,
		"show_timestamps": cm.config.ShowTimestamps,
		"theme":           cm.config.Theme,
	}
}

// ImportConfig imports configuration from a map
func (cm *ConfigManager) ImportConfig(data map[string]interface{}) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if mode, ok := data["mode"].(string); ok {
		cm.config.Mode = mode
	}

	if unit, ok := data["angle_unit"].(string); ok {
		cm.config.AngleUnit = unit
	}

	if places, ok := data["decimal_places"].(float64); ok {
		cm.config.DecimalPlaces = int(places)
	}

	if limit, ok := data["history_limit"].(float64); ok {
		cm.config.HistoryLimit = int(limit)
	}

	if autoSave, ok := data["auto_save"].(bool); ok {
		cm.config.AutoSave = autoSave
	}

	if showTimestamps, ok := data["show_timestamps"].(bool); ok {
		cm.config.ShowTimestamps = showTimestamps
	}

	if theme, ok := data["theme"].(string); ok {
		cm.config.Theme = theme
	}

	return nil
}

// Reload reloads the configuration from disk
func (cm *ConfigManager) Reload() error {
	return cm.Load()
}

// Reset resets the configuration to default values
func (cm *ConfigManager) Reset() {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.config = DefaultConfig()
}

// GetDataDir returns the directory containing the config file
func (cm *ConfigManager) GetDataDir() string {
	return filepath.Dir(cm.configPath)
}

// Global config manager instance
var (
	globalConfig *ConfigManager
	configOnce   sync.Once
)

// GetGlobalConfig returns the global config manager
func GetGlobalConfig() *ConfigManager {
	configOnce.Do(func() {
		globalConfig = NewConfigManager(getDefaultConfigDir())
	})
	return globalConfig
}

// getDefaultConfigDir returns the default configuration directory
func getDefaultConfigDir() string {
	// Try APPDATA on Windows
	if appData := os.Getenv("APPDATA"); appData != "" {
		return filepath.Join(appData, "calculator")
	}

	// Try HOME on Unix-like
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".config", "calculator")
	}

	// Fallback to current directory
	return "."
}

// InitGlobalConfig initializes the global config manager with a specific directory
func InitGlobalConfig(configDir string) *ConfigManager {
	configOnce.Do(func() {
		globalConfig = NewConfigManager(configDir)
	})
	return globalConfig
}
