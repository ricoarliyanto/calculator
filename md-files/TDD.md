# TDD - Test Driven Development Strategy

## Overview
Dokumen ini menjelaskan strategi Test Driven Development (TDD) untuk aplikasi Calculator. TDD akan digunakan untuk memastikan kualitas kode dan meminimalkan bug sejak awal pengembangan.

## TDD Philosophy
TDD mengikuti siklus **Red-Green-Refactor**:
1. **Red**: Tulis test yang gagal (test belum ada implementasi)
2. **Green**: Tulis kode minimal agar test lulus
3. **Refactor**: Refactor kode untuk meningkatkan kualitas tanpa mengubah perilaku

## Testing Structure

### Directory Structure
```
calculator/
├── internal/
│   ├── calculator/
│   │   ├── standard.go
│   │   ├── scientific.go
│   │   ├── evaluator.go
│   │   └── calculator_test.go
│   ├── database/
│   │   ├── sqlite.go
│   │   ├── models.go
│   │   └── database_test.go
│   ├── history/
│   │   ├── history.go
│   │   └── history_test.go
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   └── ui/
│       ├── cli.go
│       └── cli_test.go
└── pkg/
    └── utils/
        ├── helpers.go
        └── helpers_test.go
```

## Test Categories

### 1. Unit Tests
Test untuk fungsi/method individual dalam isolasi.

#### Coverage Target: > 80%

#### Test Cases for Calculator

##### Standard Calculator Tests
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     float64
        expected float64
    }{
        {"positive numbers", 5, 3, 8},
        {"negative numbers", -5, -3, -8},
        {"mixed signs", 5, -3, 2},
        {"decimals", 5.5, 3.3, 8.8},
        {"zero", 5, 0, 5},
    }
    // ... implementation
}

func TestSubtract(t *testing.T) { /* ... */ }
func TestMultiply(t *testing.T) { /* ... */ }
func TestDivide(t *testing.T) {
    // Test normal division
    // Test division by zero (should return error)
    // Test decimal division
}

func TestPercentage(t *testing.T) { /* ... */ }
func TestParentheses(t *testing.T) { /* ... */ }
```

##### Scientific Calculator Tests
```go
func TestSin(t *testing.T) {
    tests := []struct {
        angle    float64
        unit     string // "deg" or "rad"
        expected float64
    }{
        {0, "rad", 0},
        {math.Pi / 2, "rad", 1},
        {90, "deg", 1},
    }
    // ... implementation
}

func TestCos(t *testing.T) { /* ... */ }
func TestTan(t *testing.T) { /* ... */ }
func TestLog(t *testing.T) { /* ... */ }
func TestLn(t *testing.T) { /* ... */ }
func TestExp(t *testing.T) { /* ... */ }
func TestSqrt(t *testing.T) { /* ... */ }
func TestFactorial(t *testing.T) {
    // Test positive integers
    // Test zero (should return 1)
    // Test negative (should return error)
    // Test non-integer (should return error)
}

func TestPower(t *testing.T) { /* ... */ }
func TestConstants(t *testing.T) {
    // Test PI value
    // Test E value
}
```

##### Expression Parser Tests
```go
func TestParseExpression(t *testing.T) {
    tests := []struct {
        expression string
        expected   float64
        shouldErr  bool
    }{
        {"2 + 2", 4, false},
        {"2 * (3 + 4)", 14, false},
        {"10 / 2", 5, false},
        {"invalid", 0, true},
        {"10 / 0", 0, true}, // division by zero
    }
    // ... implementation
}

func TestOperatorPrecedence(t *testing.T) {
    // Test that * and / have higher precedence than + and -
    // Test parentheses override precedence
}

func TestComplexExpression(t *testing.T) {
    // Test nested parentheses
    // Test multiple operators
    // Test scientific functions in expressions
}
```

#### Test Cases for Database

##### Database Connection Tests
```go
func TestDatabaseConnection(t *testing.T) {
    // Test successful connection
    // Test connection failure (invalid path)
    // Test connection pooling
}

func TestSchemaCreation(t *testing.T) {
    // Test tables are created
    // Test indexes are created
    // Test migration works
}
```

##### CRUD Operations Tests
```go
func TestInsertCalculation(t *testing.T) {
    // Test successful insert
    // Test insert with empty expression (should fail)
    // Test auto-increment ID
}

