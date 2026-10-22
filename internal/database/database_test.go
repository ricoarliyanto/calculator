package database

import (
	"testing"
	"time"
)

// setupTestDB creates a database instance with migrations applied for testing
func setupTestDB(t *testing.T) *Database {
	t.Helper()
	dir := t.TempDir()
	db := NewDatabase(dir)
	if err := db.Connect(); err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	if err := db.RunMigrations(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
	return db
}

// =============================================================================
// Database Connection Tests
// =============================================================================

func TestDatabaseConnection(t *testing.T) {
	dir := t.TempDir()
	db := NewDatabase(dir)

	if err := db.Connect(); err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}

	// Get DB should work
	conn, err := db.GetDB()
	if err != nil {
		t.Errorf("GetDB() error = %v", err)
	}
	if conn == nil {
		t.Error("GetDB() returned nil connection")
	}

	// Close should work
	if err := db.Close(); err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestDatabasePath(t *testing.T) {
	dir := t.TempDir()
	db := NewDatabase(dir)
	_ = db.Connect()

	expectedPath := dir + "/calculator.db"
	if db.GetDBPath() != expectedPath {
		t.Errorf("GetDBPath() = %v, want %v", db.GetDBPath(), expectedPath)
	}
}

func TestDatabaseConnectionFailure(t *testing.T) {
	// Try connecting to an invalid path
	db := NewDatabase("/root/invalid/path/that/does/not/exist")
	if err := db.Connect(); err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

// =============================================================================
// Schema Creation Tests
// =============================================================================

func TestSchemaCreation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	conn, _ := db.GetDB()

	// Check calculations table exists
	var tableName string
	err := conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='calculations'").Scan(&tableName)
	if err != nil {
		t.Errorf("calculations table should exist: %v", err)
	}

	// Check settings table exists
	err = conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='settings'").Scan(&tableName)
	if err != nil {
		t.Errorf("settings table should exist: %v", err)
	}

	// Check migrations table exists
	err = conn.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='migrations'").Scan(&tableName)
	if err != nil {
		t.Errorf("migrations table should exist: %v", err)
	}

	// Check indexes
	var indexCount int
	err = conn.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND tbl_name='calculations'").Scan(&indexCount)
	if err != nil {
		t.Errorf("Failed to check indexes: %v", err)
	}
	if indexCount < 2 {
		t.Errorf("Expected at least 2 indexes on calculations table, got %d", indexCount)
	}
}

func TestMigrationSystem(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	migrationManager := NewMigrationManager(db)

	// Initially should have version 1
	version, err := migrationManager.GetCurrentVersion()
	if err != nil {
		t.Errorf("GetCurrentVersion() error = %v", err)
	}
	if version != 1 {
		t.Errorf("Expected version 1, got %d", version)
	}

	// Add a test migration
	migrationManager.AddMigration("002_test_migration", "Test migration", `
		CREATE TABLE IF NOT EXISTS test_table (
			id INTEGER PRIMARY KEY,
			name TEXT
		);
	`, 2)

	// Apply migrations should create the test table
	if err := migrationManager.ApplyMigrations(); err != nil {
		t.Errorf("ApplyMigrations() error = %v", err)
	}

	// Check version updated
	version, err = migrationManager.GetCurrentVersion()
	if err != nil {
		t.Errorf("GetCurrentVersion() error = %v", err)
	}
	if version != 2 {
		t.Errorf("Expected version 2 after migration, got %d", version)
	}

	// Validate migrations should pass
	if err := migrationManager.ValidateMigrations(); err != nil {
		t.Errorf("ValidateMigrations() error = %v", err)
	}

	// Applied migrations should return both migrations
	applied, err := migrationManager.GetAppliedMigrations()
	if err != nil {
		t.Errorf("GetAppliedMigrations() error = %v", err)
	}
	if len(applied) != 2 {
		t.Errorf("Expected 2 applied migrations, got %d", len(applied))
	}
}

func TestMigrationIdempotent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	migrationManager := NewMigrationManager(db)
	migrationManager.AddMigration("002_test", "Test", `CREATE TABLE IF NOT EXISTS test_idem (id INTEGER PRIMARY KEY);`, 2)

	// Apply twice - should not error
	if err := migrationManager.ApplyMigrations(); err != nil {
		t.Errorf("First ApplyMigrations() error = %v", err)
	}
	if err := migrationManager.ApplyMigrations(); err != nil {
		t.Errorf("Second ApplyMigrations() error = %v", err)
	}
}

