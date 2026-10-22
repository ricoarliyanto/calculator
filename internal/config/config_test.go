package config

import (
	"os"
	"path/filepath"
	"testing"
)

// =============================================================================
// Config Default Tests
// =============================================================================

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Mode != ModeStandard {
		t.Errorf("Expected default mode 'standard', got '%s'", cfg.Mode)
	}
	if cfg.AngleUnit != AngleUnitRad {
		t.Errorf("Expected default angle unit 'rad', got '%s'", cfg.AngleUnit)
	}
	if cfg.DecimalPlaces != 10 {
		t.Errorf("Expected default decimal places 10, got %d", cfg.DecimalPlaces)
	}
	if cfg.HistoryLimit != 1000 {
		t.Errorf("Expected default history limit 1000, got %d", cfg.HistoryLimit)
	}
	if !cfg.AutoSave {
		t.Error("Expected default auto save true")
	}
	if !cfg.ShowTimestamps {
		t.Error("Expected default show timestamps true")
	}
	if cfg.Theme != "light" {
		t.Errorf("Expected default theme 'light', got '%s'", cfg.Theme)
	}
}

// Test constants for config package
const (
	ModeStandard     = "standard"
	ModeScientific   = "scientific"
	AngleUnitRad     = "rad"
	AngleUnitDeg     = "deg"
)

// =============================================================================
// Config Manager Tests
// =============================================================================

func TestNewConfigManager(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if cm == nil {
		t.Fatal("NewConfigManager() returned nil")
	}
}

func TestConfigLoadDefaults(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	err := cm.Load()
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}
	cfg := cm.GetConfig()
	if cfg.Mode != ModeStandard {
		t.Errorf("Expected default mode, got '%s'", cfg.Mode)
	}
}

func TestConfigSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	// Change settings
	if err := cm.SetMode(ModeScientific); err != nil {
		t.Fatalf("SetMode() error = %v", err)
	}
	if err := cm.SetDecimalPlaces(5); err != nil {
		t.Fatalf("SetDecimalPlaces() error = %v", err)
	}
	if err := cm.SetHistoryLimit(500); err != nil {
		t.Fatalf("SetHistoryLimit() error = %v", err)
	}
	if err := cm.SetTheme("dark"); err != nil {
		t.Fatalf("SetTheme() error = %v", err)
	}
	cm.SetAutoSave(true)
	cm.SetShowTimestamps(true)
	if err := cm.SetAngleUnit(AngleUnitDeg); err != nil {
		t.Fatalf("SetAngleUnit() error = %v", err)
	}

	if err := cm.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Check file exists
	configPath := filepath.Join(dir, "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should exist after Save()")
	}

	// Create new manager to load from file
	cm2 := NewConfigManager(dir)
	err := cm2.Load()
	if err != nil {
		t.Errorf("Load() error = %v", err)
	}
	cfg := cm2.GetConfig()
	if cfg.Mode != ModeScientific {
		t.Errorf("Expected mode 'scientific', got '%s'", cfg.Mode)
	}
	if cfg.DecimalPlaces != 5 {
		t.Errorf("Expected decimal places 5, got %d", cfg.DecimalPlaces)
	}
	if cfg.HistoryLimit != 500 {
		t.Errorf("Expected history limit 500, got %d", cfg.HistoryLimit)
	}
	if cfg.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got '%s'", cfg.Theme)
	}
	if !cfg.AutoSave {
		t.Error("Expected auto save to be true")
	}
	if !cfg.ShowTimestamps {
		t.Error("Expected show timestamps to be true")
	}
	if cfg.AngleUnit != AngleUnitDeg {
		t.Errorf("Expected angle unit 'deg', got '%s'", cfg.AngleUnit)
	}
}

func TestConfigLoadWithExistingFile(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	// Save config
	if err := cm.SetMode(ModeScientific); err != nil {
		t.Fatalf("SetMode() error = %v", err)
	}
	if err := cm.SetDecimalPlaces(8); err != nil {
		t.Fatalf("SetDecimalPlaces() error = %v", err)
	}
	if err := cm.Save(); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Load config on same manager
	err := cm.Load()
	if err != nil {
		t.Errorf("Reload() error = %v", err)
	}
	cfg := cm.GetConfig()
	if cfg.Mode != ModeScientific {
		t.Errorf("Expected mode 'scientific', got '%s'", cfg.Mode)
	}
	if cfg.DecimalPlaces != 8 {
		t.Errorf("Expected decimal places 8, got %d", cfg.DecimalPlaces)
	}
}

