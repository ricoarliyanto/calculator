package calculator

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

// Calculator mode type
type Mode string

const (
	ModeStandard   Mode = "standard"
	ModeScientific Mode = "scientific"
)

// Errors
var (
	ErrDivisionByZero = errors.New("division by zero")
	ErrDomainError    = errors.New("domain error")
	ErrInvalidInput   = errors.New("invalid input")
	ErrInvalidExpression = errors.New("invalid expression")
)

// Mathematical constants
const (
	Pi = math.Pi
	E  = math.E
)

// Calculator represents a calculator with mode support and memory
type Calculator struct {
	mode          Mode
	memory        float64
	undoStack     []float64
	redoStack     []float64
	currentResult float64
	useRadians    bool
	decimalPlaces int
}

// NewCalculator creates a new calculator instance
func NewCalculator(mode Mode) *Calculator {
	return &Calculator{
		mode:          mode,
		memory:        0,
		undoStack:     make([]float64, 0),
		redoStack:     make([]float64, 0),
		currentResult: 0,
		useRadians:    true,
		decimalPlaces: 10,
	}
}

// GetMode returns the current calculator mode
func (c *Calculator) GetMode() Mode {
	return c.mode
}

// SwitchMode switches the calculator mode
func (c *Calculator) SwitchMode(mode Mode) {
	c.mode = mode
}

// SetAngleUnit sets whether to use radians or degrees
func (c *Calculator) SetAngleUnit(useRadians bool) {
	c.useRadians = useRadians
}

// SetDecimalPlaces sets the number of decimal places for display
func (c *Calculator) SetDecimalPlaces(places int) {
	c.decimalPlaces = places
}

// GetAvailableOperations returns available operations for the current mode
func (c *Calculator) GetAvailableOperations() []string {
	ops := []string{"+", "-", "*", "/", "(", ")"}

	if c.mode == ModeScientific {
		ops = append(ops,
			"sin", "cos", "tan", "asin", "acos", "atan",
			"log", "ln", "exp", "sqrt", "!", "^", "%",
			"pi", "e",
			"MC", "MR", "M+", "M-",
		)
	}

	if c.mode == ModeStandard {
		ops = append(ops, "%")
	}

	return ops
}

// Evaluate evaluates a mathematical expression
func (c *Calculator) Evaluate(expression string) (float64, error) {
	// Set angle unit for trigonometric functions
	angleUnit := "rad"
	if !c.useRadians {
		angleUnit = "deg"
	}

	result, err := EvaluateWithContext(expression, angleUnit)
	if err != nil {
		return 0, err
	}

	// Save to undo stack
	c.undoStack = append(c.undoStack, c.currentResult)
	c.currentResult = result
	c.redoStack = make([]float64, 0) // Clear redo stack on new operation

	return result, nil
}

// GetCurrentResult returns the current result
func (c *Calculator) GetCurrentResult() float64 {
	return c.currentResult
}

// Undo undoes the last operation
func (c *Calculator) Undo() {
	if len(c.undoStack) > 0 {
		lastIndex := len(c.undoStack) - 1
		c.redoStack = append(c.redoStack, c.currentResult)
		c.currentResult = c.undoStack[lastIndex]
		c.undoStack = c.undoStack[:lastIndex]
	}
}

// Redo redoes the last undone operation
func (c *Calculator) Redo() {
	if len(c.redoStack) > 0 {
		lastIndex := len(c.redoStack) - 1
		c.undoStack = append(c.undoStack, c.currentResult)
		c.currentResult = c.redoStack[lastIndex]
		c.redoStack = c.redoStack[:lastIndex]
	}
}

// MemoryClear clears the memory
func (c *Calculator) MemoryClear() {
	c.memory = 0
}

// MemoryRecall recalls the memory value
func (c *Calculator) MemoryRecall() float64 {
	return c.memory
}

// MemoryAdd adds to memory
func (c *Calculator) MemoryAdd(value float64) {
	c.memory += value
}

// MemorySubtract subtracts from memory
func (c *Calculator) MemorySubtract(value float64) {
	c.memory -= value
}

// =============================================================================
// Static Expression Evaluation
// =============================================================================

// Evaluate evaluates a mathematical expression string (uses radians for trig functions)
func Evaluate(expression string) (float64, error) {
	return EvaluateWithContext(expression, "rad")
}

// EvaluateWithContext evaluates with a specific angle unit
func EvaluateWithContext(expression string, angleUnit string) (float64, error) {
	if expression == "" {
		return 0, ErrInvalidExpression
	}

	// Preprocess: Replace special characters and symbols
	expression = strings.ReplaceAll(expression, " ", "")
	expression = strings.ReplaceAll(expression, "π", "pi")
	expression = strings.ReplaceAll(expression, "Π", "pi")

	// Tokenize the expression
	tokens, err := tokenize(expression)
	if err != nil {
		return 0, err
	}

	// Parse and evaluate using Shunting Yard algorithm
	return evaluateTokens(tokens, angleUnit)
}