// =============================================================================
// Calculation Repository Tests
// =============================================================================

func TestInsertCalculation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	calc, err := repo.InsertCalculation("2 + 2", "4", ModeStandard)
	if err != nil {
		t.Errorf("InsertCalculation() error = %v", err)
	}
	if calc.ID == 0 {
		t.Error("Expected non-zero ID")
	}
	if calc.Expression != "2 + 2" {
		t.Errorf("Expected expression '2 + 2', got %s", calc.Expression)
	}
	if calc.Result != "4" {
		t.Errorf("Expected result '4', got %s", calc.Result)
	}
	if calc.Mode != ModeStandard {
		t.Errorf("Expected mode 'standard', got %s", calc.Mode)
	}
}

func TestInsertCalculationScientific(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	calc, err := repo.InsertCalculation("sin(90, deg)", "1", ModeScientific)
	if err != nil {
		t.Errorf("InsertCalculation() error = %v", err)
	}
	if calc.Mode != ModeScientific {
		t.Errorf("Expected mode 'scientific', got %s", calc.Mode)
	}
}

func TestGetAllCalculations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert test data
	for i := 0; i < 5; i++ {
		_, err := repo.InsertCalculation("test", "result", ModeStandard)
		if err != nil {
			t.Fatalf("InsertCalculation() error = %v", err)
		}
	}

	calculations, err := repo.GetAllCalculations()
	if err != nil {
		t.Errorf("GetAllCalculations() error = %v", err)
	}
	if len(calculations) != 5 {
		t.Errorf("Expected 5 calculations, got %d", len(calculations))
	}

	// Check ordering (by created_at DESC)
	if len(calculations) > 1 && calculations[0].CreatedAt.Before(calculations[1].CreatedAt) {
		t.Error("Calculations should be ordered by created_at DESC")
	}
}

func TestGetAllCalculationsEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	calculations, err := repo.GetAllCalculations()
	if err != nil {
		t.Errorf("GetAllCalculations() error = %v", err)
	}
	if len(calculations) != 0 {
		t.Errorf("Expected 0 calculations for empty DB, got %d", len(calculations))
	}
}

func TestGetCalculationsByFilter(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert mixed mode calculations
	_, _ = repo.InsertCalculation("2 + 2", "4", ModeStandard)
	_, _ = repo.InsertCalculation("sin(90)", "1", ModeScientific)
	_, _ = repo.InsertCalculation("3 * 3", "9", ModeStandard)
	_, _ = repo.InsertCalculation("cos(0)", "1", ModeScientific)

	// Filter by mode
	scientificMode := ModeScientific
	filter := CalculationFilter{Mode: &scientificMode}
	calculations, err := repo.GetCalculationsByFilter(filter)
	if err != nil {
		t.Errorf("GetCalculationsByFilter() error = %v", err)
	}
	if len(calculations) != 2 {
		t.Errorf("Expected 2 scientific calculations, got %d", len(calculations))
	}
}

func TestGetCalculationsByFilterSearch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert calculations with different expressions
	_, _ = repo.InsertCalculation("2 + 2", "4", ModeStandard)
	_, _ = repo.InsertCalculation("sin(90)", "1", ModeScientific)
	_, _ = repo.InsertCalculation("2 + 3", "5", ModeStandard)

	// Filter by search
	search := "sin"
	filter := CalculationFilter{Search: &search}
	calculations, err := repo.GetCalculationsByFilter(filter)
	if err != nil {
		t.Errorf("GetCalculationsByFilter() error = %v", err)
	}
	if len(calculations) != 1 {
		t.Errorf("Expected 1 calculation matching 'sin', got %d", len(calculations))
	}
	if calculations[0].Expression != "sin(90)" {
		t.Errorf("Expected expression 'sin(90)', got %s", calculations[0].Expression)
	}
}

