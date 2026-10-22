package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ricoarliyanto/calculator/internal/database"
)

// HistoryManager manages calculation history
type HistoryManager struct {
	calcRepo  *database.CalculationRepository
	settingsRepo *database.SettingsRepository
}

// NewHistoryManager creates a new history manager
func NewHistoryManager(db *database.Database) *HistoryManager {
	return &HistoryManager{
		calcRepo:  database.NewCalculationRepository(db),
		settingsRepo: database.NewSettingsRepository(db),
	}
}

// AddCalculation adds a calculation to history
func (h *HistoryManager) AddCalculation(expression, result string, mode database.CalculationMode) (*database.Calculation, error) {
	return h.calcRepo.InsertCalculation(expression, result, mode)
}

// GetHistory retrieves all calculation history
func (h *HistoryManager) GetHistory() ([]database.Calculation, error) {
	return h.calcRepo.GetAllCalculations()
}

// GetHistoryByFilter retrieves filtered history
func (h *HistoryManager) GetHistoryByFilter(filter database.CalculationFilter) ([]database.Calculation, error) {
	return h.calcRepo.GetCalculationsByFilter(filter)
}

// GetHistoryByMode retrieves history filtered by mode
func (h *HistoryManager) GetHistoryByMode(mode database.CalculationMode) ([]database.Calculation, error) {
	filter := database.CalculationFilter{
		Mode: &mode,
	}
	return h.calcRepo.GetCalculationsByFilter(filter)
}

// GetHistoryByDateRange retrieves history within a date range
func (h *HistoryManager) GetHistoryByDateRange(start, end time.Time) ([]database.Calculation, error) {
	filter := database.CalculationFilter{
		StartDate: &start,
		EndDate:   &end,
	}
	return h.calcRepo.GetCalculationsByFilter(filter)
}

// SearchHistory searches history by expression or result
func (h *HistoryManager) SearchHistory(query string) ([]database.Calculation, error) {
	filter := database.CalculationFilter{
		Search: &query,
	}
	return h.calcRepo.GetCalculationsByFilter(filter)
}

// DeleteCalculation deletes a calculation from history
func (h *HistoryManager) DeleteCalculation(id int64) error {
	return h.calcRepo.DeleteCalculation(id)
}

// ClearHistory clears all calculation history
func (h *HistoryManager) ClearHistory() error {
	return h.calcRepo.ClearAllCalculations()
}

// GetHistoryLimit returns the maximum number of history entries to keep
func (h *HistoryManager) GetHistoryLimit() (int, error) {
	return h.settingsRepo.GetIntSetting("history_limit", 1000)
}

// SetHistoryLimit sets the maximum number of history entries to keep
func (h *HistoryManager) SetHistoryLimit(limit int) error {
	return h.settingsRepo.SaveSetting("history_limit", fmt.Sprintf("%d", limit))
}

// ExportHistory exports history in the specified format
func (h *HistoryManager) ExportHistory(format database.ExportFormat) ([]byte, error) {
	return h.calcRepo.ExportCalculations(format)
}

// ExportHistoryToFile exports history to a file
func (h *HistoryManager) ExportHistoryToFile(format database.ExportFormat, dir string) (string, error) {
	data, err := h.ExportHistory(format)
	if err != nil {
		return "", fmt.Errorf("failed to export history: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create export directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("calculator_history_%s.%s", timestamp, strings.ToLower(string(format)))
	filepath := filepath.Join(dir, filename)

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write export file: %w", err)
	}

	return filepath, nil
}

// ImportHistoryFromFile imports history from a JSON file
func (h *HistoryManager) ImportHistoryFromFile(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	var calculations []map[string]interface{}
	if err := json.Unmarshal(data, &calculations); err != nil {
		return fmt.Errorf("failed to parse import file: %w", err)
	}

	// Sort by created_at descending to maintain order
	sort.Slice(calculations, func(i, j int) bool {
		iTime := calculations[i]["created_at"].(string)
		jTime := calculations[j]["created_at"].(string)
		return iTime > jTime
	})

	for _, calc := range calculations {
		expression, ok1 := calc["expression"].(string)
		result, ok2 := calc["result"].(string)
		modeStr, ok3 := calc["mode"].(string)
		mode := database.ModeStandard

		if ok3 && modeStr == "scientific" {
			mode = database.ModeScientific
		}

		if ok1 && ok2 {
			_, err := h.calcRepo.InsertCalculation(expression, result, mode)
			if err != nil {
				// Log error but continue with other entries
				continue
			}
		}
	}

	return nil
}

// GetStats returns statistics about the history
func (h *HistoryManager) GetStats() (map[string]interface{}, error) {
	history, err := h.GetHistory()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_calculations": len(history),
		"by_mode": map[string]int{
			"standard":   0,
			"scientific": 0,
		},
	}

	for _, calc := range history {
		switch calc.Mode {
		case database.ModeStandard:
			stats["by_mode"].(map[string]int)["standard"]++
		case database.ModeScientific:
			stats["by_mode"].(map[string]int)["scientific"]++
		}
	}

	return stats, nil
}