func TestConfigReset(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	// Change settings
	_ = cm.SetMode(ModeScientific)
	_ = cm.SetDecimalPlaces(2)
	_ = cm.SetHistoryLimit(100)
	_ = cm.SetTheme("dark")
	cm.SetAutoSave(false)
	cm.SetShowTimestamps(false)
	_ = cm.SetAngleUnit(AngleUnitDeg)

	// Reset
	cm.Reset()

	// Verify reset to defaults
	if cm.GetMode() != ModeStandard {
		t.Errorf("Expected mode 'standard' after reset, got '%s'", cm.GetMode())
	}
	if cm.GetDecimalPlaces() != 10 {
		t.Errorf("Expected decimal places 10 after reset, got %d", cm.GetDecimalPlaces())
	}
	if cm.GetHistoryLimit() != 1000 {
		t.Errorf("Expected history limit 1000 after reset, got %d", cm.GetHistoryLimit())
	}
	if cm.GetTheme() != "light" {
		t.Errorf("Expected theme 'light' after reset, got '%s'", cm.GetTheme())
	}
	if !cm.GetAutoSave() {
		t.Error("Expected auto save true after reset")
	}
	if !cm.GetShowTimestamps() {
		t.Error("Expected show timestamps true after reset")
	}
	if cm.GetAngleUnit() != AngleUnitRad {
		t.Errorf("Expected angle unit 'rad' after reset, got '%s'", cm.GetAngleUnit())
	}
}

// =============================================================================
// Config Getter/Setter Tests
// =============================================================================

func TestConfigModeSettersAndGetters(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetMode(ModeScientific); err != nil {
		t.Errorf("SetMode() error = %v", err)
	}
	if cm.GetMode() != ModeScientific {
		t.Errorf("Expected mode 'scientific', got '%s'", cm.GetMode())
	}

	if err := cm.SetMode(ModeStandard); err != nil {
		t.Errorf("SetMode() error = %v", err)
	}
	if cm.GetMode() != ModeStandard {
		t.Errorf("Expected mode 'standard', got '%s'", cm.GetMode())
	}
}

func TestConfigDecimalPlacesSettersAndGetters(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetDecimalPlaces(20); err != nil {
		t.Errorf("SetDecimalPlaces() error = %v", err)
	}
	if cm.GetDecimalPlaces() != 20 {
		t.Errorf("Expected decimal places 20, got %d", cm.GetDecimalPlaces())
	}
}

func TestConfigHistoryLimitSettersAndGetters(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetHistoryLimit(500); err != nil {
		t.Errorf("SetHistoryLimit() error = %v", err)
	}
	if cm.GetHistoryLimit() != 500 {
		t.Errorf("Expected history limit 500, got %d", cm.GetHistoryLimit())
	}
}

func TestConfigThemeSettersAndGetters(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetTheme("dark"); err != nil {
		t.Errorf("SetTheme() error = %v", err)
	}
	if cm.GetTheme() != "dark" {
		t.Errorf("Expected theme 'dark', got '%s'", cm.GetTheme())
	}

	if err := cm.SetTheme("light"); err != nil {
		t.Errorf("SetTheme() error = %v", err)
	}
	if cm.GetTheme() != "light" {
		t.Errorf("Expected theme 'light', got '%s'", cm.GetTheme())
	}
}

func TestConfigAutoSaveSettersAndGetters(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	cm.SetAutoSave(true)
	if !cm.GetAutoSave() {
		t.Error("Expected auto save true")
	}

	cm.SetAutoSave(false)
	if cm.GetAutoSave() {
		t.Error("Expected auto save false")
	}
}

func TestConfigShowTimestampsSettersAndGetters(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	cm.SetShowTimestamps(true)
	if !cm.GetShowTimestamps() {
		t.Error("Expected show timestamps true")
	}

	cm.SetShowTimestamps(false)
	if cm.GetShowTimestamps() {
		t.Error("Expected show timestamps false")
	}
}

func TestConfigAngleUnitSettersAndGetters(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetAngleUnit(AngleUnitDeg); err != nil {
		t.Errorf("SetAngleUnit() error = %v", err)
	}
	if cm.GetAngleUnit() != AngleUnitDeg {
		t.Errorf("Expected angle unit 'deg', got '%s'", cm.GetAngleUnit())
	}
}

// =============================================================================
// Config Validation Tests
// =============================================================================

func TestConfigValidateValidMode(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetMode(ModeScientific); err != nil {
		t.Errorf("SetMode() error = %v", err)
	}
	if err := cm.Validate(); err != nil {
		t.Errorf("Validate() for mode '%s' error = %v", ModeScientific, err)
	}
}