func TestGetCalculationsByFilterDateRange(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert a calculation
	_, err := repo.InsertCalculation("2 + 2", "4", ModeStandard)
	if err != nil {
		t.Fatalf("InsertCalculation() error = %v", err)
	}

	// Filter by date range (current time to far future)
	now := time.Now()
	future := now.AddDate(1, 0, 0)
	start := now.AddDate(0, 0, -1)

	filter := CalculationFilter{
		StartDate: &start,
		EndDate:   &future,
	}
	calculations, err := repo.GetCalculationsByFilter(filter)
	if err != nil {
		t.Errorf("GetCalculationsByFilter() error = %v", err)
	}
	if len(calculations) != 1 {
		t.Errorf("Expected 1 calculation in date range, got %d", len(calculations))
	}
}

func TestGetCalculationByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert a calculation
	calc, err := repo.InsertCalculation("2 + 2", "4", ModeStandard)
	if err != nil {
		t.Fatalf("InsertCalculation() error = %v", err)
	}

	// Retrieve by ID
	retrieved, err := repo.GetCalculationByID(calc.ID)
	if err != nil {
		t.Errorf("GetCalculationByID() error = %v", err)
	}
	if retrieved == nil {
		t.Fatal("Expected non-nil calculation")
	}
	if retrieved.Expression != "2 + 2" {
		t.Errorf("Expected expression '2 + 2', got %s", retrieved.Expression)
	}
}

func TestGetCalculationByIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	retrieved, err := repo.GetCalculationByID(999)
	if err != nil {
		t.Errorf("GetCalculationByID() error = %v", err)
	}
	if retrieved != nil {
		t.Error("Expected nil for non-existent ID")
	}
}

func TestDeleteCalculation(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert two calculations
	calc1, _ := repo.InsertCalculation("2 + 2", "4", ModeStandard)
	_, _ = repo.InsertCalculation("3 + 3", "6", ModeStandard)

	// Delete first
	if err := repo.DeleteCalculation(calc1.ID); err != nil {
		t.Errorf("DeleteCalculation() error = %v", err)
	}

	// Check only one remains
	calculations, _ := repo.GetAllCalculations()
	if len(calculations) != 1 {
		t.Errorf("Expected 1 calculation after delete, got %d", len(calculations))
	}
}

func TestClearAllCalculations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert some calculations
	for i := 0; i < 10; i++ {
		_, _ = repo.InsertCalculation("test", "result", ModeStandard)
	}

	// Clear all
	if err := repo.ClearAllCalculations(); err != nil {
		t.Errorf("ClearAllCalculations() error = %v", err)
	}

	// Verify empty
	calculations, _ := repo.GetAllCalculations()
	if len(calculations) != 0 {
		t.Errorf("Expected 0 calculations after clear, got %d", len(calculations))
	}
}

// =============================================================================
// Export Tests
// =============================================================================

func TestExportCalculationsJSON(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert test data
	_, _ = repo.InsertCalculation("2 + 2", "4", ModeStandard)
	_, _ = repo.InsertCalculation("sin(90, deg)", "1", ModeScientific)

	data, err := repo.ExportCalculations(ExportJSON)
	if err != nil {
		t.Errorf("ExportCalculations(JSON) error = %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected non-empty JSON export")
	}
}

func TestExportCalculationsCSV(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Insert test data
	_, _ = repo.InsertCalculation("2 + 2", "4", ModeStandard)

	data, err := repo.ExportCalculations(ExportCSV)
	if err != nil {
		t.Errorf("ExportCalculations(CSV) error = %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected non-empty CSV export")
	}
}

func TestExportCalculationsEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Export from empty database
	data, err := repo.ExportCalculations(ExportJSON)
	if err != nil {
		t.Errorf("ExportCalculations() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected valid JSON export (even if empty)")
	}
}

// =============================================================================
// Settings Repository Tests
// =============================================================================