func TestGetAllCalculations(t *testing.T) {
    // Test empty database
    // Test with multiple records
    // Test ordering by date
}

func TestGetCalculationsByFilter(t *testing.T) {
    // Test filter by mode
    // Test filter by date
    // Test filter by operator
    // Test combined filters
}

func TestDeleteCalculation(t *testing.T) {
    // Test successful delete
    // Test delete non-existent ID
}

func TestClearAllCalculations(t *testing.T) {
    // Test clear operation
    // Test database is empty after clear
}

func TestExportCalculations(t *testing.T) {
    // Test export to CSV
    // Test export to JSON
    // Test export empty database
}
```

##### Settings Management Tests
```go
func TestSaveSetting(t *testing.T) {
    // Test save new setting
    // Test update existing setting
    // Test save empty key (should fail)
}

func TestLoadSetting(t *testing.T) {
    // Test load existing setting
    // Test load non-existent setting (return default)
}

func TestDefaultSettings(t *testing.T) {
    // Test default values are applied
    // Test default mode is standard
}
```

#### Test Cases for Config

```go
func TestLoadConfig(t *testing.T) {
    // Test load valid config
    // Test load invalid config (use defaults)
    // Test load non-existent config (create default)
}

func TestSaveConfig(t *testing.T) {
    // Test save config
    // Test config persistence
}

func TestConfigValidation(t *testing.T) {
    // Test valid config
    // Test invalid decimal places (negative)
    // Test invalid angle unit
}
```

### 2. Integration Tests
Test untuk memverifikasi bahwa komponen bekerja bersama dengan benar.

#### Integration Test Cases

```go
func TestCompleteCalculationFlow(t *testing.T) {
    // 1. Load config
    // 2. Connect to database
    // 3. Perform calculation
    // 4. Save to database
    // 5. Retrieve from database
    // 6. Verify result matches
}

func TestHistoryPersistence(t *testing.T) {
    // 1. Perform multiple calculations
    // 2. Save to database
    // 3. Close database
    // 4. Reopen database
    // 5. Verify history is preserved
}

func TestModeSwitching(t *testing.T) {
    // 1. Start in standard mode
    // 2. Perform standard calculation
    // 3. Switch to scientific mode
    // 4. Perform scientific calculation
    // 5. Switch back to standard
    // 6. Verify mode persistence
}

func TestExportWorkflow(t *testing.T) {
    // 1. Add calculations to database
    // 2. Export to CSV
    // 3. Verify CSV file content
    // 4. Export to JSON
    // 5. Verify JSON file content
}
```

### 3. Edge Cases Tests
Test untuk kondisi batas dan kasus extreme.

```go
func TestEdgeCases(t *testing.T) {
    // Large numbers
    t.Run("Large Numbers", func(t *testing.T) {
        // Test with very large numbers
        // Test with very small numbers
    })

    // Invalid input
    t.Run("Invalid Input", func(t *testing.T) {
        // Test empty input
        // Test special characters
        // Test malformed expressions
    })

    // Boundary values
    t.Run("Boundary Values", func(t *testing.T) {
        // Test factorial(0) = 1
        // Test factorial(1) = 1
        // Test log(1) = 0
        // Test sqrt(0) = 0
    })

    // Precision
    t.Run("Precision", func(t *testing.T) {
        // Test decimal precision
        // Test floating point errors
        // Test rounding
    })
}
```

### 4. Performance Tests
Test untuk memverifikasi performa aplikasi.

```go
func BenchmarkCalculation(b *testing.B) {
    calculator := NewCalculator()
    for i := 0; i < b.N; i++ {
        calculator.Evaluate("2 + 2")
    }
}

func BenchmarkDatabaseInsert(b *testing.B) {
    db := SetupTestDatabase()
    defer db.Close()
    for i := 0; i < b.N; i++ {
        db.InsertCalculation("test", "result", "standard")
    }
}

func BenchmarkDatabaseQuery(b *testing.B) {
    db := SetupTestDatabase()
    defer db.Close()
    // Pre-populate with data
    for i := 0; i < 1000; i++ {
        db.InsertCalculation("test", "result", "standard")
    }
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        db.GetAllCalculations()
    }
}
```

## Test Utilities

### Test Helper Functions

```go
// setupTestDatabase creates a temporary database for testing
func setupTestDatabase(t *testing.T) *Database {
    db, err := NewDatabase(":memory:")
    if err != nil {
        t.Fatalf("Failed to create test database: %v", err)
    }
    return db
}

