# Calculator Application

Aplikasi Calculator berbasis CLI (Command Line Interface) yang dibangun dengan Golang dan SQLite. Mendukung mode Standard dan Scientific dengan fitur history lengkap.

## Features

### Standard Calculator
- Operasi dasar: tambah (+), kurang (-), kali (*), bagi (/)
- Support parentheses ( )
- Angka desimal
- Operasi persentase (%)
- Clear dan delete function

### Scientific Calculator
- Fungsi trigonometri: sin, cos, tan (derajat dan radian)
- Inverse trigonometri: asin, acos, atan
- Logaritma: log (base 10), ln (natural log)
- Eksponen: exp, power (^), square root (√)
- Factorial (n!)
- Konstanta: π (pi), e (Euler's number)
- Memory functions: MC, MR, M+, M-
- Modulo operation (%)

### History Management
- Simpan setiap perhitungan ke database SQLite
- View history perhitungan
- Filter history berdasarkan tanggal, operator, atau mode
- Delete history entry tertentu
- Clear all history
- Export history ke CSV/JSON

### Additional Features
- Undo/Redo perhitungan
- Mode switching (Standard ↔ Scientific)
- Konfigurasi pengaturan
- Keyboard shortcuts
- Interactive REPL mode

## Installation

### Prerequisites
- Go 1.21 atau higher
- GCC (for CGO/SQLite on Windows)

### Build from Source

```bash
# Clone repository
git clone https://github.com/username/calculator.git
cd calculator

# Install dependencies
go mod download

# Build binary
go build -o calculator.exe ./cmd/calculator

# Run
./calculator.exe
```

### Install on Windows

```bash
# Build for Windows
go build -o calculator.exe ./cmd/calculator

# Add to PATH (optional)
# Copy calculator.exe to directory in PATH
```

### Build for macOS
```bash
# macOS requires a macOS machine or cross-compilation toolchain
GOSO=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o calculator-darwin ./cmd/calculator
```

### Pre-built Binaries

Pre-built binaries are available for Linux and Windows:

| Platform | Architecture | Binary |
|----------|-------------|--------|
| Linux    | amd64       | `calculator-linux-amd64` |
| Windows  | amd64       | `calculator-windows-amd64.exe` |
| macOS    | amd64       | Build from source (see above) |

## Troubleshooting

### Common Issues

#### "cannot find package" or build errors

Ensure you have the correct module path and dependencies:
```bash
go mod tidy
go mod download
```

#### CGO/SQLite errors on Linux

Install GCC and SQLite development headers:
```bash
# Ubuntu/Debian
sudo apt-get install gcc libsqlite3-dev

# Fedora/RHEL
sudo dnf install gcc sqlite-devel
```

#### Windows: "gcc: error: unrecognized command-line option '-mthreads'"

This error occurs when cross-compiling for Windows on Linux. Install the MinGW cross-compiler:
```bash
sudo apt-get install gcc-mingw-w64-x86-64
```

Then build with:
```bash
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o calculator.exe ./cmd/calculator
```

#### Database connection errors

If you see "database connection not established":
1. Ensure the data directory is writable
2. Check if SQLite library is installed
3. Try clearing the database file (deletes history):
   ```bash
   rm ~/.config/calculator/calculator.db
   ```

#### Settings not saving

The config file is saved to:
- Linux: `~/.config/calculator/config.yaml`
- Windows: `%APPDATA%\calculator\config.yaml`

Ensure the directory exists and is writable. Use `calc config path` to check the location.

#### "command not found" after installation

Ensure the binary is in your PATH:
```bash
# Linux/macOS
sudo cp calculator /usr/local/bin/

# Windows (PowerShell as Administrator)
move calculator.exe C:\Windows\System32\
```

## Usage

### Interactive Mode

```bash
# Start calculator (default: standard mode)
calculator

# Start in standard mode
calculator --mode standard

# Start in scientific mode
calculator --mode scientific
```

### Commands

#### Calculator Commands (Interactive Mode)
```
> 2 + 2                    # Basic calculation
4
> (3 + 4) * 5              # Parentheses
35
> sin(90)                  # Scientific function (deg mode)
1
> mode scientific          # Switch to scientific mode
> undo                     # Undo last calculation
> redo                     # Redo undone calculation
> clear                    # Clear current input
> exit                     # Exit calculator
```

#### History Commands
```bash
# View all history
calculator history

# Filter by mode
calculator history --mode scientific

# Filter by date
calculator history --from "2024-01-01" --to "2024-12-31"

# Export to CSV
calculator history --export csv --output history.csv

# Export to JSON
calculator history --export json --output history.json

# Clear all history
calculator history --clear

# Delete specific entry
calculator history --delete 5
```

#### Settings Commands
```bash
# View settings
calculator settings

# Set default mode
calculator settings set mode scientific

# Set angle unit (deg/rad)
calculator settings set angle-unit deg

# Set decimal places
calculator settings set decimal-places 4

# Reset to defaults
calculator settings reset
```

#### Config Commands
```bash
# View config file location
calculator config path

# Edit config
calculator config edit

# Reset config
calculator config reset
```

### Configuration

Config file disimpan di:
- Windows: `%APPDATA%/calculator/config.yaml`
- Linux/macOS: `~/.config/calculator/config.yaml`

Example config.yaml:
```yaml
calculator:
  mode: standard
  angle-unit: deg
  decimal-places: 4
  history-limit: 1000
  auto-save: true
  show-timestamps: true
```

## Project Structure

```
calculator/
├── cmd/
│   └── calculator/
│       └── main.go           # Entry point
├── internal/
│   ├── calculator/
│   │   ├── standard.go       # Standard calculator operations
│   │   ├── scientific.go     # Scientific calculator operations
│   │   ├── evaluator.go      # Expression parser (Shunting Yard)
│   │   ├── calculator_test.go # Unit + edge case tests
│   │   └── benchmark_test.go  # Performance benchmarks
│   ├── database/
│   │   ├── sqlite.go         # Database connection manager
│   │   ├── models.go         # Data models
│   │   ├── migrations.go     # Migration manager
│   │   ├── calculations.go   # Calculation repository
│   │   ├── settings.go       # Settings repository
│   │   └── database_test.go  # Unit tests
│   ├── history/
│   │   ├── history.go        # History management
│   │   └── history_test.go   # Unit tests
│   ├── config/
│   │   ├── config.go         # Configuration management
│   │   └── config_test.go    # Unit tests
│   └── ui/
│       └── cli.go            # CLI interface (Cobra)
├── pkg/
│   └── utils/
│       ├── helpers.go        # Utility functions
│       └── helpers_test.go   # Tests
├── data/
│   └── calculator.db         # SQLite database
├── md-files/
│   ├── PRD.md                # Product Requirements
│   ├── PLANNING.md           # Development Planning
│   ├── TDD.md                # Test Strategy
│   ├── ARCHITECTURE.md       # Architecture Documentation
│   ├── API.md                # API Documentation
│   └── README.md             # Original README
├── go.mod
├── go.sum
├── coverage.out            # Generated coverage report
└── .gitignore
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/calculator

# Run benchmarks
go test -bench=. ./...

# Run tests with verbose output
go test -v ./...
```

### Test Coverage

Phase 5 test coverage results:

| Package      | Coverage |
|-------------|----------|
| calculator  | 89.8%    |
| config      | 82.3%    |
| database    | 80.5%    |
| history     | 82.5%    |
| utils       | N/A      |

Overall test suite covers **all core packages with >80% coverage**.

### Code Quality

```bash
# Run go vet
go vet ./...

# Format code
go fmt ./...

# Run golangci-lint (if installed)
golangci-lint run
```

### Build

```bash
# Build for current platform
go build ./cmd/calculator

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o calculator.exe ./cmd/calculator

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o calculator ./cmd/calculator

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o calculator ./cmd/calculator
```

## Documentation

- [PRD.md](PRD.md) - Product Requirements Document
- [PLANNING.md](PLANNING.md) - Development Planning
- [TDD.md](TDD.md) - Test Driven Development Strategy
- [ARCHITECTURE.md](ARCHITECTURE.md) - Architecture Documentation
- [API.md](API.md) - API Documentation

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see LICENSE file for details

## Roadmap

- [ ] GUI version (desktop application)
- [ ] Web interface
- [ ] Mobile app
- [ ] Cloud sync for history
- [ ] Custom functions
- [ ] Graph plotting
- [ ] Unit converter

## Support

For issues and questions:
- Open an issue on GitHub
- Check documentation in md-files/
- Contact: rico@example.com

## Acknowledgments

Built with:
- Go (Golang)
- SQLite
- Cobra CLI framework
