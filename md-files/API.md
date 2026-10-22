# API Documentation - Calculator Application

## Overview
Dokumen ini mendokumentasikan API internal dan fungsi-fungsi publik yang tersedia di aplikasi Calculator.

## Package: internal/calculator

### Standard Calculator

#### Add
```go
func Add(a, b float64) float64
```
**Description**: Menambahkan dua angka

**Parameters**:
- `a` (float64): Angka pertama
- `b` (float64): Angka kedua

**Returns**: float64 - Hasil penjumlahan

**Example**:
```go
result := Add(5, 3) // returns 8
```

#### Subtract
```go
func Subtract(a, b float64) float64
```
**Description**: Mengurangkan dua angka

**Parameters**:
- `a` (float64): Angka pertama
- `b` (float64): Angka kedua

**Returns**: float64 - Hasil pengurangan

**Example**:
```go
result := Subtract(10, 4) // returns 6
```

#### Multiply
```go
func Multiply(a, b float64) float64
```
**Description**: Mengalikan dua angka

**Parameters**:
- `a` (float64): Angka pertama
- `b` (float64): Angka kedua

**Returns**: float64 - Hasil perkalian

**Example**:
```go
result := Multiply(6, 7) // returns 42
```

#### Divide
```go
func Divide(a, b float64) (float64, error)
```
**Description**: Membagi dua angka

**Parameters**:
- `a` (float64): Angka pertama (dividend)
- `b` (float64): Angka kedua (divisor)

**Returns**:
- float64 - Hasil pembagian
- error - Error jika pembagian dengan nol

**Errors**:
- `CALC_DIVISION_BY_ZERO` - Jika b adalah 0

**Example**:
```go
result, err := Divide(10, 2) // returns 5, nil
result, err := Divide(10, 0) // returns 0, error
```

#### Percentage
```go
func Percentage(value, percent float64) float64
```
**Description**: Menghitung persentase dari nilai

**Parameters**:
- `value` (float64): Nilai dasar
- `percent` (float64): Persentase yang dihitung

**Returns**: float64 - Hasil persentase

**Example**:
```go
result := Percentage(100, 25) // returns 25
```

### Scientific Calculator

#### Sin
```go
func Sin(angle float64, unit string) (float64, error)
```
**Description**: Menghitung sinus dari sudut

**Parameters**:
- `angle` (float64): Sudut dalam derajat atau radian
- `unit` (string): Unit sudut ("deg" atau "rad")

**Returns**:
- float64 - Nilai sinus
- error - Error jika unit tidak valid

**Errors**:
- `CALC_INVALID_UNIT` - Jika unit bukan "deg" atau "rad"

**Example**:
```go
result, err := Sin(90, "deg") // returns 1, nil
result, err := Sin(math.Pi/2, "rad") // returns 1, nil
```

#### Cos
```go
func Cos(angle float64, unit string) (float64, error)
```
**Description**: Menghitung cosinus dari sudut

**Parameters**:
- `angle` (float64): Sudut dalam derajat atau radian
- `unit` (string): Unit sudut ("deg" atau "rad")

**Returns**:
- float64 - Nilai cosinus
- error - Error jika unit tidak valid

**Example**:
```go
result, err := Cos(0, "deg") // returns 1, nil
```

#### Tan
```go
func Tan(angle float64, unit string) (float64, error)
```
**Description**: Menghitung tangen dari sudut

**Parameters**:
- `angle` (float64): Sudut dalam derajat atau radian
- `unit` (string): Unit sudut ("deg" atau "rad")

**Returns**:
- float64 - Nilai tangen
- error - Error jika unit tidak valid atau sudut invalid

**Errors**:
- `CALC_INVALID_UNIT` - Jika unit bukan "deg" atau "rad"
- `CALC_DOMAIN_ERROR` - Jika sudut menyebabkan undefined (cos=0)

**Example**:
```go
result, err := Tan(45, "deg") // returns 1, nil
```

