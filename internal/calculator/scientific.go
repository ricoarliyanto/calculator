package calculator

import (
	"errors"
	"math"
)

// Scientific calculator operations

// ErrDomainError is returned for domain errors (e.g., sqrt of negative)
var ErrDomainError = errors.New("domain error")

// ErrInvalidInput is returned for invalid input
var ErrInvalidInput = errors.New("invalid input")

// Trigonometric constants
const (
	DegToRad = math.Pi / 180
	RadToDeg = 180 / math.Pi
)

// Mathematical constants
const (
	Pi = math.Pi
	E  = math.E
)

// Sin calculates sine of an angle
func Sin(angle float64, unit string) (float64, error) {
	if unit == "deg" || unit == "degree" {
		angle = angle * DegToRad
	}
	return math.Sin(angle), nil
}

// Cos calculates cosine of an angle
func Cos(angle float64, unit string) (float64, error) {
	if unit == "deg" || unit == "degree" {
		angle = angle * DegToRad
	}
	return math.Cos(angle), nil
}

// Tan calculates tangent of an angle
func Tan(angle float64, unit string) (float64, error) {
	if unit == "deg" || unit == "degree" {
		angle = angle * DegToRad
	}
	return math.Tan(angle), nil
}

// ASin calculates inverse sine (arcsine)
func ASin(value float64, unit string) (float64, error) {
	if value < -1 || value > 1 {
		return 0, ErrDomainError
	}
	result := math.Asin(value)
	if unit == "deg" || unit == "degree" {
		result = result * RadToDeg
	}
	return result, nil
}

// ACos calculates inverse cosine (arccosine)
func ACos(value float64, unit string) (float64, error) {
	if value < -1 || value > 1 {
		return 0, ErrDomainError
	}
	result := math.Acos(value)
	if unit == "deg" || unit == "degree" {
		result = result * RadToDeg
	}
	return result, nil
}

// ATan calculates inverse tangent (arctangent)
func ATan(value float64, unit string) (float64, error) {
	result := math.Atan(value)
	if unit == "deg" || unit == "degree" {
		result = result * RadToDeg
	}
	return result, nil
}

// Log calculates logarithm base 10
func Log(value float64) (float64, error) {
	if value <= 0 {
		return 0, ErrDomainError
	}
	return math.Log10(value), nil
}

// Ln calculates natural logarithm (base e)
func Ln(value float64) (float64, error) {
	if value <= 0 {
		return 0, ErrDomainError
	}
	return math.Log(value), nil
}

// Exp calculates exponential function (e^x)
func Exp(value float64) (float64, error) {
	return math.Exp(value), nil
}

// Sqrt calculates square root
func Sqrt(value float64) (float64, error) {
	if value < 0 {
		return 0, ErrDomainError
	}
	return math.Sqrt(value), nil
}

// Factorial calculates factorial of a non-negative integer
func Factorial(n int64) (float64, error) {
	if n < 0 {
		return 0, ErrDomainError
	}
	if n > 20 {
		// 21! overflows float64
		return 0, ErrInvalidInput
	}
	
	if n == 0 || n == 1 {
		return 1, nil
	}
	
	result := float64(1)
	for i := int64(2); i <= n; i++ {
		result *= float64(i)
	}
	return result, nil
}

// Power calculates base^exponent
func Power(base, exponent float64) (float64, error) {
	return math.Pow(base, exponent), nil
}

// Modulo calculates remainder of division
func Modulo(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return math.Mod(a, b), nil
}