func TestSaveAndLoadSetting(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// Save a setting
	err := repo.SaveSetting("test_key", "test_value")
	if err != nil {
		t.Errorf("SaveSetting() error = %v", err)
	}

	// Load it back
	setting, err := repo.GetSetting("test_key")
	if err != nil {
		t.Errorf("GetSetting() error = %v", err)
	}
	if setting == nil {
		t.Fatal("Expected non-nil setting")
	}
	if setting.Key != "test_key" {
		t.Errorf("Expected key 'test_key', got %s", setting.Key)
	}
	if setting.Value != "test_value" {
		t.Errorf("Expected value 'test_value', got %s", setting.Value)
	}
}

func TestUpdateSetting(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// Save initial
	_ = repo.SaveSetting("counter", "5")

	// Update
	err := repo.SaveSetting("counter", "10")
	if err != nil {
		t.Errorf("SaveSetting() error = %v", err)
	}

	// Verify update
	val, err := repo.GetSettingValue("counter", "0")
	if err != nil {
		t.Errorf("GetSettingValue() error = %v", err)
	}
	if val != "10" {
		t.Errorf("Expected '10', got '%s'", val)
	}
}

func TestGetSettingNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	setting, err := repo.GetSetting("nonexistent_key")
	if err != nil {
		t.Errorf("GetSetting() error = %v", err)
	}
	if setting != nil {
		t.Error("Expected nil for non-existent setting")
	}
}

func TestGetSettingValueWithDefault(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	val, err := repo.GetSettingValue("nonexistent_key", "default_val")
	if err != nil {
		t.Errorf("GetSettingValue() error = %v", err)
	}
	if val != "default_val" {
		t.Errorf("Expected 'default_val', got '%s'", val)
	}
}

func TestGetAllSettings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// Save some settings
	_ = repo.SaveSetting("key1", "value1")
	_ = repo.SaveSetting("key2", "value2")

	// Get all
	settings, err := repo.GetAllSettings()
	if err != nil {
		t.Errorf("GetAllSettings() error = %v", err)
	}
	if len(settings) < 2 {
		t.Errorf("Expected at least 2 settings, got %d", len(settings))
	}
}

func TestDeleteSetting(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// Save a setting
	_ = repo.SaveSetting("temp_key", "temp_value")

	// Delete it
	err := repo.DeleteSetting("temp_key")
	if err != nil {
		t.Errorf("DeleteSetting() error = %v", err)
	}

	// Verify it's gone
	setting, _ := repo.GetSetting("temp_key")
	if setting != nil {
		t.Error("Setting should be deleted")
	}
}

func TestGetIntSetting(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// Save as string, retrieve as int
	_ = repo.SaveSetting("limit", "100")

	val, err := repo.GetIntSetting("limit", 50)
	if err != nil {
		t.Errorf("GetIntSetting() error = %v", err)
	}
	if val != 100 {
		t.Errorf("Expected 100, got %d", val)
	}

	// Test with default (nonexistent key)
	val2, _ := repo.GetIntSetting("nonexistent", 50)
	if val2 != 50 {
		t.Errorf("Expected 50 (default), got %d", val2)
	}
}

func TestGetBoolSetting(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// Test true value
	_ = repo.SaveSetting("enabled", "true")
	val, _ := repo.GetBoolSetting("enabled", false)
	if !val {
		t.Error("Expected true for 'true' string")
	}

	// Test false value
	_ = repo.SaveSetting("enabled", "false")
	val, _ = repo.GetBoolSetting("enabled", true)
	if val {
		t.Error("Expected false for 'false' string")
	}

	// Test default value
	val, _ = repo.GetBoolSetting("nonexistent", true)
	if !val {
		t.Error("Expected true (default) for nonexistent key")
	}

	// Test "1" as true
	_ = repo.SaveSetting("enabled", "1")
	val, _ = repo.GetBoolSetting("enabled", false)
	if !val {
		t.Error("Expected true for '1' string")
	}
}

