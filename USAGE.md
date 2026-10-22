# Calculator - Usage Guide

## Quick Start

```bash
# Build the calculator
go build -o calculator ./cmd/calculator

# Run in interactive mode (default)
./calculator

# Evaluate a single expression (non-interactive)
./calculator "2 + 2"
```

## Command Line Options

```bash
# View all available commands
./calculator --help

# View version
./calculator version

# Run with specific mode
./calculator --mode scientific

# Show debug logs
./calculator --debug
```

## Interactive Mode

Start the calculator in interactive REPL mode:

```bash
./calculator
```

### Available Commands in Interactive Mode

#### Calculator Operations

| Command | Description | Example |
|---------|-------------|---------|
| `expression` | Evaluate any math expression | `> (10 + 5) * 2` |
| `mode` | Switch between standard/scientific | `> mode scientific` |
| `undo` | Undo the last calculation | `> undo` |
| `redo` | Redo an undone calculation | `> redo` |
| `clear` | Clear current entry | `> clear` |
| `help` | Show available commands | `> help` |
| `exit` | Exit the calculator | `> exit` |

#### Memory Operations

| Command | Description | Example |
|---------|-------------|---------|
| `mc` or `memory-clear` | Clear memory | `> mc` |
| `mr` or `memory-recall` | Recall memory value | `> mr` |
| `m+` or `memory-add` | Add to memory | `> m+ 5` |
| `m-` or `memory-subtract` | Subtract from memory | `> m- 3` |

#### History Commands

| Command | Description | Example |
|---------|-------------|---------|
| `history` | View calculation history | `> history` |
| `history --mode <mode>` | Filter history by mode | `> history --mode scientific` |
| `history --search <query>` | Search history | `> history --search "sin"` |
| `history --clear` | Clear all history | `> history --clear` |
| `history --delete <id>` | Delete specific entry | `> history --delete 5` |
| `history --export <format>` | Export to file | `> history --export json` |

### Expression Syntax

#### Standard Mode
- Basic arithmetic: `+`, `-`, `*`, `/`
- Parentheses: `(`, `)`
- Decimal numbers: `3.14`
- Percentage: `100 + 10%` → `110`
- Modulo: `10 % 3` → `1`

```bash
> 2 + 2
4
> (3 + 4) * 5
35
> 100 + 10%
110
```

#### Scientific Mode
All standard operations plus:

| Function | Description | Example |
|----------|-------------|---------|
| `sin(x)` | Sine | `> sin(90, deg)` → `1` |
| `cos(x)` | Cosine | `> cos(0, deg)` → `1` |
| `tan(x)` | Tangent | `> tan(45, deg)` → `1` |
| `asin(x)` | Arc sine | `> asin(1, deg)` → `90` |
| `acos(x)` | Arc cosine | `> acos(1, deg)` → `0` |
| `atan(x)` | Arc tangent | `> atan(1, deg)` → `45` |
| `log(x)` | Log base 10 | `> log(100)` → `2` |
| `ln(x)` | Natural log | `> ln(e)` → `1` |
| `exp(x)` | e^x | `> exp(1)` → `2.718...` |
| `sqrt(x)` | Square root | `> sqrt(16)` → `4` |
| `n!` | Factorial | `> 5!` → `120` |
| `x^y` | Power | `> 2^10` |
| `%` | Modulo | `> 10 % 3` → `1` |
| `pi` | Pi constant | `> pi` → `3.14159...` |
| `e` | Euler's number | `> e` → `2.71828...` |

### Angle Units

Scientific functions accept an optional angle unit argument:
```bash
> sin(90, deg)      # sin(90 degrees) = 1
> cos(0, deg)       # cos(0 degrees) = 1
> sin(1.5708)       # sin(1.5708 radians) ≈ 1
> sin(90, rad)      # sin(90 radians) ≈ 0.894
```

### Complex Expressions

```bash
> 2 * (3 + (4 * (5 - 1)))    # Nested parentheses
38
> sin(90, deg) + cos(0, deg) * 2    # Mixed operations
3
> 2^10 + 10!    # Power and factorial
3638811
```

## CLI Commands

```bash
# Calculate an expression directly (non-interactive)
./calculator "sin(90, deg) + cos(0, deg)"

# View history
./calculator history

# View settings
./calculator settings

# Manage settings
./calculator settings set mode scientific
./calculator settings set angle-unit deg
./calculator settings set decimal-places 4
./calculator settings reset

# Manage config
./calculator config path
./calculator config reset

# Export history
./calculator history --export json --output my_history.json
./calculator history --export csv --output my_history.csv

# Clear history
./calculator history --clear

# View version
./calculator version
```

## Configuration

The calculator stores configuration in a YAML file:

### Default Locations
- **Linux/macOS**: `~/.config/calculator/config.yaml`
- **Windows**: `%APPDATA%\calculator\config.yaml`

### Config File Format

```yaml
mode: standard
angle-unit: rad
decimal-places: 10
history-limit: 1000
auto-save: true
show-timestamps: true
theme: light
```

### Settings Reference

| Setting | Values | Default | Description |
|---------|--------|---------|-------------|
| `mode` | `standard`, `scientific` | `standard` | Default calculator mode |
| `angle-unit` | `rad`, `deg` | `rad` | Angle unit for trig functions |
| `decimal-places` | 0-20 | `10` | Number of decimal places for results |
| `history-limit` | 0-100000 | `1000` | Maximum history entries to keep |
| `auto-save` | `true`, `false` | `true` | Auto-save calculations to history |
| `show-timestamps` | `true`, `false` | `true` | Show timestamps in history |
| `theme` | `light`, `dark` | `light` | UI color theme |

## Database

The calculator uses SQLite for history storage:

- **Location**: Same as config directory
- **File**: `calculator.db`
- **Tables**: `calculations`, `settings`, `migrations`

### History Filtering

```bash
./calculator history --mode scientific
./calculator history --search "sin"
./calculator history --from "2024-01-01" --to "2024-12-31"
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Run benchmarks
go test -bench=. ./internal/calculator/ -benchmem
```

### Building from Source

```bash
# Prerequisites: Go 1.21+, GCC (for CGO/SQLite)

# Build
CGO_ENABLED=1 go build -o calculator ./cmd/calculator

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 go build -o calculator.exe ./cmd/calculator

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o calculator ./cmd/calculator

# Cross-compile for macOS
GOOS=darwin GOARCH=amd64 CGO_ENABLED=1 go build -o calculator ./cmd/calculator
```