func TestConfigValidateInvalidMode(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.Validate(); err != nil {
		// Should be valid with default
		t.Errorf("Validate() default should pass, error = %v", err)
	}
	// Now make it invalid
	cm.config.Mode = "invalid_mode"
	if err := cm.Validate(); err == nil {
		t.Error("Validate() for invalid mode should return error")
	}
}

func TestConfigValidateNegativeDecimalPlaces(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetDecimalPlaces(-1); err == nil {
		t.Error("SetDecimalPlaces(-1) should return error")
	}
}

func TestConfigValidateNegativeHistoryLimit(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetHistoryLimit(-1); err == nil {
		t.Error("SetHistoryLimit(-1) should return error")
	}
}

func TestConfigValidateInvalidAngleUnit(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetAngleUnit("invalid"); err == nil {
		t.Error("SetAngleUnit('invalid') should return error")
	}
}

func TestConfigValidateInvalidTheme(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	if err := cm.SetTheme("invalid_theme"); err == nil {
		t.Error("SetTheme('invalid_theme') should return error")
	}
}

func TestConfigSaveToInvalidDir(t *testing.T) {
	cm := NewConfigManager("/root/invalid/path/calculator")
	err := cm.Save()
	if err == nil {
		t.Error("Expected error when saving to invalid directory")
	}
}

func TestConfigLoadFromNonExistentFile(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	// No file exists, should return defaults
	err := cm.Load()
	if err != nil {
		t.Errorf("Load() from non-existent file should return defaults, error = %v", err)
	}
	cfg := cm.GetConfig()
	if cfg.Mode != ModeStandard {
		t.Errorf("Expected default mode, got '%s'", cfg.Mode)
	}
}

func TestConfigFileExists(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	// File should not exist initially
	if cm.ConfigFileExists() {
		t.Error("Config should not exist initially")
	}

	// Save config
	_ = cm.Save()

	// File should exist now
	if !cm.ConfigFileExists() {
		t.Error("Config should exist after Save()")
	}
}

func TestConfigGetConfig(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	cfg := cm.GetConfig()
	if cfg == nil {
		t.Fatal("GetConfig() returned nil")
	}
	if cfg.Mode != ModeStandard {
		t.Errorf("Expected default mode, got '%s'", cfg.Mode)
	}
}

func TestConfigGetConfigPath(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	expectedPath := filepath.Join(dir, "config.yaml")
	if cm.GetConfigPath() != expectedPath {
		t.Errorf("Expected config path '%s', got '%s'", expectedPath, cm.GetConfigPath())
	}
}

func TestConfigEnsureConfigFile(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	// EnsureConfigFile should create the file
	err := cm.EnsureConfigFile()
	if err != nil {
		t.Errorf("EnsureConfigFile() error = %v", err)
	}
	if !cm.ConfigFileExists() {
		t.Error("Config file should exist after EnsureConfigFile()")
	}

	// Load from the created file
	err = cm.Load()
	if err != nil {
		t.Errorf("Load() after EnsureConfigFile() error = %v", err)
	}
}

func TestConfigExportConfig(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	_ = cm.SetMode(ModeScientific)
	cm.SetAutoSave(true)
	cm.SetShowTimestamps(false)

	data := cm.ExportConfig()
	if data == nil {
		t.Fatal("ExportConfig() returned nil")
	}
	if data["mode"] != ModeScientific {
		t.Errorf("Expected mode 'scientific', got '%v'", data["mode"])
	}
	if data["history_limit"] != 1000 {
		t.Errorf("Expected history_limit 1000, got %v", data["history_limit"])
	}
}

func TestConfigImportConfig(t *testing.T) {
	dir := t.TempDir()
	cm := NewConfigManager(dir)

	data := map[string]interface{}{
		"mode":            ModeScientific,
		"angle_unit":      AngleUnitDeg,
		"decimal_places":  float64(5),
		"history_limit":   float64(500),
		"auto_save":       true,
		"show_timestamps": false,
		"theme":           "dark",
	}

	err := cm.ImportConfig(data)
	if err != nil {
		t.Errorf("ImportConfig() error = %v", err)
	}

	cfg := cm.GetConfig()
	if cfg.Mode != ModeScientific {
		t.Errorf("Expected mode 'scientific', got '%s'", cfg.Mode)
	}
	if cfg.DecimalPlaces != 5 {
		t.Errorf("Expected decimal places 5, got %d", cfg.DecimalPlaces)
	}
	if cfg.Theme != "dark" {
		t.Errorf("Expected theme 'dark', got '%s'", cfg.Theme)
	}
}
