# ARCHITECTURE - Calculator Application

## Overview
Dokumen ini menjelaskan arsitektur teknis dari aplikasi Calculator, termasuk desain komponen, data flow, dan pattern yang digunakan.

## Architecture Pattern
Aplikasi menggunakan **Layered Architecture** dengan pemisahan concern yang jelas:

```
┌─────────────────────────────────────────────────┐
│           Presentation Layer (CLI)              │
│                  internal/ui                    │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│          Application Layer (Logic)              │
│              internal/calculator                 │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│           Data Access Layer (Database)          │
│              internal/database                   │
└────────────────────┬────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────┐
│              Data Layer (SQLite)                 │
│                   data/calculator.db             │
└─────────────────────────────────────────────────┘
```

## Component Architecture

### 1. Presentation Layer (CLI)
**Location**: `internal/ui/`

**Responsibilities**:
- Handle user input dari command line
- Display output ke user
- Parse command line arguments
- Manage interactive REPL loop
- Format dan presentasi data

**Key Components**:
- `cli.go` - Main CLI interface
- Command handlers untuk setiap command
- Input parser dan validator
- Output formatter

**Dependencies**:
- Application layer (calculator logic)
- Config layer
- History layer

### 2. Application Layer (Calculator Logic)
**Location**: `internal/calculator/`

**Responsibilities**:
- Implementasi business logic calculator
- Expression parsing dan evaluation
- Mode management (standard/scientific)
- Undo/Redo logic
- Validasi operasi

**Key Components**:
- `standard.go` - Standard calculator operations
- `scientific.go` - Scientific calculator operations
- `evaluator.go` - Expression parser dan evaluator
- `memory.go` - Memory operations (MC, MR, M+, M-)

**Dependencies**:
- Math library (math package)
- No external dependencies

### 3. Data Access Layer (Database)
**Location**: `internal/database/`

**Responsibilities**:
- Database connection management
- CRUD operations
- Query execution
- Transaction management
- Data validation

**Key Components**:
- `sqlite.go` - SQLite connection manager
- `models.go` - Data models (Calculation, Setting)
- `migrations.go` - Database schema migrations

**Dependencies**:
- SQLite driver (mattn/go-sqlite3)
- No application layer dependencies

### 4. History Layer
**Location**: `internal/history/`

**Responsibilities**:
- History management logic
- Filtering dan searching
- Export functionality
- History aggregation

**Key Components**:
- `history.go` - History operations
- `filter.go` - Filter logic
- `export.go` - Export to CSV/JSON

**Dependencies**:
- Database layer
- Calculator layer (untuk formatting)

### 5. Config Layer
**Location**: `internal/config/`

**Responsibilities**:
- Configuration management
- Load/save config files
- Validation config
- Default values

**Key Components**:
- `config.go` - Config struct dan operations
- `validator.go` - Config validation

**Dependencies**:
- File system
- YAML parser

### 6. Utility Layer
**Location**: `pkg/utils/`

**Responsibilities**:
- Shared utility functions
- Helper functions
- Common operations

**Key Components**:
- `helpers.go` - Utility functions
- `formatters.go` - Formatting utilities
- `validators.go` - Validation utilities

**Dependencies**:
- None (pure functions)

## Data Flow

### Calculation Flow

```
User Input (CLI)
    │
    ▼
┌─────────────────┐
│  Input Parser   │
│  (ui/cli.go)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Expression      │
│ Validator       │
│ (calculator/    │
│  evaluator.go)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Calculator      │
│ Engine          │
│ (calculator/    │
│  standard.go    │
│  scientific.go) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Result          │
│ Formatter       │
│ (utils/         │
│  formatters.go) │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Database Save   │
│ (database/      │
│  sqlite.go)     │
└────────┬────────┘
         │
         ▼
    Output (CLI)
```

### History Retrieval Flow

```
User Request (CLI)
    │
    ▼
┌─────────────────┐
│  Command Parser │
│  (ui/cli.go)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Filter Builder  │
│ (history/       │
│  filter.go)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Database Query  │
│ (database/      │
│  sqlite.go)     │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Result          │
│ Formatter       │
│ (history/       │
│  export.go)     │
└────────┬────────┘
         │
         ▼
    Output (CLI/File)
```

## Design Patterns

### 1. Repository Pattern
Database layer menggunakan Repository pattern untuk abstraksi data access:

```go
type Database interface {
    InsertCalculation(expression, result, mode string) error
    GetAllCalculations() ([]Calculation, error)
    GetCalculationsByFilter(filter Filter) ([]Calculation, error)
    DeleteCalculation(id int) error
    ClearAllCalculations() error
}
```

### 2. Strategy Pattern
Calculator operations menggunakan Strategy pattern untuk mode switching:

```go
type Calculator interface {
    Evaluate(expression string) (float64, error)
    GetAvailableOperations() []string
}

type StandardCalculator struct { /* ... */ }
type ScientificCalculator struct { /* ... */ }
```

### 3. Builder Pattern
Filter builder untuk history queries:

```go
type FilterBuilder struct {
    filter Filter
}

func (b *FilterBuilder) WithMode(mode string) *FilterBuilder
func (b *FilterBuilder) WithDateRange(start, end time.Time) *FilterBuilder
func (b *FilterBuilder) Build() Filter
```