// tokenType represents the type of a token
type tokenType int

const (
	tokenNumber tokenType = iota
	tokenOperator
	tokenFunction
	tokenLeftParen
	tokenRightParen
	tokenComma
)

// Token represents a token in the expression
type Token struct {
	Value string
	Type  tokenType
}

// tokenize converts an expression string into tokens
func tokenize(expression string) ([]Token, error) {
	var tokens []Token
	i := 0

	for i < len(expression) {
		char := rune(expression[i])

		// Skip whitespace
		if unicode.IsSpace(char) {
			i++
			continue
		}

		// Handle numbers (including decimals and scientific notation)
		if unicode.IsDigit(char) || char == '.' {
			start := i
			for i < len(expression) && (unicode.IsDigit(rune(expression[i])) || expression[i] == '.') {
				i++
			}
			tokens = append(tokens, Token{
				Value: expression[start:i],
				Type:  tokenNumber,
			})
			continue
		}

		// Handle functions (sin, cos, tan, log, ln, sqrt, etc.) and constants (pi, e)
		if unicode.IsLetter(char) {
			start := i
			for i < len(expression) && unicode.IsLetter(rune(expression[i])) {
				i++
			}
			name := expression[start:i]

			// Check if it's a constant
			if isConstant(name) {
				tokens = append(tokens, Token{
					Value: name,
					Type:  tokenNumber,
				})
				continue
			}

			// Check if it's a known function
			if isFunction(name) {
				tokens = append(tokens, Token{
					Value: name,
					Type:  tokenFunction,
				})
				continue
			}

			// Unknown identifier
			return nil, fmt.Errorf("unknown identifier: %s", name)
		}

		// Handle operators and special characters
		switch char {
		case '(', ')', ',':
			tokenType := tokenLeftParen
			if char == ')' {
				tokenType = tokenRightParen
			} else if char == ',' {
				tokenType = tokenComma
			}
			tokens = append(tokens, Token{
				Value: string(char),
				Type:  tokenType,
			})
			i++
			continue

		case '+', '-', '*', '/', '^':
			// Check for unary minus
			if char == '-' && (len(tokens) == 0 ||
				isUnaryContext(tokens[len(tokens)-1])) {
				// This is a unary minus
				// We'll handle this in the evaluation phase
				tokens = append(tokens, Token{
					Value: "u-",
					Type:  tokenOperator,
				})
				i++
				continue
			}
			
			// Regular binary operator
			tokens = append(tokens, Token{
				Value: string(char),
				Type:  tokenOperator,
			})
			i++
			continue

		case '%':
			// Percentage as unary postfix operator
			tokens = append(tokens, Token{
				Value: "percent",
				Type:  tokenFunction,
			})
			i++
			continue

		case '!':
			// Factorial as unary postfix operator
			tokens = append(tokens, Token{
				Value: "factorial",
				Type:  tokenFunction,
			})
			i++
			continue
		}

		return nil, fmt.Errorf("invalid character: %c", char)
	}

	return tokens, nil
}

// isUnaryContext checks if the previous token indicates a unary operator context
func isUnaryContext(prev Token) bool {
	return prev.Type == tokenOperator && prev.Value != ")" ||
		prev.Type == tokenLeftParen ||
		prev.Type == tokenFunction
}

// isFunction checks if a string is a known function
func isFunction(name string) bool {
	functions := []string{"sin", "cos", "tan", "asin", "acos", "atan", "log", "ln", "exp", "sqrt", "abs"}
	for _, f := range functions {
		if name == f {
			return true
		}
	}
	return false
}

// isConstant checks if a string is a known constant
func isConstant(name string) bool {
	return name == "pi" || name == "e"
}

// operatorPrecedence returns the precedence of an operator
func operatorPrecedence(op string) int {
	switch op {
	case "^":
		return 4
	case "*", "/":
		return 3
	case "+", "-":
		return 2
	case "u-":
		return 5 // Highest precedence for unary minus
	default:
		return 0
	}
}

// isRightAssociative returns true for right-associative operators
func isRightAssociative(op string) bool {
	return op == "^"
}

