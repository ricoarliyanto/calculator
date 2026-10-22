# PRD - Calculator Application

## 1. Overview
Aplikasi Calculator adalah aplikasi kalkulator berbasis CLI (Command Line Interface) yang dibangun dengan bahasa pemrograman Go (Golang) dan menggunakan database SQLite untuk menyimpan history perhitungan. Aplikasi ini mendukung dua mode: Standard Calculator dan Scientific Calculator.

## 2. Objectives
- Menyediakan kalkulator fungsional dengan mode standar dan scientific
- Menyimpan history perhitungan di database SQLite
- Mendukung perpindahan antar mode dengan mudah
- Menyediakan fitur tambahan seperti undo, redo, dan export history

## 3. Target Users
- Pengguna yang membutuhkan kalkulator dengan fitur dasar dan scientific
- Pengguna yang ingin menyimpan history perhitungan
- Pengguna yang prefer CLI application

## 4. Functional Requirements

### 4.1 Standard Calculator Features
- **Basic Operations**: Addition (+), Subtraction (-), Multiplication (*), Division (/)
- **Parentheses Support**: Mengelompokkan operasi dengan tanda kurung
- **Decimal Numbers**: Mendukung angka desimal
- **Clear Function**: Reset input dan hasil
- **Delete Function**: Menghapus karakter terakhir input
- **Percentage Operation**: Menghitung persentase

### 4.2 Scientific Calculator Features
- **Trigonometric Functions**: sin, cos, tan (derajat dan radian)
- **Inverse Trigonometric**: asin, acos, atan
- **Logarithmic Functions**: log (base 10), ln (natural log)
- **Exponential Functions**: exp, power (^), square root (√)
- **Factorial**: n!
- **Constants**: π (pi), e (Euler's number)
- **Memory Functions**: MC, MR, M+, M-
- **Modulo Operation**: %

### 4.3 History Features
- **Save History**: Setiap perhitungan disimpan ke database
- **View History**: Menampilkan semua history perhitungan
- **Filter History**: Filter berdasarkan tanggal, operator, atau mode
- **Delete History**: Menghapus entry history tertentu
- **Clear All History**: Menghapus seluruh history
- **Export History**: Export history ke file CSV/JSON

### 4.4 Mode Switching
- **Switch Mode**: Perpindahan antara Standard dan Scientific mode
- **Persistent Mode**: Mode terakhir yang digunakan disimpan
- **Mode Indicator**: Tampilan indikator mode aktif

### 4.5 Additional Features
- **Undo/Redo**: Kembalikan ke perhitungan sebelumnya
- **Keyboard Shortcuts**: Shortcut untuk operasi umum
- **Config File**: Konfigurasi pengaturan (decimal places, angle unit, dll)
- **Help Command**: Menampilkan bantuan dan dokumentasi

## 5. Non-Functional Requirements

### 5.1 Performance
- Response time < 100ms untuk operasi kalkulasi
- Database query < 50ms untuk history operations

### 5.2 Reliability
- Error handling yang baik untuk invalid input
- Data integrity untuk database operations
- Graceful error recovery

### 5.3 Usability
- CLI yang intuitive dan user-friendly
- Clear error messages
- Comprehensive help documentation

### 5.4 Maintainability
- Clean code architecture
- Well-documented code
- Modular design
- Test coverage > 80%

## 6. Technical Requirements

### 6.1 Technology Stack
- **Language**: Go (Golang) 1.21+
- **Database**: SQLite 3
- **Database Driver**: mattn/go-sqlite3
- **CLI Framework**: cobra or urfave/cli
- **Testing**: Go testing framework
- **Logging**: logrus or zap

### 6.2 Database Schema

#### Table: calculations
```sql
CREATE TABLE calculations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    expression TEXT NOT NULL,
    result TEXT NOT NULL,
    mode TEXT NOT NULL, -- 'standard' or 'scientific'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Table: settings
```sql
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 6.3 Project Structure
```
calculator/
├── cmd/
│   └── calculator/
│       └── main.go
├── internal/
│   ├── calculator/
│   │   ├── standard.go
│   │   ├── scientific.go
│   │   └── evaluator.go
│   ├── database/
│   │   ├── sqlite.go
│   │   └── models.go
│   ├── history/
│   │   └── history.go
│   ├── config/
│   │   └── config.go
│   └── ui/
│       └── cli.go
├── pkg/
│   └── utils/
│       └── helpers.go
├── data/
│   └── calculator.db
├── go.mod
├── go.sum
├── PRD.md
└── PLANNING.md
```

## 7. User Stories

### US1: Basic Calculation
Sebagai user, saya ingin melakukan perhitungan dasar (tambah, kurang, kali, bagi) untuk mendapatkan hasil yang akurat.

### US2: Scientific Calculation
Sebagai user, saya ingin melakukan perhitungan scientific (trigonometri, logaritma, dll) untuk keperluan akademik atau profesional.

### US3: History Viewing
Sebagai user, saya ingin melihat history perhitungan saya untuk melihat kembali hasil perhitungan sebelumnya.

### US4: Mode Switching
Sebagai user, saya ingin berganti antara mode standar dan scientific sesuai kebutuhan.

### US5: Export History
Sebagai user, saya ingin export history perhitungan ke file untuk backup atau analisis.

## 8. Acceptance Criteria

### AC1: Standard Calculator
- [ ] Semua operasi dasar berfungsi dengan benar
- [ ] Input dengan desimal ditangani dengan benar
- [ ] Parentheses support berfungsi
- [ ] Error handling untuk division by zero

### AC2: Scientific Calculator
- [ ] Semua fungsi trigonometri berfungsi
- [ ] Logaritma dan eksponen berfungsi
- [ ] Konstanta (π, e) tersedia
- [ ] Factorial berfungsi untuk bilangan bulat positif

### AC3: History
- [ ] Setiap perhitungan disimpan ke database
- [ ] History dapat ditampilkan
- [ ] History dapat difilter
- [ ] History dapat dihapus
- [ ] History dapat di-export

### AC4: Mode Switching
- [ ] Perpindahan mode berfungsi
- [ ] Mode terakhir disimpan
- [ ] Indikator mode ditampilkan

## 9. Future Enhancements
- GUI version (desktop application)
- Web interface
- Mobile app
- Cloud sync for history
- Custom functions
- Graph plotting
- Unit converter