#### Log
```go
func Log(value float64) (float64, error)
```
**Description**: Menghitung logaritma base 10

**Parameters**:
- `value` (float64): Nilai yang dihitung

**Returns**:
- float64 - Nilai logaritma
- error - Error jika value <= 0

**Errors**:
- `CALC_DOMAIN_ERROR` - Jika value <= 0

**Example**:
```go
result, err := Log(100) // returns 2, nil
```

#### Ln
```go
func Ln(value float64) (float64, error)
```
**Description**: Menghitung logaritma natural (base e)

**Parameters**:
- `value` (float64): Nilai yang dihitung

**Returns**:
- float64 - Nilai logaritma natural
- error - Error jika value <= 0

**Errors**:
- `CALC_DOMAIN_ERROR` - Jika value <= 0

**Example**:
```go
result, err := Ln(math.E) // returns 1, nil
```

#### Exp
```go
func Exp(value float64) float64
```
**Description**: Menghitung e pangkat value

**Parameters**:
- `value` (float64): Eksponen

**Returns**: float64 - Hasil e^value

**Example**:
```go
result := Exp(1) // returns 2.718281828459045
```

#### Power
```go
func Power(base, exponent float64) (float64, error)
```
**Description**: Menghitung base pangkat exponent

**Parameters**:
- `base` (float64): Basis
- `exponent` (float64): Eksponen

**Returns**:
- float64 - Hasil pangkat
- error - Error untuk kasus invalid

**Errors**:
- `CALC_DOMAIN_ERROR` - Jika base < 0 dan exponent bukan integer

**Example**:
```go
result, err := Power(2, 3) // returns 8, nil
```

#### Sqrt
```go
func Sqrt(value float64) (float64, error)
```
**Description**: Menghitung akar kuadrat

**Parameters**:
- `value` (float64): Nilai yang dihitung

**Returns**:
- float64 - Nilai akar kuadrat
- error - Error jika value < 0

**Errors**:
- `CALC_DOMAIN_ERROR` - Jika value < 0

**Example**:
```go
result, err := Sqrt(16) // returns 4, nil
```

#### Factorial
```go
func Factorial(n int) (int, error)
```
**Description**: Menghitung faktorial dari n

**Parameters**:
- `n` (int): Bilangan bulat non-negatif

**Returns**:
- int - Nilai faktorial
- error - Error jika n < 0

**Errors**:
- `CALC_DOMAIN_ERROR` - Jika n < 0

**Example**:
```go
result, err := Factorial(5) // returns 120, nil
```

### Expression Evaluator

#### Evaluate
```go
func Evaluate(expression string, mode string) (float64, error)
```
**Description**: Mengevaluasi ekspresi matematika

**Parameters**:
- `expression` (string): Ekspresi matematika
- `mode` (string): Mode calculator ("standard" atau "scientific")

**Returns**:
- float64 - Hasil evaluasi
- error - Error jika ekspresi invalid

**Errors**:
- `CALC_INVALID_EXPRESSION` - Ekspresi syntax error
- `CALC_DIVISION_BY_ZERO` - Pembagian dengan nol
- `CALC_INVALID_FUNCTION` - Fungsi tidak valid untuk mode
- `CALC_DOMAIN_ERROR` - Domain error

**Example**:
```go
result, err := Evaluate("2 + 2 * 3", "standard") // returns 8, nil
result, err := Evaluate("sin(90)", "scientific") // returns 1, nil
```

#### ParseExpression
```go
func ParseExpression(expression string) ([]Token, error)
```
**Description**: Parsing ekspresi menjadi tokens

**Parameters**:
- `expression` (string): Ekspresi matematika

**Returns**:
- []Token - List of tokens
- error - Error jika parsing gagal

**Example**:
```go
tokens, err := ParseExpression("2 + 2")
// tokens: [{NUMBER, 2}, {OPERATOR, +}, {NUMBER, 2}]
```

### Memory Functions

