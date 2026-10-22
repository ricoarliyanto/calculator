package calculator

import "testing"

// =============================================================================
// Performance Benchmarks
// =============================================================================

// BenchmarkSimpleAddition benchmarks simple addition evaluation
func BenchmarkSimpleAddition(b *testing.B) {
	calc := NewCalculator(ModeStandard)
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate("2 + 2")
	}
}

// BenchmarkComplexExpression benchmarks complex nested expression evaluation
func BenchmarkComplexExpression(b *testing.B) {
	calc := NewCalculator(ModeStandard)
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate("2 * (3 + 4) * (5 - 1) + 6 / 2")
	}
}

// BenchmarkScientificExpression benchmarks scientific expression evaluation
func BenchmarkScientificExpression(b *testing.B) {
	calc := NewCalculator(ModeScientific)
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate("sin(90, deg) + cos(0, deg) * 2 + tan(45, deg)")
	}
}

// BenchmarkLongExpression benchmarks very long expression evaluation
func BenchmarkLongExpression(b *testing.B) {
	calc := NewCalculator(ModeStandard)
	expr := "1 + 2 - 3 * 4 / 5 + 6 * 7 + 8 - 9 / 3 + 4 * 5 + 6 - 7 + 8 * 9 + 1 - 2 + 3 * 4 + 5 - 6 + 7 * 8 + 9 / 3 + 1"
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate(expr)
	}
}

// BenchmarkFactorial benchmarks factorial computation
func BenchmarkFactorial(b *testing.B) {
	calc := NewCalculator(ModeScientific)
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate("10!")
	}
}

// BenchmarkPower benchmarks power computation
func BenchmarkPower(b *testing.B) {
	calc := NewCalculator(ModeScientific)
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate("2^10")
	}
}

// BenchmarkNestedParentheses benchmarks deeply nested parentheses
func BenchmarkNestedParentheses(b *testing.B) {
	calc := NewCalculator(ModeStandard)
	expr := "((((1 + 2))))"
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate(expr)
	}
}

// BenchmarkManyOperators benchmarks expression with many operators
func BenchmarkManyOperators(b *testing.B) {
	calc := NewCalculator(ModeStandard)
	expr := "1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10 - 1 - 2 - 3 - 4 - 5"
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate(expr)
	}
}

// BenchmarkTrigonometricFunctions benchmarks trigonometric functions
func BenchmarkTrigonometricFunctions(b *testing.B) {
	calc := NewCalculator(ModeScientific)
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate("sin(1) + cos(1) + tan(1) + asin(1) + acos(1) + atan(1)")
	}
}

// BenchmarkMemoryOperations benchmarks memory store and recall operations
func BenchmarkMemoryOperations(b *testing.B) {
	calc := NewCalculator(ModeStandard)
	for i := 0; i < b.N; i++ {
		calc.MemoryAdd(42.0)
		calc.MemoryRecall()
	}
}

// BenchmarkUndoRedo benchmarks undo/redo operations
func BenchmarkUndoRedo(b *testing.B) {
	calc := NewCalculator(ModeStandard)
	for i := 0; i < b.N; i++ {
		_, _ = calc.Evaluate("1 + 2")
		calc.Undo()
		calc.Redo()
	}
}

// BenchmarkTokenizer benchmarks the tokenizer
func BenchmarkTokenizer(b *testing.B) {
	expr := "sin(90, deg) + cos(0, deg) * 2^3 - 4! + ln(e) + log(100) + sqrt(16)"
	for i := 0; i < b.N; i++ {
		_, _ = tokenize(expr)
	}
}

// BenchmarkShuntingYard benchmarks the shunting yard algorithm
func BenchmarkShuntingYard(b *testing.B) {
	expr := "2 * (3 + 4) - 5 / (7 - 2) + 3 * 4"
	for i := 0; i < b.N; i++ {
		tokens, _ := tokenize(expr)
		_ = tokens
	}
}

// BenchmarkRPNEvaluation benchmarks the RPN evaluation
func BenchmarkRPNEvaluation(b *testing.B) {
	expr := "2 * (3 + 4) - 5 / (7 - 2) + 3 * 4"
	tokens, _ := tokenize(expr)
	for i := 0; i < b.N; i++ {
		_, _ = evaluateTokens(tokens, "rad")
	}
}
