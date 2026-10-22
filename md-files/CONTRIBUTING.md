# Contributing to Calculator Application

Terima kasih atas ketertarikan Anda untuk berkontribusi pada proyek Calculator! Dokumen ini menjelaskan cara berkontribusi.

## How to Contribute

### Reporting Bugs

Sebelum melaporkan bug:
1. Cek existing issues di GitHub
2. Pastikan bug belum dilaporkan
3. Gunakan template issue yang disediakan

Saat melaporkan bug:
- Gunakan judul yang jelas dan deskriptif
- Jelaskan langkah untuk reproduce bug
- Sertakan error message dan stack trace
- Sertakan environment details (OS, Go version)
- Sertakan screenshot jika relevan

### Suggesting Enhancements

Sebelum mengusulkan enhancement:
1. Cek existing issues dan PRs
2. Pastikan usulan belum ada
3. Pertimbangkan apakah enhancement sesuai dengan scope proyek

Saat mengusulkan enhancement:
- Jelaskan masalah yang ingin diselesaikan
- Jelaskan solusi yang diusulkan
- Pertimbangkan alternatif lain
- Sertakan contoh penggunaan

### Pull Requests

#### Before Submitting PR

1. Fork repository
2. Buat branch untuk fitur/fix Anda
3. Ikuti coding standards
4. Tulis tests untuk changes Anda
5. Pastikan semua tests pass
6. Update documentation jika perlu
7. Commit dengan pesan yang jelas

#### PR Guidelines

- Satu PR per fitur/fix
- Branch naming:
  - `feature/feature-name` untuk fitur baru
  - `fix/bug-description` untuk bug fix
  - `docs/update-description` untuk dokumentasi
  - `refactor/component-name` untuk refactoring
- Commit message format:
  ```
  type(scope): subject

  body

  footer
  ```
  Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

  Example:
  ```
  feat(calculator): add power function for scientific mode

  Add power function to support exponentiation operations.
  Includes unit tests and documentation updates.

  Closes #123
  ```

#### Code Review Process

1. PR akan direview oleh maintainer
2. Reviewer mungkin meminta changes
3. Lakukan changes yang diminta
4. Push changes ke branch yang sama
5. Tunggu approval dan merge

## Development Setup

### Prerequisites
- Go 1.21+
- GCC (for CGO/SQLite on Windows)
- Git

### Setup Development Environment

```bash
# Fork dan clone repository
git clone https://github.com/your-username/calculator.git
cd calculator

# Add upstream remote
git remote add upstream https://github.com/original-owner/calculator.git

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build ./cmd/calculator
```

### Keeping Your Fork Updated

```bash
# Fetch upstream changes
git fetch upstream

# Merge upstream changes to your main branch
git checkout main
git merge upstream/main

# Or rebase
git rebase upstream/main
```

## Coding Standards

### Go Code Style

Ikuti [Effective Go](https://golang.org/doc/effective_go) dan [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

#### Naming Conventions
- Use `camelCase` untuk variables dan functions
- Use `PascalCase` untuk exported types dan functions
- Use `UPPER_SNAKE_CASE` untuk constants
- Use descriptive names, avoid abbreviations

#### Example
```go
// Good
func CalculateAverage(numbers []float64) float64 {
    // ...
}

// Bad
func calcAvg(nums []float64) float64 {
    // ...
}
```

### Code Organization

#### Package Structure
- Internal packages untuk code yang tidak di-export
- Public packages untuk reusable code
- Separate concerns: calculator, database, history, config, ui

#### File Organization
- One type per file (kecuali types yang terkait erat)
- Related functions dalam file yang sama
- Test files: `filename_test.go`

### Comments

#### Package Comments
```go
// Package calculator provides calculator operations
// for standard and scientific modes.
package calculator
```

#### Function Comments
```go
// Add performs addition of two numbers.
// It returns the sum of a and b.
func Add(a, b float64) float64 {
    return a + b
}
```

#### Complex Logic Comments
```go
// Calculate factorial using recursion.
// Base case: factorial of 0 or 1 is 1.
// Recursive case: n * factorial(n-1)
func Factorial(n int) (int, error) {
    if n < 0 {
        return 0, errors.New("factorial of negative number")
    }
    if n <= 1 {
        return 1, nil
    }
    return n * Factorial(n-1)
}
```

### Error Handling

#### Always Handle Errors
```go
// Good
result, err := Calculate(expression)
if err != nil {
    log.Printf("Calculation failed: %v", err)
    return err
}

// Bad
result, _ := Calculate(expression) // Don't ignore errors!
```

#### Use Custom Error Types
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

### Testing

#### Write Tests First (TDD)
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     float64
        expected float64
    }{
        {"positive numbers", 5, 3, 8},
        {"negative numbers", -5, -3, -8},
        {"zero", 5, 0, 5},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%f, %f) = %f; want %f", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