### 4. Factory Pattern
Calculator factory untuk membuat instance calculator:

```go
func NewCalculator(mode string) Calculator {
    switch mode {
    case "scientific":
        return &ScientificCalculator{}
    default:
        return &StandardCalculator{}
    }
}
```

### 5. Singleton Pattern
Database connection menggunakan Singleton pattern:

```go
var (
    dbInstance *Database
    dbOnce     sync.Once
)

func GetDatabase() *Database {
    dbOnce.Do(func() {
        dbInstance = NewDatabase()
    })
    return dbInstance
}
```

## Data Models

### Calculation Model
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

### Setting Model
```go
type Setting struct {
    Key       string    `json:"key"`
    Value     string    `json:"value"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Config Model
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

## Error Handling Strategy

### Error Types
```go
type CalculatorError struct {
    Code    string
    Message string
    Err     error
}

func (e *CalculatorError) Error() string {
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *CalculatorError) Unwrap() error {
    return e.Err
}
```

### Error Codes
- `CALC_INVALID_EXPRESSION` - Invalid expression syntax
- `CALC_DIVISION_BY_ZERO` - Division by zero
- `CALC_INVALID_FUNCTION` - Invalid function name
- `CALC_DOMAIN_ERROR` - Domain error (e.g., sqrt of negative)
- `DB_CONNECTION_ERROR` - Database connection error
- `DB_QUERY_ERROR` - Database query error
- `CONFIG_ERROR` - Configuration error
- `IO_ERROR` - File I/O error

### Error Handling Flow
```
Error Occurs
    │
    ▼
┌─────────────────┐
│ Error Wrap      │
│ (add context)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Log Error      │
│ (structured)   │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Return to User  │
│ (user-friendly)│
└─────────────────┘
```

## Concurrency Model

### Thread Safety
- Database connection menggunakan connection pool
- Calculator operations are stateless (thread-safe)
- Config menggunakan sync.RWMutex untuk concurrent access

### Goroutine Usage
- Export operations bisa dijalankan di background goroutine
- History queries menggunakan batch processing untuk large datasets

## Security Considerations

### Input Validation
- All user input validated before processing
- SQL injection prevention using parameterized queries
- Path traversal prevention for file operations

### Data Protection
- Database file permissions set correctly
- Config file permissions set correctly
- No sensitive data logged

### Error Messages
- Error messages tidak expose internal details
- Stack traces tidak ditampilkan ke user

## Performance Optimization

### Database Optimization
- Indexes on frequently queried columns (created_at, mode)
- Connection pooling
- Prepared statements for repeated queries
- Batch inserts for bulk operations

### Memory Optimization
- Stream large result sets
- Limit history records in memory
- Use efficient data structures

### Calculation Optimization
- Cache frequently used results
- Pre-compile scientific functions
- Lazy evaluation where possible

## Scalability Considerations

### Database Scaling
- SQLite suitable for single-user application
- Can migrate to PostgreSQL/MySQL for multi-user
- Abstract database layer for easy migration

### Feature Scaling
- Modular design allows easy feature addition
- Plugin architecture for custom functions
- Extensible command system

## Testing Strategy

### Unit Tests
- Test each component in isolation
- Mock external dependencies
- Achieve >80% coverage

### Integration Tests
- Test component interactions
- Use in-memory database
- Test complete workflows

### End-to-End Tests
- Test user scenarios
- Test CLI commands
- Verify data persistence

## Deployment Architecture

### Local Deployment
```
calculator.exe
├── config.yaml
└── data/
    └── calculator.db
```

### Portable Deployment
```
calculator/
├── calculator.exe
├── config/
│   └── config.yaml
└── data/
    └── calculator.db
```

## Monitoring and Logging

### Logging Strategy
- Structured logging using logrus
- Log levels: DEBUG, INFO, WARN, ERROR
- Log rotation for large logs
- Sensitive data not logged

### Metrics
- Calculation count
- Average calculation time
- Database query time
- Error rate

## Technology Rationale

### Go (Golang)
- **Performance**: Fast execution, low memory footprint
- **Concurrency**: Built-in goroutines and channels
- **Type Safety**: Strong typing reduces bugs
- **Cross-platform**: Easy to build for multiple platforms
- **Standard Library**: Rich built-in libraries

### SQLite
- **Zero Configuration**: No database server needed
- **Portable**: Single file database
- **Reliable**: ACID compliant
- **Performance**: Fast for single-user applications
- **SQL Support**: Standard SQL queries

### Cobra CLI
- **Feature-rich**: Built-in help, arguments, flags
- **Flexibility**: Easy to create complex commands
- **Community**: Well-maintained and popular
- **Integration**: Works well with Go ecosystem

## Future Architecture Improvements

### Potential Enhancements
1. **Plugin System**: Allow custom calculator functions
2. **API Layer**: REST API for remote access
3. **Caching Layer**: Redis for performance
4. **Message Queue**: For async operations
5. **Microservices**: Split into separate services

### Migration Path
- Keep current architecture for v1
- Evaluate needs before adding complexity
- Maintain backward compatibility