#### MemoryClear
```go
func MemoryClear()
```
**Description**: Menghapus nilai memory

#### MemoryRecall
```go
func MemoryRecall() float64
```
**Description**: Mengambil nilai dari memory

**Returns**: float64 - Nilai memory

#### MemoryStore
```go
func MemoryStore(value float64)
```
**Description**: Menyimpan nilai ke memory

**Parameters**:
- `value` (float64): Nilai yang disimpan

#### MemoryAdd
```go
func MemoryAdd(value float64)
```
**Description**: Menambahkan nilai ke memory

**Parameters**:
- `value` (float64): Nilai yang ditambahkan

#### MemorySubtract
```go
func MemorySubtract(value float64)
```
**Description**: Mengurangkan nilai dari memory

**Parameters**:
- `value` (float64): Nilai yang dikurangkan

### Undo/Redo

#### Undo
```go
func Undo() (float64, error)
```
**Description**: Mengembalikan ke perhitungan sebelumnya

**Returns**:
- float64 - Hasil perhitungan sebelumnya
- error - Error jika tidak ada history undo

#### Redo
```go
func Redo() (float64, error)
```
**Description**: Mengulang perhitungan yang di-undo

**Returns**:
- float64 - Hasil perhitungan
- error - Error jika tidak ada history redo

---

## Package: internal/database

### Database Connection

#### NewDatabase
```go
func NewDatabase(path string) (*Database, error)
```
**Description**: Membuat koneksi database baru

**Parameters**:
- `path` (string): Path ke file database (":memory:" untuk in-memory)

**Returns**:
- *Database - Instance database
- error - Error jika koneksi gagal

**Example**:
```go
db, err := NewDatabase("data/calculator.db")
db, err := NewDatabase(":memory:") // in-memory database
```

#### Close
```go
func (db *Database) Close() error
```
**Description**: Menutup koneksi database

**Returns**: error - Error jika close gagal

### CRUD Operations

#### InsertCalculation
```go
func (db *Database) InsertCalculation(expression, result, mode string) error
```
**Description**: Menyimpan perhitungan ke database

**Parameters**:
- `expression` (string): Ekspresi matematika
- `result` (string): Hasil perhitungan
- `mode` (string): Mode calculator

**Returns**: error - Error jika insert gagal

**Example**:
```go
err := db.InsertCalculation("2 + 2", "4", "standard")
```

#### GetAllCalculations
```go
func (db *Database) GetAllCalculations() ([]Calculation, error)
```
**Description**: Mengambil semua history perhitungan

**Returns**:
- []Calculation - List of calculations
- error - Error jika query gagal

**Example**:
```go
calculations, err := db.GetAllCalculations()
```

#### GetCalculationsByFilter
```go
func (db *Database) GetCalculationsByFilter(filter Filter) ([]Calculation, error)
```
**Description**: Mengambil history dengan filter

**Parameters**:
- `filter` (Filter): Filter criteria

**Returns**:
- []Calculation - List of calculations
- error - Error jika query gagal

**Example**:
```go
filter := Filter{
    Mode: "scientific",
    Limit: 10,
}
calculations, err := db.GetCalculationsByFilter(filter)
```

#### DeleteCalculation
```go
func (db *Database) DeleteCalculation(id int) error
```
**Description**: Menghapus perhitungan berdasarkan ID

**Parameters**:
- `id` (int): ID perhitungan

**Returns**: error - Error jika delete gagal

**Example**:
```go
err := db.DeleteCalculation(5)
```

#### ClearAllCalculations
```go
func (db *Database) ClearAllCalculations() error
```
**Description**: Menghapus semua history perhitungan

**Returns**: error - Error jika clear gagal

### Settings

#### SaveSetting
```go
func (db *Database) SaveSetting(key, value string) error
```
**Description**: Menyimpan setting

**Parameters**:
- `key` (string): Key setting
- `value` (string): Value setting

**Returns**: error - Error jika save gagal

