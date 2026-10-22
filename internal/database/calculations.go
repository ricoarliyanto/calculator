package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// CalculationRepository handles calculation database operations
type CalculationRepository struct {
	db *Database
}

// NewCalculationRepository creates a new calculation repository
func NewCalculationRepository(db *Database) *CalculationRepository {
	return &CalculationRepository{db: db}
}

// InsertCalculation inserts a new calculation record
func (r *CalculationRepository) InsertCalculation(expression, result string, mode CalculationMode) (*Calculation, error) {
	conn, err := r.db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	query := `INSERT INTO calculations (expression, result, mode, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	res, err := conn.Exec(query, expression, result, mode, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert calculation: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	calc := &Calculation{
		ID:         id,
		Expression: expression,
		Result:     result,
		Mode:       mode,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return calc, nil
}

// GetAllCalculations retrieves all calculations ordered by created_at descending
func (r *CalculationRepository) GetAllCalculations() ([]Calculation, error) {
	conn, err := r.db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	rows, err := conn.Query(`SELECT id, expression, result, mode, created_at, updated_at FROM calculations ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to query calculations: %w", err)
	}
	defer rows.Close()

	var calculations []Calculation
	for rows.Next() {
		calc := Calculation{}
		var modeStr string

		if err := rows.Scan(&calc.ID, &calc.Expression, &calc.Result, &modeStr, &calc.CreatedAt, &calc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan calculation: %w", err)
		}

		calc.Mode = CalculationMode(modeStr)
		calculations = append(calculations, calc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating calculations: %w", err)
	}

	return calculations, nil
}

// GetCalculationsByFilter retrieves calculations based on filter criteria
func (r *CalculationRepository) GetCalculationsByFilter(filter CalculationFilter) ([]Calculation, error) {
	conn, err := r.db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	var query strings.Builder
	var args []interface{}

	query.WriteString(`SELECT id, expression, result, mode, created_at, updated_at FROM calculations WHERE 1=1`)

	if filter.Mode != nil {
		query.WriteString(` AND mode = ?`)
		args = append(args, string(*filter.Mode))
	}

	if filter.StartDate != nil {
		query.WriteString(` AND created_at >= ?`)
		args = append(args, filter.StartDate.Format(time.RFC3339))
	}

	if filter.EndDate != nil {
		query.WriteString(` AND created_at <= ?`)
		args = append(args, filter.EndDate.Format(time.RFC3339))
	}

	if filter.Search != nil && *filter.Search != "" {
		query.WriteString(` AND (expression LIKE ? OR result LIKE ?)`)
		searchPattern := fmt.Sprintf("%%%s%%", *filter.Search)
		args = append(args, searchPattern, searchPattern)
	}

	query.WriteString(` ORDER BY created_at DESC`)

	rows, err := conn.Query(query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query filtered calculations: %w", err)
	}
	defer rows.Close()

	var calculations []Calculation
	for rows.Next() {
		calc := Calculation{}
		var modeStr string

		if err := rows.Scan(&calc.ID, &calc.Expression, &calc.Result, &modeStr, &calc.CreatedAt, &calc.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan calculation: %w", err)
		}

		calc.Mode = CalculationMode(modeStr)
		calculations = append(calculations, calc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating calculations: %w", err)
	}

	return calculations, nil
}

// GetCalculationByID retrieves a single calculation by ID
func (r *CalculationRepository) GetCalculationByID(id int64) (*Calculation, error) {
	conn, err := r.db.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	calc := &Calculation{}
	var modeStr string

	if err := conn.QueryRow(`SELECT id, expression, result, mode, created_at, updated_at FROM calculations WHERE id = ?`, id).Scan(
		&calc.ID, &calc.Expression, &calc.Result, &modeStr, &calc.CreatedAt, &calc.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query calculation by id: %w", err)
	}

	calc.Mode = CalculationMode(modeStr)
	return calc, nil
}

// DeleteCalculation deletes a calculation by ID
func (r *CalculationRepository) DeleteCalculation(id int64) error {
	conn, err := r.db.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if _, err := conn.Exec(`DELETE FROM calculations WHERE id = ?`, id); err != nil {
		return fmt.Errorf("failed to delete calculation: %w", err)
	}

	return nil
}

// ClearAllCalculations deletes all calculations
func (r *CalculationRepository) ClearAllCalculations() error {
	conn, err := r.db.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	if _, err := conn.Exec(`DELETE FROM calculations`); err != nil {
		return fmt.Errorf("failed to clear all calculations: %w", err)
	}

	return nil
}

// ExportCalculations exports calculations in the specified format
func (r *CalculationRepository) ExportCalculations(format ExportFormat) ([]byte, error) {
	calculations, err := r.GetAllCalculations()
	if err != nil {
		return nil, fmt.Errorf("failed to get calculations for export: %w", err)
	}

	switch format {
	case ExportJSON:
		var data []map[string]interface{}
		for _, calc := range calculations {
			data = append(data, map[string]interface{}{
				"id":         calc.ID,
				"expression": calc.Expression,
				"result":     calc.Result,
				"mode":       calc.Mode,
				"created_at": calc.CreatedAt.Format(time.RFC3339),
			})
		}
		return json.MarshalIndent(data, "", "  ")
	case ExportCSV:
		var sb strings.Builder
		sb.WriteString("ID,Expression,Result,Mode,CreatedAt\n")
		for _, calc := range calculations {
			sb.WriteString(fmt.Sprintf("%d,\"%s\",\"%s\",%s,%s\n",
				calc.ID,
				strings.ReplaceAll(calc.Expression, "\"", "\"\" "),
				strings.ReplaceAll(calc.Result, "\"", "\"\" "),
				calc.Mode,
				calc.CreatedAt.Format(time.RFC3339)))
		}
		return []byte(sb.String()), nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}
