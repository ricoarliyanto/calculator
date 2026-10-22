package calculator

import (
	"math"
	"testing"
)

// =============================================================================
// Standard Calculator Tests
// =============================================================================

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
		{"zero with positive", 5, 0, 5},
		{"zero with negative", -5, 0, -5},
		{"large numbers", 1e10, 1e10, 2e10},
		{"small numbers", 1e-10, 1e-10, 2e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Add(tt.a, tt.b)
			if err != nil {
				t.Errorf("Add() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 5, 3, 2},
		{"negative numbers", -5, -3, -2},
		{"mixed signs", 5, -3, 8},
		{"decimals", 5.5, 3.3, 2.2},
		{"zero with positive", 5, 0, 5},
		{"zero with negative", -5, 0, -5},
		{"result negative", 3, 5, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Subtract(tt.a, tt.b)
			if err != nil {
				t.Errorf("Subtract() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 5, 3, 15},
		{"negative numbers", -5, -3, 15},
		{"mixed signs", 5, -3, -15},
		{"decimals", 5.5, 2, 11.0},
		{"zero with positive", 5, 0, 0},
		{"zero with negative", -5, 0, 0},
		{"by one", 5, 1, 5},
		{"fractional", 0.5, 0.5, 0.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Multiply(tt.a, tt.b)
			if err != nil {
				t.Errorf("Multiply() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
		wantErr  bool
	}{
		{"positive numbers", 6, 3, 2, false},
		{"negative numbers", -6, -3, 2, false},
		{"mixed signs", 6, -3, -2, false},
		{"decimals", 6.6, 2, 3.3, false},
		{"zero numerator", 0, 5, 0, false},
		{"one", 5, 1, 5, false},
		{"fractional", 1, 2, 0.5, false},
		{"division by zero", 5, 0, 0, true},
		{"negative division by zero", -5, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Divide(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Divide(%v, %v) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("Divide(%v, %v) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		percent  float64
		expected float64
	}{
		{"50 percent of 100", 100, 50, 50},
		{"10 percent of 200", 200, 10, 20},
		{"0 percent of 100", 100, 0, 0},
		{"100 percent of 50", 50, 100, 50},
		{"25 percent of 80", 80, 25, 20},
		{"negative base", -100, 50, -50},
		{"negative percentage", 100, -50, -50},
		{"decimal values", 50.5, 50, 25.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Percentage(tt.value, tt.percent)
			if err != nil {
				t.Errorf("Percentage() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Percentage(%v, %v) = %v, want %v", tt.value, tt.percent, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Scientific Calculator Tests
// =============================================================================

func TestSin(t *testing.T) {
	tests := []struct {
		name     string
		angle    float64
		unit     string
		expected float64
	}{
		{"0 radians", 0, "rad", 0},
		{"PI/2 radians", math.Pi / 2, "rad", 1},
		{"PI radians", math.Pi, "rad", 0},
		{"3*PI/2 radians", 3 * math.Pi / 2, "rad", -1},
		{"0 degrees", 0, "deg", 0},
		{"90 degrees", 90, "deg", 1},
		{"180 degrees", 180, "deg", 0},
		{"270 degrees", 270, "deg", -1},
		{"360 degrees", 360, "deg", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Sin(tt.angle, tt.unit)
			if err != nil {
				t.Errorf("Sin() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Sin(%v, %v) = %v, want %v", tt.angle, tt.unit, result, tt.expected)
			}
		})
	}
}

func TestCos(t *testing.T) {
	tests := []struct {
		name     string
		angle    float64
		unit     string
		expected float64
	}{
		{"0 radians", 0, "rad", 1},
		{"PI/2 radians", math.Pi / 2, "rad", 0},
		{"PI radians", math.Pi, "rad", -1},
		{"0 degrees", 0, "deg", 1},
		{"90 degrees", 90, "deg", 0},
		{"180 degrees", 180, "deg", -1},
		{"360 degrees", 360, "deg", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Cos(tt.angle, tt.unit)
			if err != nil {
				t.Errorf("Cos() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Cos(%v, %v) = %v, want %v", tt.angle, tt.unit, result, tt.expected)
			}
		})
	}
}

func TestTan(t *testing.T) {
	tests := []struct {
		name     string
		angle    float64
		unit     string
		expected float64
	}{
		{"0 radians", 0, "rad", 0},
		{"PI/4 radians", math.Pi / 4, "rad", 1},
		{"0 degrees", 0, "deg", 0},
		{"45 degrees", 45, "deg", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Tan(tt.angle, tt.unit)
			if err != nil {
				t.Errorf("Tan() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Tan(%v, %v) = %v, want %v", tt.angle, tt.unit, result, tt.expected)
			}
		})
	}
}

func TestASin(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		unit     string
		expected float64
		wantErr  bool
	}{
		{"0 radians", 0, "rad", 0, false},
		{"1 radians", 1, "rad", math.Pi / 2, false},
		{"0.5 radians", 0.5, "rad", math.Pi / 6, false},
		{"1 degrees", 1, "deg", 90, false},
		{"out of range > 1", 2, "rad", 0, true},
		{"out of range < -1", -2, "rad", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ASin(tt.value, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ASin(%v, %v) error = %v, wantErr %v", tt.value, tt.unit, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Convert to same unit for comparison
				expected := tt.expected
				if tt.unit == "deg" {
					expected = expected * 180 / math.Pi
				}
				if !almostEqual(result, expected) {
					t.Errorf("ASin(%v, %v) = %v, want %v", tt.value, tt.unit, result, expected)
				}
			}
		})
	}
}

func TestACos(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		unit     string
		expected float64
		wantErr  bool
	}{
		{"1 radians", 1, "rad", 0, false},
		{"0.5 radians", 0.5, "rad", math.Pi / 3, false},
		{"0 radians", 0, "rad", math.Pi / 2, false},
		{"out of range > 1", 2, "rad", 0, true},
		{"out of range < -1", -2, "rad", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ACos(tt.value, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ACos(%v, %v) error = %v, wantErr %v", tt.value, tt.unit, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				expected := tt.expected
				if tt.unit == "deg" {
					expected = expected * 180 / math.Pi
				}
				if !almostEqual(result, expected) {
					t.Errorf("ACos(%v, %v) = %v, want %v", tt.value, tt.unit, result, expected)
				}
			}
		})
	}
}

func TestATan(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		unit     string
		expected float64
	}{
		{"0", 0, "rad", 0},
		{"1 radians", 1, "rad", math.Pi / 4},
		{"0 degrees", 0, "deg", 0},
		{"1 degrees", 1, "deg", 45},
		{"negative", -1, "rad", -math.Pi / 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ATan(tt.value, tt.unit)
			if err != nil {
				t.Errorf("ATan() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				expected := tt.expected
				if tt.unit == "deg" {
					expected = expected * 180 / math.Pi
				}
				if !almostEqual(result, expected) {
					t.Errorf("ATan(%v, %v) = %v, want %v", tt.value, tt.unit, result, expected)
				}
			}
		})
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
		wantErr  bool
	}{
		{"1", 1, 0, false},
		{"10", 10, 1, false},
		{"100", 100, 2, false},
		{"0.1", 0.1, -1, false},
		{"negative", -1, 0, true},
		{"zero", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Log(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Log(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("Log(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestLn(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
		wantErr  bool
	}{
		{"1", 1, 0, false},
		{"e", math.E, 1, false},
		{"e squared", math.E * math.E, 2, false},
		{"negative", -1, 0, true},
		{"zero", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Ln(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ln(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("Ln(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestExp(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"0", 0, 1},
		{"1", 1, math.E},
		{"2", 2, math.E * math.E},
		{"negative", -1, 1 / math.E},
		{"0.5", 0.5, math.Sqrt(math.E)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Exp(tt.value)
			if err != nil {
				t.Errorf("Exp() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Exp(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
		wantErr  bool
	}{
		{"0", 0, 0, false},
		{"1", 1, 1, false},
		{"4", 4, 2, false},
		{"9", 9, 3, false},
		{"2", 2, math.Sqrt(2), false},
		{"negative", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Sqrt(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sqrt(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("Sqrt(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected float64
		wantErr  bool
	}{
		{"0", 0, 1, false},
		{"1", 1, 1, false},
		{"2", 2, 2, false},
		{"3", 3, 6, false},
		{"5", 5, 120, false},
		{"10", 10, 3628800, false},
		{"negative", -1, 0, true},
		{"large number", 21, 0, true}, // Should error for very large numbers
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Factorial(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Factorial(%v) error = %v, wantErr %v", tt.value, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("Factorial(%v) = %v, want %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestPower(t *testing.T) {
	tests := []struct {
		name     string
		base     float64
		exponent float64
		expected float64
	}{
		{"2^3", 2, 3, 8},
		{"2^-1", 2, -1, 0.5},
		{"3^2", 3, 2, 9},
		{"5^0", 5, 0, 1},
		{"0^5", 0, 5, 0},
		{"10^1", 10, 1, 10},
		{"2^0.5", 2, 0.5, math.Sqrt(2)},
		{"decimal base", 2.5, 2, 6.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Power(tt.base, tt.exponent)
			if err != nil {
				t.Errorf("Power() error = %v, wantErr nil", err)
				return
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Power(%v, %v) = %v, want %v", tt.base, tt.exponent, result, tt.expected)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"Pi", Pi},
		{"E", E},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify constants are defined and have reasonable values
			if tt.value <= 0 {
				t.Errorf("%s should be positive, got %v", tt.name, tt.value)
			}
		})
	}
}

// =============================================================================
// Expression Parser Tests
// =============================================================================

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name      string
		expression string
		expected  float64
		wantErr   bool
	}{
		// Basic operations
		{"simple addition", "2 + 2", 4, false},
		{"simple subtraction", "5 - 3", 2, false},
		{"simple multiplication", "2 * 3", 6, false},
		{"simple division", "6 / 2", 3, false},
		
		// Operator precedence
		{"multiplication first", "2 + 3 * 4", 14, false},
		{"parentheses override", "(2 + 3) * 4", 20, false},
		{"complex precedence", "2 + 3 * 4 - 5", 9, false},
		{"nested parentheses", "((2 + 3) * 4) + 1", 21, false},
		
		// Decimal numbers
		{"decimal addition", "2.5 + 3.5", 6, false},
		{"decimal multiplication", "2.5 * 2", 5, false},
		{"decimal division", "5.5 / 2", 2.75, false},
		
		// Negative numbers
		{"negative result", "3 - 5", -2, false},
		{"negative in expression", "2 + (-3)", -1, false},
		
		// Division by zero
		{"division by zero", "5 / 0", 0, true},
		
		// Invalid expressions
		{"invalid expression", "2 + * 3", 0, true},
		{"empty expression", "", 0, true},
		{"invalid characters", "2 + abc", 0, true},
		
		// Percentage
		{"percentage", "100 * 50%", 50, false},
		
		// Complex expressions
		{"complex 1", "2 * (3 + 4) / 2", 7, false},
		{"complex 2", "(10 + 5) * (20 - 8) / 3", 90, false},
		{"complex 3", "2 + 3 * 4 - 5 / 2", 11.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate(%q) error = %v, wantErr %v", tt.expression, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.expression, result, tt.expected)
			}
		})
	}
}

func TestEvaluateScientific(t *testing.T) {
	tests := []struct {
		name      string
		expression string
		expected  float64
		wantErr   bool
	}{
		// Scientific functions (assuming radians)
		{"sin 0", "sin(0)", 0, false},
		{"cos 0", "cos(0)", 1, false},
		{"sqrt 4", "sqrt(4)", 2, false},
		{"log 100", "log(100)", 2, false},
		{"ln e", "ln(e)", 1, false},
		{"exp 1", "exp(1)", math.E, false},
		{"pi", "pi", math.Pi, false},
		{"e", "e", math.E, false},
		{"factorial 5", "5!", 120, false},
		{"power", "2^3", 8, false},
		
		// Trigonometric with degrees
		{"sin 90 deg", "sin(90, deg)", 1, false},
		{"cos 0 deg", "cos(0, deg)", 1, false},
		
		// Invalid scientific
		{"invalid function", "foo(5)", 0, true},
		{"invalid factorial", "-5!", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate(%q) error = %v, wantErr %v", tt.expression, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !almostEqual(result, tt.expected) {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.expression, result, tt.expected)
			}
		})
	}
}

func TestOperatorPrecedence(t *testing.T) {
	// Test that * and / have higher precedence than + and -
	result, err := Evaluate("2 + 3 * 4")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if !almostEqual(result, 14) {
		t.Errorf("Operator precedence test failed: 2 + 3 * 4 = %v, want 14", result)
	}

	// Test parentheses override precedence
	result, err = Evaluate("(2 + 3) * 4")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if !almostEqual(result, 20) {
		t.Errorf("Parentheses override test failed: (2 + 3) * 4 = %v, want 20", result)
	}

	// Test nested parentheses
	result, err = Evaluate("2 * (3 + (4 * 5))")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if !almostEqual(result, 46) {
		t.Errorf("Nested parentheses test failed: 2 * (3 + (4 * 5)) = %v, want 46", result)
	}

	// Test left-to-right for same precedence
	result, err = Evaluate("10 - 3 - 2")
	if err != nil {
		t.Fatalf("Evaluate() error = %v", err)
	}
	if !almostEqual(result, 5) {
		t.Errorf("Left-to-right test failed: 10 - 3 - 2 = %v, want 5", result)
	}
}

func TestComplexExpression(t *testing.T) {
	// Test complex nested expressions
	tests := []struct {
		name      string
		expression string
		expected  float64
	}{
		{"nested 1", "2 * (3 + (4 * (5 - 1)))", 98, false},
		{"nested 2", "(10 + 20) * (30 - (40 / 2))", 600, false},
		{"multiple ops", "2 + 3 * 4 - 5 / 2 + 1", 12.5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.expression)
			if err != nil {
				t.Fatalf("Evaluate(%q) error = %v", tt.expression, err)
			}
			if !almostEqual(result, tt.expected) {
				t.Errorf("Complex expression %q = %v, want %v", tt.expression, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// Mode Management Tests
// =============================================================================

func TestCalculatorModes(t *testing.T) {
	// Test standard mode
	stdCalc := NewCalculator(ModeStandard)
	if stdCalc.GetMode() != ModeStandard {
		t.Errorf("Standard mode calculator has wrong mode: %v", stdCalc.GetMode())
	}

	// Test scientific mode
	sciCalc := NewCalculator(ModeScientific)
	if sciCalc.GetMode() != ModeScientific {
		t.Errorf("Scientific mode calculator has wrong mode: %v", sciCalc.GetMode())
	}
}

func TestSwitchMode(t *testing.T) {
	calc := NewCalculator(ModeStandard)
	
	// Start in standard mode
	if calc.GetMode() != ModeStandard {
		t.Errorf("Initial mode should be standard")
	}

	// Switch to scientific
	calc.SwitchMode(ModeScientific)
	if calc.GetMode() != ModeScientific {
		t.Errorf("Mode should be scientific after switch")
	}

	// Switch back to standard
	calc.SwitchMode(ModeStandard)
	if calc.GetMode() != ModeStandard {
		t.Errorf("Mode should be standard after switch back")
	}
}

func TestGetAvailableOperations(t *testing.T) {
	stdCalc := NewCalculator(ModeStandard)
	sciCalc := NewCalculator(ModeScientific)

	stdOps := stdCalc.GetAvailableOperations()
	sciOps := sciCalc.GetAvailableOperations()

	// Scientific should have more operations than standard
	if len(sciOps) <= len(stdOps) {
		t.Errorf("Scientific mode should have more operations than standard")
	}

	// Standard should have basic operations
	for _, op := range []string{"+", "-", "*", "/"} {
		if !contains(stdOps, op) {
			t.Errorf("Standard mode should have operation %q", op)
		}
	}

	// Scientific should have scientific functions
	for _, op := range []string{"sin", "cos", "tan", "sqrt", "log", "ln"} {
		if !contains(sciOps, op) {
			t.Errorf("Scientific mode should have operation %q", op)
		}
	}
}

// =============================================================================
// Undo/Redo Tests
// =============================================================================

func TestUndoRedo(t *testing.T) {
	calc := NewCalculator(ModeStandard)

	// Perform calculations
	_, _ = calc.Evaluate("2 + 2")
	_, _ = calc.Evaluate("3 * 4")
	_, _ = calc.Evaluate("5 - 1")

	// Get current result
	current := calc.GetCurrentResult()
	if !almostEqual(current, 4) {
		t.Errorf("Current result should be 4, got %v", current)
	}

	// Undo once
	calc.Undo()
	current = calc.GetCurrentResult()
	if !almostEqual(current, 12) {
		t.Errorf("After undo, result should be 12, got %v", current)
	}

	// Undo again
	calc.Undo()
	current = calc.GetCurrentResult()
	if !almostEqual(current, 4) {
		t.Errorf("After second undo, result should be 4, got %v", current)
	}

	// Redo
	calc.Redo()
	current = calc.GetCurrentResult()
	if !almostEqual(current, 12) {
		t.Errorf("After redo, result should be 12, got %v", current)
	}

	// Multiple undos beyond history
	calc.Undo()
	calc.Undo()
	calc.Undo() // Should not crash
	calc.Redo()
	current = calc.GetCurrentResult()
	if !almostEqual(current, 4) {
		t.Errorf("After excessive undo/redo, result should be 4, got %v", current)
	}
}

func TestMemoryOperations(t *testing.T) {
	calc := NewCalculator(ModeScientific)

	// Clear memory initially
	calc.MemoryClear()
	
	// Test M+
	calc.MemoryAdd(10)
	result := calc.MemoryRecall()
	if !almostEqual(result, 10) {
		t.Errorf("Memory recall should be 10, got %v", result)
	}

	// Test M+
	calc.MemoryAdd(5)
	result = calc.MemoryRecall()
	if !almostEqual(result, 15) {
		t.Errorf("Memory recall should be 15, got %v", result)
	}

	// Test M-
	calc.MemorySubtract(3)
	result = calc.MemoryRecall()
	if !almostEqual(result, 12) {
		t.Errorf("Memory recall should be 12, got %v", result)
	}

	// Test MC
	calc.MemoryClear()
	result = calc.MemoryRecall()
	if !almostEqual(result, 0) {
		t.Errorf("Memory recall should be 0 after clear, got %v", result)
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// almostEqual checks if two float64 values are approximately equal
func almostEqual(a, b float64) bool {
	const tolerance = 1e-10
	return math.Abs(a-b) < tolerance
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