#### LoadSetting
```go
func (db *Database) LoadSetting(key string) (string, error)
```
**Description**: Mengambil setting

**Parameters**:
- `key` (string): Key setting

**Returns**:
- string - Value setting
- error - Error jika load gagal

---

## Package: internal/history

### History Management

#### GetHistory
```go
func GetHistory(filter Filter) ([]Calculation, error)
```
**Description**: Mengambil history dengan filter

**Parameters**:
- `filter` (Filter): Filter criteria

**Returns**:
- []Calculation - List of calculations
- error - Error jika gagal

#### FilterHistory
```go
func FilterHistory(calculations []Calculation, filter Filter) []Calculation
```
**Description**: Filter list of calculations

**Parameters**:
- `calculations` ([]Calculation): List calculations
- `filter` (Filter): Filter criteria

**Returns**: []Calculation - Filtered calculations

### Export

#### ExportToCSV
```go
func ExportToCSV(calculations []Calculation, filename string) error
```
**Description**: Export history ke file CSV

**Parameters**:
- `calculations` ([]Calculation): List calculations
- `filename` (string): Nama file output

**Returns**: error - Error jika export gagal

**Example**:
```go
err := ExportToCSV(calculations, "history.csv")
```

#### ExportToJSON
```go
func ExportToJSON(calculations []Calculation, filename string) error
```
**Description**: Export history ke file JSON

**Parameters**:
- `calculations` ([]Calculation): List calculations
- `filename` (string): Nama file output

**Returns**: error - Error jika export gagal

**Example**:
```go
err := ExportToJSON(calculations, "history.json")
```

---

## Package: internal/config

### Config Management

#### LoadConfig
```go
func LoadConfig(path string) (*Config, error)
```
**Description**: Load config dari file

**Parameters**:
- `path` (string): Path ke config file

**Returns**:
- *Config - Config object
- error - Error jika load gagal

**Example**:
```go
config, err := LoadConfig("config.yaml")
```

#### SaveConfig
```go
func SaveConfig(config *Config, path string) error
```
**Description**: Save config ke file

**Parameters**:
- `config` (*Config): Config object
- `path` (string): Path ke config file

**Returns**: error - Error jika save gagal

#### GetDefaultConfig
```go
func GetDefaultConfig() *Config
```
**Description**: Mendapatkan config dengan default values

**Returns**: *Config - Config dengan default values

#### ValidateConfig
```go
func ValidateConfig(config *Config) error
```
**Description**: Validasi config

**Parameters**:
- `config` (*Config): Config object

**Returns**: error - Error jika config invalid

---

## Package: internal/ui

### CLI Interface

#### Run
```go
func Run() error
```
**Description**: Menjalankan CLI application

**Returns**: error - Error jika terjadi error

#### StartInteractiveMode
```go
func StartInteractiveMode(mode string) error
```
**Description**: Memulai mode interaktif

**Parameters**:
- `mode` (string): Mode awal ("standard" atau "scientific")

**Returns**: error - Error jika terjadi error

#### ExecuteCommand
```go
func ExecuteCommand(args []string) error
```
**Description**: Execute command dari arguments

**Parameters**:
- `args` ([]string): Command arguments

**Returns**: error - Error jika command gagal

---

## Package: pkg/utils

### Formatters

#### FormatResult
```go
func FormatResult(result float64, decimalPlaces int) string
```
**Description**: Format hasil perhitungan

**Parameters**:
- `result` (float64): Hasil perhitungan
- `decimalPlaces` (int): Jumlah desimal

**Returns**: string - Formatted result

**Example**:
```go
formatted := FormatResult(3.14159, 2) // returns "3.14"
```

#### FormatTimestamp
```go
func FormatTimestamp(t time.Time) string
```
**Description**: Format timestamp

**Parameters**:
- `t` (time.Time): Timestamp

**Returns**: string - Formatted timestamp

### Validators

#### ValidateExpression
```go
func ValidateExpression(expression string) error
```
**Description**: Validasi ekspresi matematika

