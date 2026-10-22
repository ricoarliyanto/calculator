# Example Usage Scenarios

This document provides practical example scenarios for using the Calculator application.

## Scenario 1: Basic Arithmetic

```bash
# Start the calculator
./calculator

# In interactive mode:
> 123 + 456
579
> 7 * 8
56
> 100 / 4
25
> 10 - 3
7
> exit
```

## Scenario 2: Complex Expressions with Parentheses

```bash
# Start in scientific mode
./calculator --mode scientific

# In interactive mode:
> (10 + 5) * (20 - 8) / 3
60
> 2 ^ 10
1024
> ((1 + 2) * (3 + 4))
21
> 100 + 10%
110
> exit
```

## Scenario 3: Scientific Calculations

```bash
# Start in scientific mode
./calculator --mode scientific

# In interactive mode:
> sin(90, deg)
1
> cos(0, deg)
1
> tan(45, deg)
1
> asin(1, deg)
90
> sqrt(16)
4
> 5!
120
> log(100)
2
> ln(e)
1
> pi
3.141592653589793
> exp(1)
2.718281828459045
> exit
```

## Scenario 4: Trigonometric Functions in Radians

```bash
# Start in scientific mode (default angle unit is radians)
./calculator --mode scientific

# In interactive mode:
> sin(pi/2)
1
> cos(0)
1
> tan(pi/4)
1
> exit
```

## Scenario 5: Memory Operations

```bash
# Start in standard mode
./calculator

# In interactive mode:
> 10
10
> m+
> 5
5
> m+
> mr
15
> 3
3
> m-
> mr
12
> mc
> mr
0
> exit
```

## Scenario 6: History Management

```bash
# Do some calculations first
./calculator
> 2 + 2
> 3 * 3
> sin(90, deg)  (scientific mode)
> 10 / 0        (triggers error)
> exit

# View history
./calculator history

# Filter by scientific mode
./calculator history --mode scientific

# Search history
./calculator history --search "sin"

# Export to JSON
./calculator history --export json

# Export to CSV
./calculator history --export csv

# Delete a specific entry
./calculator history --delete 3

# Clear all history
./calculator history --clear

# View history statistics
./calculator history --stats
```

## Scenario 7: Settings Management

```bash
# View current settings
./calculator settings

# Change calculator mode
./calculator settings set mode=scientific

# Change angle unit
./calculator settings set angle-unit=deg

# Change decimal places
./calculator settings set decimal-places=4

# Enable/disable auto-save
./calculator settings set auto-save=true

# Change theme
./calculator settings set theme=dark

# Reset to defaults
./calculator settings reset

# View specific setting
./calculator settings get mode
```

## Scenario 8: Configuration Management

```bash
# View config file location
./calculator config path

# Ensure config file exists with defaults
./calculator config ensure

# Reset config to defaults
./calculator config reset

# View config file contents
cat ~/.config/calculator/config.yaml
```

## Scenario 9: Undo and Redo

```bash
# Start in standard mode
./calculator

# In interactive mode:
> 2 + 2
4
> 3 * 3
9
> 4 - 1
3
> undo
3       (result reverted to previous: 9)
> undo
4       (result reverted to 4)
> redo
3       (result restored to 3)
> redo
9       (result restored to 9)
> exit
```

## Scenario 10: Non-Interactive Mode

The calculator can be used in non-interactive mode for scripting:

```bash
# Get version information
./calculator version

# View settings
./calculator settings

# View history
./calculator history

# Export history to file
./calculator history --export json

# Clear history (use with caution!)
./calculator history --clear
```

## Scenario 11: Error Handling

```bash
# Division by zero
> 1 / 0
Error: division by zero

# Invalid expression
> 2 ++ 2
Error: invalid expression

# Unbalanced parentheses
> (2 + 3
Error: unbalanced parentheses

# Scientific function in standard mode
> sin(90)    (in standard mode)
Error: scientific function in standard mode

# Unknown identifier
> abc
Error: unknown identifier: abc

# Factorial of negative number (via unary)
> 5!
120

# Square root of negative number
> sqrt(-1)
Error: math domain error
```

## Scenario 12: Date Range History Filtering

```bash
# View calculations from the last 7 days
./calculator history --from "2024-01-01" --to "2024-12-31"

# View calculations from today onwards
./calculator history --from "$(date +%Y-%m-%d)"

# View calculations from last week
./calculator history --from "2024-01-01" --to "2024-01-07"
```

## Scenario 13: Data Directory Management

```bash
# Use a custom data directory
./calculator --data-dir /tmp/my_calculator_data

# This is useful for:
# - Running multiple calculator instances with different data
# - Using a temporary directory for testing
# - Portable installation on USB drives
```
