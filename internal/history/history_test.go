package history

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ricoarliyanto/calculator/internal/database"
)

// setupTestHistory creates a HistoryManager with a test database
func setupTestHistory(t *testing.T) *HistoryManager {
	t.Helper()
	dir := t.TempDir()
	db := database.NewDatabase(dir)
	if err := db.Connect(); err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	if err := db.RunMigrations(); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}
	return NewHistoryManager(db)
}

// =============================================================================
// AddCalculation Tests
// =============================================================================

func TestAddCalculation(t *testing.T) {
	hm := setupTestHistory(t)

	calc, err := hm.AddCalculation("2 + 2", "4", database.ModeStandard)
	if err != nil {
		t.Errorf("AddCalculation() error = %v", err)
	}
	if calc.ID == 0 {
		t.Error("Expected non-zero ID")
	}
}

func TestAddMultipleCalculations(t *testing.T) {
	hm := setupTestHistory(t)

	inputs := []struct {
		expr, result string
		mode         database.CalculationMode
	}{
		{"2 + 2", "4", database.ModeStandard},
		{"sin(90)", "1", database.ModeScientific},
		{"3 * 3", "9", database.ModeStandard},
		{"cos(0)", "1", database.ModeScientific},
	}

	for _, tc := range inputs {
		_, err := hm.AddCalculation(tc.expr, tc.result, tc.mode)
		if err != nil {
			t.Errorf("AddCalculation(%q) error = %v", tc.expr, err)
		}
	}
}

// =============================================================================
// GetHistory Tests
// =============================================================================

func TestGetHistory(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)
	_, _ = hm.AddCalculation("sin(90)", "1", database.ModeScientific)

	history, err := hm.GetHistory()
	if err != nil {
		t.Errorf("GetHistory() error = %v", err)
	}
	if len(history) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(history))
	}
}

func TestGetHistoryWithLimit(t *testing.T) {
	hm := setupTestHistory(t)

	// Insert exactly 5 calculations
	for i := 0; i < 5; i++ {
		_, _ = hm.AddCalculation("test", "result", database.ModeStandard)
	}

	history, err := hm.GetHistory()
	if err != nil {
		t.Errorf("GetHistory() error = %v", err)
	}
	if len(history) != 5 {
		t.Errorf("Expected 5 history entries, got %d", len(history))
	}
}

func TestGetHistoryEmpty(t *testing.T) {
	hm := setupTestHistory(t)

	history, err := hm.GetHistory()
	if err != nil {
		t.Errorf("GetHistory() error = %v", err)
	}
	if len(history) != 0 {
		t.Errorf("Expected empty history, got %d entries", len(history))
	}
}

// =============================================================================
// GetHistoryByFilter Tests
// =============================================================================

func TestGetHistoryByMode(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)
	_, _ = hm.AddCalculation("sin(90)", "1", database.ModeScientific)
	_, _ = hm.AddCalculation("3 * 3", "9", database.ModeStandard)
	_, _ = hm.AddCalculation("cos(0)", "1", database.ModeScientific)

	history, err := hm.GetHistoryByMode(database.ModeScientific)
	if err != nil {
		t.Errorf("GetHistoryByMode() error = %v", err)
	}
	if len(history) != 2 {
		t.Errorf("Expected 2 scientific history entries, got %d", len(history))
	}

	for _, h := range history {
		if h.Mode != database.ModeScientific {
			t.Errorf("Expected mode 'scientific', got '%s'", h.Mode)
		}
	}
}

func TestGetHistoryByDateRange(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)

	now := time.Now()
	start := now.AddDate(0, 0, -1)
	end := now.AddDate(0, 0, 1)

	history, err := hm.GetHistoryByDateRange(start, end)
	if err != nil {
		t.Errorf("GetHistoryByDateRange() error = %v", err)
	}
	if len(history) != 1 {
		t.Errorf("Expected 1 history entry in date range, got %d", len(history))
	}

	// Test outside range
	history, err = hm.GetHistoryByDateRange(now.AddDate(1, 0, 0), now.AddDate(2, 0, 0))
	if err != nil {
		t.Errorf("GetHistoryByDateRange() error = %v", err)
	}
	if len(history) != 0 {
		t.Errorf("Expected 0 history entries outside date range, got %d", len(history))
	}
}