**Parameters**:
- `expression` (string): Ekspresi yang divalidasi

**Returns**: error - Error jika ekspresi invalid

#### ValidateMode
```go
func ValidateMode(mode string) error
```
**Description**: Validasi mode

**Parameters**:
- `mode` (string): Mode yang divalidasi

**Returns**: error - Error jika mode invalid

---

## Data Structures

### Calculation
```go
type Calculation struct {
    ID         int       `json:"id"`
    Expression string    `json:"expression"`
    Result     string    `json:"result"`
    Mode       string    `json:"mode"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
}
```

### Setting
```go
type Setting struct {
    Key       string    `json:"key"`
    Value     string    `json:"value"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Filter
```go
type Filter struct {
    Mode      string    `json:"mode,omitempty"`
    FromDate  time.Time `json:"from_date,omitempty"`
    ToDate    time.Time `json:"to_date,omitempty"`
    Operator  string    `json:"operator,omitempty"`
    Limit     int       `json:"limit,omitempty"`
    Offset    int       `json:"offset,omitempty"`
}
```

### Config
```go
type Config struct {
    Calculator struct {
        Mode           string `yaml:"mode"`
        AngleUnit      string `yaml:"angle-unit"`
        DecimalPlaces  int    `yaml:"decimal-places"`
        HistoryLimit   int    `yaml:"history-limit"`
        AutoSave       bool   `yaml:"auto-save"`
        ShowTimestamps bool   `yaml:"show-timestamps"`
    } `yaml:"calculator"`
}
```

### Token
```go
type Token struct {
    Type  TokenType
    Value string
}

type TokenType int

const (
    TokenNumber TokenType = iota
    TokenOperator
    TokenFunction
    TokenParenthesis
    TokenComma
)
```

---

## Error Types

### CalculatorError
```go
type CalculatorError struct {
    Code    string
    Message string
    Err     error
}
```

**Error Codes**:
- `CALC_INVALID_EXPRESSION` - Invalid expression syntax
- `CALC_DIVISION_BY_ZERO` - Division by zero
- `CALC_INVALID_FUNCTION` - Invalid function name
- `CALC_DOMAIN_ERROR` - Domain error
- `CALC_INVALID_UNIT` - Invalid angle unit
- `DB_CONNECTION_ERROR` - Database connection error
- `DB_QUERY_ERROR` - Database query error
- `CONFIG_ERROR` - Configuration error
- `IO_ERROR` - File I/O error

---

## Usage Examples

### Complete Calculation Flow
```go
package main

import (
    "fmt"
    "calculator/internal/calculator"
    "calculator/internal/database"
    "calculator/internal/config"
)

func main() {
    // Load config
    cfg, err := config.LoadConfig("config.yaml")
    if err != nil {
        cfg = config.GetDefaultConfig()
    }

    // Connect to database
    db, err := database.NewDatabase("data/calculator.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Perform calculation
    result, err := calculator.Evaluate("2 + 2 * 3", "standard")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Result: %f\n", result)

    // Save to database
    err = db.InsertCalculation("2 + 2 * 3", fmt.Sprintf("%f", result), "standard")
    if err != nil {
        panic(err)
    }
}
```

### Scientific Calculation
```go
// Calculate sin(90) in degrees
result, err := calculator.Sin(90, "deg")
if err != nil {
    panic(err)
}
fmt.Printf("sin(90°) = %f\n", result)

// Calculate logarithm
result, err = calculator.Log(100)
if err != nil {
    panic(err)
}
fmt.Printf("log(100) = %f\n", result)
```

### History Export
```go
// Get calculations
calculations, err := db.GetAllCalculations()
if err != nil {
    panic(err)
}

// Export to CSV
err = history.ExportToCSV(calculations, "history.csv")
if err != nil {
    panic(err)
}

// Export to JSON
err = history.ExportToJSON(calculations, "history.json")
if err != nil {
    panic(err)
}
```