// evaluateTokens evaluates tokens using Shunting Yard algorithm
func evaluateTokens(tokens []Token, angleUnit string) (float64, error) {
	output := make([]Token, 0)
	ops := make([]Token, 0)

	for _, token := range tokens {
		switch token.Type {
		case tokenNumber:
			// Handle constants
			if isConstant(token.Value) {
				var value float64
				switch token.Value {
				case "pi":
					value = Pi
				case "e":
					value = E
				default:
					return 0, fmt.Errorf("unknown constant: %s", token.Value)
				}
				output = append(output, Token{
					Value: fmt.Sprintf("%v", value),
					Type:  tokenNumber,
				})
			} else {
				output = append(output, token)
			}

		case tokenFunction:
			ops = append(ops, token)

		case tokenComma:
			// Process until left paren
			for len(ops) > 0 && ops[len(ops)-1].Type != tokenLeftParen {
				output = append(output, ops[len(ops)-1])
				ops = ops[:len(ops)-1]
			}
			if len(ops) == 0 || ops[len(ops)-1].Type != tokenLeftParen {
				return 0, ErrInvalidExpression
			}

		case tokenLeftParen:
			ops = append(ops, token)

		case tokenRightParen:
			// Process until left paren
			foundLeftParen := false
			for len(ops) > 0 && !foundLeftParen {
				top := ops[len(ops)-1]
				ops = ops[:len(ops)-1]
				if top.Type == tokenLeftParen {
					foundLeftParen = true
				} else {
					output = append(output, top)
				}
			}
			if !foundLeftParen {
				return 0, ErrInvalidExpression
			}
			// Check if there's a function before the paren
			if len(ops) > 0 && ops[len(ops)-1].Type == tokenFunction {
				output = append(output, ops[len(ops)-1])
				ops = ops[:len(ops)-1]
			}

		case tokenOperator:
			// Process operators with higher or equal precedence
			for len(ops) > 0 {
				top := ops[len(ops)-1]
				if top.Type != tokenOperator && top.Type != tokenFunction {
					break
				}
				
				shouldPop := false
				if isRightAssociative(token.Value) {
					shouldPop = operatorPrecedence(token.Value) < operatorPrecedence(top.Value)
				} else {
					shouldPop = operatorPrecedence(token.Value) <= operatorPrecedence(top.Value)
				}
				
				if shouldPop {
					output = append(output, top)
					ops = ops[:len(ops)-1]
				} else {
					break
				}
			}
			ops = append(ops, token)
		}
	}

	// Process remaining operators
	for len(ops) > 0 {
		top := ops[len(ops)-1]
		ops = ops[:len(ops)-1]
		if top.Type == tokenLeftParen || top.Type == tokenRightParen {
			return 0, ErrInvalidExpression
		}
		output = append(output, top)
	}

	// Evaluate RPN (Reverse Polish Notation)
	return evaluateRPN(output, angleUnit)
}

// evaluateRPN evaluates tokens in Reverse Polish Notation
func evaluateRPN(tokens []Token, angleUnit string) (float64, error) {
	stack := make([]float64, 0)

	for _, token := range tokens {
		switch token.Type {
		case tokenNumber:
			value, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid number: %s", token.Value)
			}
			stack = append(stack, value)

		case tokenOperator:
			if token.Value == "u-" {
				// Unary minus
				if len(stack) < 1 {
					return 0, ErrInvalidExpression
				}
				a := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				stack = append(stack, -a)
			} else {
				// Binary operator
				if len(stack) < 2 {
					return 0, ErrInvalidExpression
				}
				b := stack[len(stack)-1]
				a := stack[len(stack)-2]
				stack = stack[:len(stack)-2]

				result, err := applyOperator(a, b, token.Value)
				if err != nil {
					return 0, err
				}
				stack = append(stack, result)
			}

		case tokenFunction:
			if len(stack) < 1 {
				return 0, ErrInvalidExpression
			}
			arg := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			result, err := applyFunction(token.Value, arg, angleUnit)
			if err != nil {
				return 0, err
			}
			stack = append(stack, result)

		default:
			return 0, ErrInvalidExpression
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}

// applyOperator applies a binary operator to two operands
func applyOperator(a, b float64, operator string) (float64, error) {
	switch operator {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	case "^":
		return math.Pow(a, b), nil
	default:
		return 0, fmt.Errorf("unknown operator: %s", operator)
	}
}

// applyFunction applies a function to an argument
func applyFunction(funcName string, arg float64, angleUnit string) (float64, error) {
	switch funcName {
	case "sin":
		return Sin(arg, angleUnit)
	case "cos":
		return Cos(arg, angleUnit)
	case "tan":
		return Tan(arg, angleUnit)
	case "asin":
		return ASin(arg, angleUnit)
	case "acos":
		return ACos(arg, angleUnit)
	case "atan":
		return ATan(arg, angleUnit)
	case "log":
		return Log(arg)
	case "ln":
		return Ln(arg)
	case "exp":
		return Exp(arg)
	case "sqrt":
		return Sqrt(arg)
	case "abs":
		return math.Abs(arg), nil
	case "percent":
		// Percentage: convert to decimal (divide by 100)
		return arg / 100, nil
	case "factorial":
		// Factorial - arg should be an integer
		if arg < 0 || arg != math.Floor(arg) {
			return 0, ErrDomainError
		}
		if arg > 20 {
			return 0, ErrInvalidInput // 21! overflows float64
		}
		return Factorial(int64(arg))
	default:
		return 0, fmt.Errorf("unknown function: %s", funcName)
	}
}