func TestSearchHistory(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)
	_, _ = hm.AddCalculation("sin(90)", "1", database.ModeScientific)
	_, _ = hm.AddCalculation("3 * 3", "9", database.ModeStandard)

	history, err := hm.SearchHistory("sin")
	if err != nil {
		t.Errorf("SearchHistory() error = %v", err)
	}
	if len(history) != 1 {
		t.Errorf("Expected 1 result matching 'sin', got %d", len(history))
	}
}

func TestSearchHistoryNoMatch(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)

	history, err := hm.SearchHistory("xyz")
	if err != nil {
		t.Errorf("SearchHistory() error = %v", err)
	}
	if len(history) != 0 {
		t.Errorf("Expected 0 results for 'xyz', got %d", len(history))
	}
}

// =============================================================================
// Delete and Clear Tests
// =============================================================================

func TestDeleteCalculation(t *testing.T) {
	hm := setupTestHistory(t)

	calc, _ := hm.AddCalculation("2 + 2", "4", database.ModeStandard)

	err := hm.DeleteCalculation(calc.ID)
	if err != nil {
		t.Errorf("DeleteCalculation() error = %v", err)
	}

	history, _ := hm.GetHistory()
	if len(history) != 0 {
		t.Errorf("Expected 0 history entries after delete, got %d", len(history))
	}
}

func TestDeleteCalculationNotFound(t *testing.T) {
	hm := setupTestHistory(t)

	err := hm.DeleteCalculation(999)
	if err != nil {
		t.Errorf("DeleteCalculation() for non-existent ID should not error, got %v", err)
	}
}

func TestClearHistory(t *testing.T) {
	hm := setupTestHistory(t)

	for i := 0; i < 5; i++ {
		_, _ = hm.AddCalculation("test", "result", database.ModeStandard)
	}

	err := hm.ClearHistory()
	if err != nil {
		t.Errorf("ClearHistory() error = %v", err)
	}

	history, _ := hm.GetHistory()
	if len(history) != 0 {
		t.Errorf("Expected 0 history entries after clear, got %d", len(history))
	}
}

func TestClearHistoryEmpty(t *testing.T) {
	hm := setupTestHistory(t)

	err := hm.ClearHistory()
	if err != nil {
		t.Errorf("ClearHistory() on empty history error = %v", err)
	}
}

// =============================================================================
// Export Tests
// =============================================================================

func TestExportHistoryJSON(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)
	_, _ = hm.AddCalculation("sin(90)", "1", database.ModeScientific)

	data, err := hm.ExportHistory(database.ExportJSON)
	if err != nil {
		t.Errorf("ExportHistory(JSON) error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("Expected non-empty JSON export")
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("Invalid JSON format: %v", err)
	}
	if len(parsed) != 2 {
		t.Errorf("Expected 2 entries in JSON export, got %d", len(parsed))
	}
}

func TestExportHistoryCSV(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)

	data, err := hm.ExportHistory(database.ExportCSV)
	if err != nil {
		t.Errorf("ExportHistory(CSV) error = %v", err)
	}
	if len(data) == 0 {
		t.Fatal("Expected non-empty CSV export")
	}
}

func TestExportHistoryEmpty(t *testing.T) {
	hm := setupTestHistory(t)

	data, err := hm.ExportHistory(database.ExportJSON)
	if err != nil {
		t.Errorf("ExportHistory() error = %v", err)
	}
	if len(data) == 0 {
		t.Error("Expected valid JSON export (even if empty)")
	}
}

func TestExportHistoryToFileJSON(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)

	exportDir := t.TempDir()
	filePath, err := hm.ExportHistoryToFile(database.ExportJSON, exportDir)
	if err != nil {
		t.Errorf("ExportHistoryToFile() error = %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read exported file: %v", err)
	}

	var parsed []map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("Invalid JSON in exported file: %v", err)
	}
}

func TestExportHistoryToFileCSV(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)

	exportDir := t.TempDir()
	filePath, err := hm.ExportHistoryToFile(database.ExportCSV, exportDir)
	if err != nil {
		t.Errorf("ExportHistoryToFile() error = %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Errorf("Failed to read exported file: %v", err)
	}

	records, err := csv.NewReader(bytes.NewReader(data)).ReadAll()
	if err != nil {
		t.Errorf("Invalid CSV format: %v", err)
	}
	// Header + 1 data row
	if len(records) < 2 {
		t.Errorf("Expected at least 2 CSV rows (header + data), got %d", len(records))
	}
}

