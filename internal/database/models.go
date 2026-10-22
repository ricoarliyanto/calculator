package database

import (
	"time"
)

// CalculationMode represents the calculator mode
type CalculationMode string

const (
	ModeStandard   CalculationMode = "standard"
	ModeScientific CalculationMode = "scientific"
)

// Calculation represents a single calculation record
type Calculation struct {
	ID          int64           `json:"id" db:"id"`
	Expression  string          `json:"expression" db:"expression"`
	Result      string          `json:"result" db:"result"`
	Mode        CalculationMode `json:"mode" db:"mode"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// Setting represents an application setting
type Setting struct {
	ID        int64     `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CalculationFilter represents filter criteria for calculations
type CalculationFilter struct {
	Mode      *CalculationMode
	StartDate *time.Time
	EndDate   *time.Time
	Search    *string
}

// ExportFormat represents the format for exporting calculations
type ExportFormat string

const (
	ExportJSON ExportFormat = "json"
	ExportCSV  ExportFormat = "csv"
)