#### Test Coverage
- Target >80% coverage
- Test both happy path and error cases
- Test edge cases

#### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestAdd ./internal/calculator

# Run benchmarks
go test -bench=. ./...
```

### Documentation

#### Update Documentation
- Update README.md untuk user-facing changes
- Update API.md untuk API changes
- Update ARCHITECTURE.md untuk architectural changes
- Add comments untuk complex logic

#### Example Documentation
```go
// Power calculates base raised to the power of exponent.
//
// Parameters:
//   - base: The base number
//   - exponent: The exponent
//
// Returns:
//   - float64: The result of base^exponent
//   - error: Error if the operation is invalid (e.g., negative base with non-integer exponent)
//
// Example:
//   result, err := Power(2, 3) // returns 8, nil
func Power(base, exponent float64) (float64, error) {
    // ...
}
```

## Git Workflow

### Branch Strategy
- `main` - Production code
- `develop` - Development code (optional)
- `feature/*` - Feature branches
- `fix/*` - Bug fix branches
- `hotfix/*` - Urgent fixes for production

### Commit Workflow

```bash
# Create feature branch
git checkout -b feature/add-power-function

# Make changes
# ...

# Stage changes
git add .

# Commit
git commit -m "feat(calculator): add power function for scientific mode"

# Push to fork
git push origin feature/add-power-function

# Create PR
```

### Commit Message Format

```
type(scope): subject

body

footer
```

#### Types
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc)
- `refactor`: Code refactoring
- `test`: Test changes
- `chore`: Build process or auxiliary tool changes

#### Examples
```
feat(database): add export to JSON functionality

Add JSON export for calculation history.
Implements #45

Closes #45
```

```
fix(calculator): handle division by zero gracefully

Return error instead of panic when dividing by zero.
Add test for error case.

Fixes #23
```

```
docs(readme): update installation instructions

Add Windows-specific installation steps.
Update Go version requirement.
```

## Code Review Checklist

Sebelum submitting PR, pastikan:
- [ ] Code follows project style guidelines
- [ ] All tests pass (`go test ./...`)
- [ ] Test coverage >80%
- [ ] Code passes linters (`go vet ./...`, `golangci-lint run`)
- [ ] Documentation updated
- [ ] Comments added untuk complex logic
- [ ] No hardcoded values (use constants)
- [ ] Error handling implemented
- [ ] No debug code or console logs
- [ ] Function and variable names are descriptive
- [ ] Functions are small and focused
- [ ] No code duplication

## Project-Specific Guidelines

### Calculator Logic
- Maintain separation between standard and scientific modes
- Use expression parser consistently
- Handle all edge cases (division by zero, domain errors, etc)
- Support undo/redo operations

### Database Operations
- Always use parameterized queries
- Handle database errors gracefully
- Close connections properly
- Use transactions for multiple operations

### CLI Interface
- Provide clear error messages
- Support both interactive and command modes
- Validate user input
- Format output consistently

### Testing
- Write unit tests untuk setiap function
- Write integration tests untuk workflows
- Test edge cases thoroughly
- Mock external dependencies

## Getting Help

Jika Anda butuh bantuan:
- Check documentation di `md-files/`
- Open issue di GitHub
- Contact maintainer
- Join community discussion (if available)

## License

Dengan berkontribusi, Anda agree bahwa kontribusi Anda akan dilisensikan di bawah MIT License yang sama dengan proyek.

## Recognition

Kontributor akan diakui di:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Terima kasih atas kontribusi Anda! 🎉