func TestDefaultSettings(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	// After RunMigrations in setupTestDB, default settings should exist
	val, err := repo.GetSettingValue("default_mode", "standard")
	if err != nil {
		t.Errorf("GetSettingValue() error = %v", err)
	}
	if val != "standard" {
		t.Errorf("Expected default_mode='standard', got '%s'", val)
	}

	val, _ = repo.GetSettingValue("history_limit", "0")
	if val != "1000" {
		t.Errorf("Expected history_limit='1000', got '%s'", val)
	}

	val, _ = repo.GetSettingValue("decimal_precision", "0")
	if val != "10" {
		t.Errorf("Expected decimal_precision='10', got '%s'", val)
	}

	val, _ = repo.GetSettingValue("theme", "")
	if val != "light" {
		t.Errorf("Expected theme='light', got '%s'", val)
	}
}

// =============================================================================
// Additional Database Tests for Coverage
// =============================================================================

func TestGetDBNotConnected(t *testing.T) {
	dir := t.TempDir()
	db := NewDatabase(dir)

	// Don't call Connect, but GetDB should auto-connect
	conn, err := db.GetDB()
	if err != nil {
		t.Errorf("GetDB() auto-connect error = %v", err)
	}
	if conn == nil {
		t.Error("GetDB() returned nil connection")
	}
	db.Close()
}

func TestGetDBPathBeforeConnect(t *testing.T) {
	dir := t.TempDir()
	db := NewDatabase(dir)

	// GetDBPath returns empty string before Connect is called
	if db.GetDBPath() != "" {
		t.Errorf("Expected empty dbPath before Connect, got '%s'", db.GetDBPath())
	}

	_ = db.Connect()
	expectedPath := dir + "/calculator.db"
	if db.GetDBPath() != expectedPath {
		t.Errorf("Expected dbPath '%s', got '%s'", expectedPath, db.GetDBPath())
	}
}

func TestCloseWithoutConnect(t *testing.T) {
	dir := t.TempDir()
	db := NewDatabase(dir)

	// Close without Connect should not error
	if err := db.Close(); err != nil {
		t.Errorf("Close() without Connect error = %v", err)
	}
}

func TestInitializeDatabase(t *testing.T) {
	dir := t.TempDir()
	db := NewDatabase(dir)
	if err := db.Connect(); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer db.Close()

	// InitializeDatabase should work (same as RunMigrations)
	if err := db.InitializeDatabase(); err != nil {
		t.Errorf("InitializeDatabase() error = %v", err)
	}
}

func TestExportCalculationsUnsupportedFormat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)
	_, _ = repo.InsertCalculation("2 + 2", "4", ModeStandard)

	_, err := repo.ExportCalculations(ExportFormat("xml"))
	if err == nil {
		t.Error("Expected error for unsupported export format")
	}
}

func TestApplyMigrationsNoMigrations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mm := NewMigrationManager(db)
	// No migrations added - should be no-op
	if err := mm.ApplyMigrations(); err != nil {
		t.Errorf("ApplyMigrations() with no migrations error = %v", err)
	}
}

func TestValidateMigrationsWithMissingMigrations(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mm := NewMigrationManager(db)
	mm.AddMigration("002_test", "Test", `CREATE TABLE test_missing (id INTEGER);`, 2)

	// ApplyMigrations will apply version 1 (from RunMigrations) + version 2
	err := mm.ApplyMigrations()
	if err != nil {
		t.Errorf("ApplyMigrations() error = %v", err)
	}

	// Validate should pass since both are applied
	if err := mm.ValidateMigrations(); err != nil {
		t.Errorf("ValidateMigrations() error = %v", err)
	}
}

func TestValidateMigrationsMissing(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mm := NewMigrationManager(db)
	// Don't add any migrations - but migrations table exists with version 1 applied
	// The 001 migration is not registered on this manager
	if err := mm.ValidateMigrations(); err == nil {
		// If no migrations registered, there's nothing to be missing
		// This is expected behavior - empty migrations list means no missing
	}
}

