package calculator

// Standard calculator operations

// Add performs addition
func Add(a, b float64) (float64, error) {
	return a + b, nil
}

// Subtract performs subtraction
func Subtract(a, b float64) (float64, error) {
	return a - b, nil
}

// Multiply performs multiplication
func Multiply(a, b float64) (float64, error) {
	return a * b, nil
}

// Divide performs division
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// Percentage calculates percentage of a value
func Percentage(value, percent float64) (float64, error) {
	return (value * percent) / 100, nil
}

// Negate negates a value
func Negate(value float64) (float64, error) {
	return -value, nil
}