// =============================================================================
// Import Tests
// =============================================================================

func TestImportHistoryFromFileJSON(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)

	// Export first
	exportDir := t.TempDir()
	_, _ = hm.ExportHistoryToFile(database.ExportJSON, exportDir)

	// Find the exported file
	files, _ := filepath.Glob(filepath.Join(exportDir, "*.json"))
	if len(files) == 0 {
		t.Fatal("No exported file found")
	}

	// Create another manager and import
	dir2 := t.TempDir()
	db2 := database.NewDatabase(dir2)
	_ = db2.Connect()
	_ = db2.RunMigrations()
	hm2 := NewHistoryManager(db2)

	err := hm2.ImportHistoryFromFile(files[0])
	if err != nil {
		t.Errorf("ImportHistoryFromFile() error = %v", err)
	}

	history, _ := hm2.GetHistory()
	if len(history) != 1 {
		t.Errorf("Expected 1 calculation after import, got %d", len(history))
	}
}

func TestImportHistoryFromNonExistentFile(t *testing.T) {
	hm := setupTestHistory(t)

	err := hm.ImportHistoryFromFile("/nonexistent/file.json")
	if err == nil {
		t.Error("Expected error for non-existent import file")
	}
}

// =============================================================================
// Stats Tests
// =============================================================================

func TestGetStats(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)
	_, _ = hm.AddCalculation("3 + 3", "6", database.ModeStandard)
	_, _ = hm.AddCalculation("sin(90)", "1", database.ModeScientific)
	_, _ = hm.AddCalculation("cos(0)", "1", database.ModeScientific)

	stats, err := hm.GetStats()
	if err != nil {
		t.Errorf("GetStats() error = %v", err)
	}
	if stats["total_calculations"].(int) != 4 {
		t.Errorf("Expected 4 total calculations, got %v", stats["total_calculations"])
	}
	byMode := stats["by_mode"].(map[string]int)
	if byMode["standard"] != 2 {
		t.Errorf("Expected 2 standard count, got %d", byMode["standard"])
	}
	if byMode["scientific"] != 2 {
		t.Errorf("Expected 2 scientific count, got %d", byMode["scientific"])
	}
}

func TestGetStatsEmpty(t *testing.T) {
	hm := setupTestHistory(t)

	stats, err := hm.GetStats()
	if err != nil {
		t.Errorf("GetStats() error = %v", err)
	}
	if stats["total_calculations"].(int) != 0 {
		t.Errorf("Expected 0 total calculations, got %v", stats["total_calculations"])
	}
}

// =============================================================================
// History Limit Tests
// =============================================================================

func TestGetSetHistoryLimit(t *testing.T) {
	hm := setupTestHistory(t)

	// Default limit
	limit, err := hm.GetHistoryLimit()
	if err != nil {
		t.Errorf("GetHistoryLimit() error = %v", err)
	}
	if limit != 1000 {
		t.Errorf("Expected default limit 1000, got %d", limit)
	}

	// Set new limit
	err = hm.SetHistoryLimit(100)
	if err != nil {
		t.Errorf("SetHistoryLimit() error = %v", err)
	}

	limit, _ = hm.GetHistoryLimit()
	if limit != 100 {
		t.Errorf("Expected new limit 100, got %d", limit)
	}
}

// =============================================================================
// GetHistoryByFilter Tests
// =============================================================================

func TestGetHistoryByFilterModeAndSearch(t *testing.T) {
	hm := setupTestHistory(t)

	_, _ = hm.AddCalculation("2 + 2", "4", database.ModeStandard)
	_, _ = hm.AddCalculation("sin(90)", "1", database.ModeScientific)
	_, _ = hm.AddCalculation("cos(90)", "0", database.ModeScientific)

	mode := database.ModeScientific
	search := "sin"
	filter := database.CalculationFilter{
		Mode:   &mode,
		Search: &search,
	}
	history, err := hm.GetHistoryByFilter(filter)
	if err != nil {
		t.Errorf("GetHistoryByFilter() error = %v", err)
	}
	if len(history) != 1 {
		t.Errorf("Expected 1 result, got %d", len(history))
	}
}