// cleanupTestDatabase removes test database
func cleanupTestDatabase(t *testing.T, db *Database) {
    err := db.Close()
    if err != nil {
        t.Logf("Warning: Failed to close test database: %v", err)
    }
}

// assertFloatEqual compares float64 with tolerance
func assertFloatEqual(t *testing.T, expected, actual float64, tolerance float64) {
    diff := math.Abs(expected - actual)
    if diff > tolerance {
        t.Errorf("Expected %v, got %v (diff: %v)", expected, actual, diff)
    }
}
```

## Test Execution

### Running All Tests
```bash
go test ./...
```

### Running Tests with Coverage
```bash
go test -cover ./...
```

### Running Tests with Coverage Report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Running Specific Package Tests
```bash
go test ./internal/calculator
go test ./internal/database
```

### Running Specific Test
```bash
go test -run TestAdd ./internal/calculator
```

### Running Benchmarks
```bash
go test -bench=. ./...
```

### Running Tests with Verbose Output
```bash
go test -v ./...
```

## Test Coverage Goals

### Minimum Coverage by Package
- `internal/calculator`: 90%
- `internal/database`: 85%
- `internal/history`: 85%
- `internal/config`: 80%
- `internal/ui`: 75%
- `pkg/utils`: 90%

### Overall Coverage Target: > 80%

## Continuous Testing Strategy

### Pre-Commit Testing
Sebelum commit, jalankan:
```bash
go test ./...
go vet ./...
golangci-lint run
```

### CI/CD Pipeline
1. Run all tests
2. Check coverage (>80%)
3. Run linters
4. Run benchmarks
5. Build binaries

## Testing Best Practices

### DO's
- ✅ Tulis test sebelum implementasi (TDD)
- ✅ Buat test yang independent dan isolated
- ✅ Gunakan descriptive test names
- ✅ Test edge cases dan error conditions
- ✅ Mock external dependencies
- ✅ Keep tests fast
- ✅ Use table-driven tests for multiple cases

### DON'Ts
- ❌ Jangan test library functions (trust the library)
- ❌ Jangan buat test yang tergantung urutan eksekusi
- ❌ Jangan hardcode test data yang fragile
- ❌ Jangan abaikan test yang flaky
- ❌ Jangan buat test yang terlalu kompleks
- ❌ Jangan test private functions (test through public API)

## Test Checklist for Each Feature

Untuk setiap feature baru, pastikan:
- [ ] Unit tests ditulis sebelum implementasi
- [ ] Unit tests mencakup happy path
- [ ] Unit tests mencakup error cases
- [ ] Unit tests mencakup edge cases
- [ ] Integration tests ditulis
- [ ] Test coverage > 80%
- [ ] Tests pass consistently
- [ ] Tests are fast (< 1s per test)
- [ ] Test names are descriptive
- [ ] Test documentation is clear

## Mocking Strategy

### When to Mock
- External services (API calls)
- File system operations (for pure logic tests)
- Time (for time-dependent tests)
- Database (for unit tests, use in-memory DB for integration)

### Mocking Libraries
- `github.com/stretchr/testify/mock` - Mocking framework
- In-memory SQLite for database tests

## Test Data Management

### Test Fixtures
```go
// test_fixtures.go
var testCalculations = []Calculation{
    {Expression: "2 + 2", Result: "4", Mode: "standard"},
    {Expression: "sin(90)", Result: "1", Mode: "scientific"},
}
```

### Test Database Seeding
```go
func seedTestDatabase(db *Database) error {
    for _, calc := range testCalculations {
        err := db.InsertCalculation(calc.Expression, calc.Result, calc.Mode)
        if err != nil {
            return err
        }
    }
    return nil
}
```

## Conclusion

Dengan mengikuti strategi TDD ini, aplikasi Calculator akan memiliki:
- Kualitas kode yang tinggi
- Bug yang minimal
- Dokumentasi melalui test
- Confidence untuk refactoring
- Kemudahan maintenance