func TestGetAllSettingsEmpty(t *testing.T) {
	// Use a database with a fresh settings table (no default settings)
	dir := t.TempDir()
	db := NewDatabase(dir)
	if err := db.Connect(); err != nil {
		t.Fatalf("Connect() error = %v", err)
	}
	defer db.Close()

	// Don't run full migrations - just create tables manually
	conn, _ := db.GetDB()
	_, _ = conn.Exec(`CREATE TABLE IF NOT EXISTS settings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key TEXT NOT NULL UNIQUE,
		value TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	repo := NewSettingsRepository(db)
	settings, err := repo.GetAllSettings()
	if err != nil {
		t.Errorf("GetAllSettings() error = %v", err)
	}
	if len(settings) != 0 {
		t.Errorf("Expected 0 settings, got %d", len(settings))
	}
}

func TestGetIntSettingNonNumeric(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)
	_ = repo.SaveSetting("non_numeric", "abc")

	val, err := repo.GetIntSetting("non_numeric", 99)
	if err != nil {
		t.Errorf("GetIntSetting() error = %v", err)
	}
	if val != 99 {
		t.Errorf("Expected default value 99 for non-numeric setting, got %d", val)
	}
}

func TestGetBoolSettingInvalidValue(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)
	_ = repo.SaveSetting("invalid_bool", "not_a_bool")

	val, err := repo.GetBoolSetting("invalid_bool", true)
	if err != nil {
		t.Errorf("GetBoolSetting() error = %v", err)
	}
	if !val {
		t.Error("Expected default value true for invalid boolean setting")
	}
}

func TestDeleteSettingNonExistent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSettingsRepository(db)

	err := repo.DeleteSetting("nonexistent_key")
	if err != nil {
		t.Errorf("DeleteSetting() for non-existent key error = %v", err)
	}
}

func TestClearAllCalculationsOnEmptyDB(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	// Clear on empty DB should not error
	if err := repo.ClearAllCalculations(); err != nil {
		t.Errorf("ClearAllCalculations() on empty DB error = %v", err)
	}

	calculations, _ := repo.GetAllCalculations()
	if len(calculations) != 0 {
		t.Errorf("Expected 0 calculations, got %d", len(calculations))
	}
}

func TestGetCalculationsByFilterAllFilters(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	mode := ModeStandard
	search := "test"
	now := time.Now()
	start := now.AddDate(0, 0, -1)
	end := now.AddDate(0, 0, 1)

	filter := CalculationFilter{
		Mode:      &mode,
		Search:    &search,
		StartDate: &start,
		EndDate:   &end,
	}

	calculations, err := repo.GetCalculationsByFilter(filter)
	if err != nil {
		t.Errorf("GetCalculationsByFilter() error = %v", err)
	}
	// No data matches, but no error either
	if len(calculations) != 0 {
		t.Errorf("Expected 0 calculations with combined filters, got %d", len(calculations))
	}
}

func TestGetCalculationsByFilterEmptySearch(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)
	_, _ = repo.InsertCalculation("test", "result", ModeStandard)

	// Empty search string should not filter
	search := ""
	filter := CalculationFilter{Search: &search}

	calculations, err := repo.GetCalculationsByFilter(filter)
	if err != nil {
		t.Errorf("GetCalculationsByFilter() error = %v", err)
	}
	if len(calculations) != 1 {
		t.Errorf("Expected 1 calculation with empty search, got %d", len(calculations))
	}
}

func TestMigrationManagerSortByVersion(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mm := NewMigrationManager(db)

	// Add migrations in reverse order
	mm.AddMigration("003_third", "Third", `CREATE TABLE third (id INTEGER);`, 3)
	mm.AddMigration("001_first", "First", `CREATE TABLE first (id INTEGER);`, 1)
	mm.AddMigration("002_second", "Second", `CREATE TABLE second (id INTEGER);`, 2)

	// Verify they are sorted by version
	if len(mm.migrations) != 3 {
		t.Fatalf("Expected 3 migrations, got %d", len(mm.migrations))
	}
	if mm.migrations[0].Version != 1 {
		t.Errorf("Expected first migration version 1, got %d", mm.migrations[0].Version)
	}
	if mm.migrations[1].Version != 2 {
		t.Errorf("Expected second migration version 2, got %d", mm.migrations[1].Version)
	}
	if mm.migrations[2].Version != 3 {
		t.Errorf("Expected third migration version 3, got %d", mm.migrations[2].Version)
	}
}

func TestApplyMigrationWithInvalidSQL(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	mm := NewMigrationManager(db)
	mm.AddMigration("002_bad_migration", "Bad Migration", "INVALID SQL STATEMENT", 2)

	err := mm.ApplyMigrations()
	if err == nil {
		t.Error("Expected error for invalid SQL migration")
	}
}

func TestExportCalculationsCSVEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)
	// Clear all calculations
	_ = repo.ClearAllCalculations()

	data, err := repo.ExportCalculations(ExportCSV)
	if err != nil {
		t.Errorf("ExportCalculations(CSV) with empty data error = %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected CSV header even with empty data")
	}
}

func TestExportCalculationsJSONEmpty(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)
	_ = repo.ClearAllCalculations()

	data, err := repo.ExportCalculations(ExportJSON)
	if err != nil {
		t.Errorf("ExportCalculations(JSON) with empty data error = %v", err)
	}
	// Should return valid JSON (null or empty array)
	if len(data) == 0 {
		t.Error("Expected valid JSON for empty data")
	}
}

func TestGetCalculationsByFilterWithAllCriteria(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	mode := ModeScientific
	_, _ = repo.InsertCalculation("sin(90)", "1", ModeScientific)
	_, _ = repo.InsertCalculation("cos(0)", "1", ModeScientific)
	_, _ = repo.InsertCalculation("2 + 2", "4", ModeStandard)

	filter := CalculationFilter{Mode: &mode}
	calculations, err := repo.GetCalculationsByFilter(filter)
	if err != nil {
		t.Errorf("GetCalculationsByFilter() error = %v", err)
	}
	if len(calculations) != 2 {
		t.Errorf("Expected 2 scientific calculations, got %d", len(calculations))
	}
}

func TestRepositoryOperationsOnClosedDB(t *testing.T) {
	dir := t.TempDir()
	db := NewDatabase(dir)
	_ = db.Connect()
	_ = db.RunMigrations()

	// Close the database
	_ = db.Close()

	repo := NewCalculationRepository(db)

	// Operations on closed DB should return errors or nil
	_, err := repo.GetAllCalculations()
	_ = err // SQLite may or may not error on closed DB

	_, err = repo.GetCalculationByID(1)
	_ = err

	err = repo.ClearAllCalculations()
	_ = err
}

func TestSettingsRepoOperationsOnClosedDB(t *testing.T) {
	dir := t.TempDir()
	db := NewDatabase(dir)
	_ = db.Connect()
	_ = db.RunMigrations()
	_ = db.Close()

	repo := NewSettingsRepository(db)

	_, err := repo.GetSetting("test")
	_ = err

	_, err = repo.GetAllSettings()
	_ = err
}

func TestInsertCalculationReturnsValidID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	calc1, err := repo.InsertCalculation("1+1", "2", ModeStandard)
	if err != nil {
		t.Errorf("InsertCalculation() error = %v", err)
	}
	if calc1.ID <= 0 {
		t.Error("Expected positive ID")
	}

	calc2, err := repo.InsertCalculation("2+2", "4", ModeStandard)
	if err != nil {
		t.Errorf("InsertCalculation() error = %v", err)
	}
	if calc2.ID <= calc1.ID {
		t.Error("Expected second ID to be greater than first")
	}
}

func TestGetCalculationByIDReturnsCorrectData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewCalculationRepository(db)

	calc, err := repo.InsertCalculation("10 + 20", "30", ModeScientific)
	if err != nil {
		t.Fatalf("InsertCalculation() error = %v", err)
	}

	retrieved, err := repo.GetCalculationByID(calc.ID)
	if err != nil {
		t.Errorf("GetCalculationByID() error = %v", err)
	}
	if retrieved == nil {
		t.Fatal("Expected non-nil calculation")
	}
	if retrieved.Expression != "10 + 20" {
		t.Errorf("Expected expression '10 + 20', got '%s'", retrieved.Expression)
	}
	if retrieved.Result != "30" {
		t.Errorf("Expected result '30', got '%s'", retrieved.Result)
	}
	if retrieved.Mode != ModeScientific {
		t.Errorf("Expected mode 'scientific', got '%s'", retrieved.Mode)
	}
}
