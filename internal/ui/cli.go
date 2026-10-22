package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ricoarliyanto/calculator/internal/calculator"
	"github.com/ricoarliyanto/calculator/internal/config"
	"github.com/ricoarliyanto/calculator/internal/database"
	"github.com/ricoarliyanto/calculator/internal/history"
	"github.com/spf13/cobra"
)

var (
	// Global variables for CLI
	dataDir        string
	configManager  *config.ConfigManager
	db            *database.Database
	historyManager *history.HistoryManager
	currentCalc   *calculator.Calculator
)

// RootCmd represents the base command when called without any subcommand
var RootCmd = &cobra.Command{
	Use:   "calc",
	Short: "A powerful CLI calculator with standard and scientific modes",
	Long: `Calculator Application

A command-line calculator that supports:
- Standard mode: Basic arithmetic (+, -, *, /)
- Scientific mode: Trigonometric, logarithmic, exponential functions
- History: Save, view, filter, and export calculation history
- Config: Manage application settings

Examples:
  calc                      # Start interactive mode
  calc --mode scientific   # Start in scientific mode
  calc history             # View calculation history
  calc history --filter=sin # Filter history by keyword
  calc settings            # View/edit settings`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize application
		return initApplication(cmd)
	},
}

// Execute executes the root command
func Execute() error {
	return RootCmd.Execute()
}

// initApplication initializes the application components
func initApplication(cmd *cobra.Command) error {
	// Determine data directory
	if dataDir == "" {
		dataDir = getDefaultDataDir()
	}

	// Initialize database
	var err error
	db, err = database.NewDatabase(dataDir), nil
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	if err := db.Connect(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.RunMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Initialize history manager
	historyManager = history.NewHistoryManager(db)

	// Initialize config manager
	configManager = config.NewConfigManager(dataDir)

	// Load config
	if err := configManager.Load(); err != nil {
		// Non-fatal, use defaults
		fmt.Printf("Warning: failed to load config: %v\n", err)
	}

	// Get mode from command or config
	mode := configManager.GetMode()
	if modeFlag, _ := cmd.Flags().GetString("mode"); modeFlag != "" {
		mode = modeFlag
	}

	// Create calculator
	currentCalc = calculator.NewCalculator(calculator.Mode(mode))

	return nil
}

// getDefaultDataDir returns the default data directory
func getDefaultDataDir() string {
	// Try APPDATA on Windows
	if appData := os.Getenv("APPDATA"); appData != "" {
		return filepath.Join(appData, "calculator")
	}

	// Try HOME on Unix-like
	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, ".local", "share", "calculator")
	}

	// Fallback to current directory
	return filepath.Join(".", "data")
}

// cleanupApplication cleans up application resources
func cleanupApplication() {
	if db != nil {
		db.Close()
	}
}

func init() {
	// Root command flags
	RootCmd.PersistentFlags().StringVar(&dataDir, "data-dir", "", "Directory for data files")
	RootCmd.PersistentFlags().StringP("mode", "m", "", "Calculator mode (standard/scientific)")

	// Add subcommands
	RootCmd.AddCommand(
		newCalcCommand(),
		newHistoryCommand(),
		newSettingsCommand(),
		newConfigCommand(),
		newVersionCommand(),
	)

	// Handle cleanup on exit
	cobra.OnFinalize(func() {
		cleanupApplication()
	})
}

// =============================================================================
// Calc Command (Interactive Mode)
// =============================================================================

func newCalcCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "calc",
		Short: "Start interactive calculator",
		Long: `Start the calculator in interactive mode.

Examples:
  calc                    # Start in default mode
  calc --mode scientific  # Start in scientific mode`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get mode from flag
			mode, _ := cmd.Flags().GetString("mode")
			if mode != "" {
				currentCalc.SwitchMode(calculator.Mode(mode))
			}

			// Start interactive mode
			startInteractiveMode()
		},
	}
}

// startInteractiveMode starts the interactive calculator
func startInteractiveMode() {
	fmt.Printf("\nCalculator Application (v%s)\n", getVersion())
	fmt.Printf("Mode: %s\n", currentCalc.GetMode())
	fmt.Printf("Type 'help' for available commands, 'quit' or 'exit' to quit\n\n")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(">>> ")
		if !scanner.Scan() {
			// EOF (Ctrl+D)
			fmt.Println()
			break
		}

		input := strings.TrimSpace(scanner.Text())

		// Skip empty input
		if input == "" {
			continue
		}

		// Handle special commands
		switch strings.ToLower(input) {
		case "quit", "exit":
			fmt.Println("Goodbye!")
			return

		case "help":
			printInteractiveHelp()
			continue

		case "clear":
			currentCalc = calculator.NewCalculator(currentCalc.GetMode())
			fmt.Println("Calculator cleared")
			continue

		case "mode":
			fmt.Printf("Current mode: %s\n", currentCalc.GetMode())
			continue

		case "mode standard":
			currentCalc.SwitchMode(calculator.ModeStandard)
			fmt.Println("Switched to standard mode")
			continue

		case "mode scientific":
			currentCalc.SwitchMode(calculator.ModeScientific)
			fmt.Println("Switched to scientific mode")
			continue

		case "history":
			printHistory()
			continue

		case "undo":
			if currentCalc.CanUndo() {
				currentCalc.Undo()
				fmt.Printf("Undo: %v\n", currentCalc.GetCurrentResult())
			} else {
				fmt.Println("Nothing to undo")
			}
			continue

		case "redo":
			currentCalc.Redo()
			fmt.Printf("Redo: %v\n", currentCalc.GetCurrentResult())
			continue

		case "mc":
			currentCalc.MemoryClear()
			fmt.Println("Memory cleared")
			continue

		case "mr":
			fmt.Printf("Memory: %v\n", currentCalc.MemoryRecall())
			continue

		case "m+":
			fmt.Print("Enter value to add to memory: ")
			if scanner.Scan() {
				var val float64
				_, err := fmt.Sscanf(scanner.Text(), "%f", &val)
				if err != nil {
					fmt.Println("Invalid number")
				} else {
					currentCalc.MemoryAdd(val)
					fmt.Printf("Added %v to memory\n", val)
				}
			}
			continue

		case "m-":
			fmt.Print("Enter value to subtract from memory: ")
			if scanner.Scan() {
				var val float64
				_, err := fmt.Sscanf(scanner.Text(), "%f", &val)
				if err != nil {
					fmt.Println("Invalid number")
				} else {
					currentCalc.MemorySubtract(val)
					fmt.Printf("Subtracted %v from memory\n", val)
				}
			}
			continue
		}

		// Evaluate expression
		result, err := currentCalc.Evaluate(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Save to history
		_, err = historyManager.AddCalculation(input, fmt.Sprintf("%v", result), database.CalculationMode(currentCalc.GetMode()))
		if err != nil {
			// Non-fatal, just log
			fmt.Printf("Warning: failed to save to history: %v\n", err)
		}

		fmt.Printf("Result: %v\n", result)
	}
}

// printInteractiveHelp prints help for interactive mode
func printInteractiveHelp() {
	fmt.Print(`Interactive Calculator Commands:

Calculator Operations:
  Enter any mathematical expression to evaluate it
  Example: 2 + 2, sin(90, deg), sqrt(16)

Special Commands:
  help        - Show this help
  quit/exit   - Quit the calculator
  clear       - Clear the calculator
  mode        - Show current mode
  mode standard   - Switch to standard mode
  mode scientific - Switch to scientific mode
  history     - Show calculation history
  undo        - Undo last operation
  redo        - Redo last undone operation

Memory Operations:
  mc          - Memory clear
  mr          - Memory recall
  m+          - Memory add
  m-          - Memory subtract

Standard Mode Operations: +, -, *, /, %, (, )
Scientific Mode Operations: sin, cos, tan, asin, acos, atan, log, ln, exp, sqrt, !, ^, pi, e

Note: Scientific functions use degrees or radians based on configuration.
`)
}

// printHistory prints calculation history
func printHistory() {
	history, err := historyManager.GetHistory()
	if err != nil {
		fmt.Printf("Error loading history: %v\n", err)
		return
	}

	if len(history) == 0 {
		fmt.Println("No history yet")
		return
	}

	fmt.Printf("\nCalculation History (last %d):\n", len(history))
	fmt.Println(strings.Repeat("=", 60))

	for i, calc := range history {
		fmt.Printf("%d. [%s] %s = %s\n", i+1, calc.Mode, calc.Expression, calc.Result)
	}
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total: %d calculations\n\n", len(history))
}

// =============================================================================
// History Command
// =============================================================================

func newHistoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "View and manage calculation history",
		Long: `View, filter, export, and clear calculation history.

Examples:
  calc history                  # View all history
  calc history --limit=10      # View last 10 entries
  calc history --filter=sin   # Filter by keyword
  calc history --export=csv   # Export to CSV
  calc history --export=json  # Export to JSON
  calc history --clear        # Clear all history`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get flags
			limit, _ := cmd.Flags().GetInt("limit")
			filter, _ := cmd.Flags().GetString("filter")
			exportFormat, _ := cmd.Flags().GetString("export")
			clear, _ := cmd.Flags().GetBool("clear")

			if clear {
				if err := clearHistory(); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("History cleared successfully")
				return
			}

			if exportFormat != "" {
				if err := exportHistory(exportFormat); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("History exported to %s format\n", exportFormat)
				return
			}

			// View history
			if err := viewHistory(limit, filter); err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().IntP("limit", "l", 0, "Limit number of entries to display (0 = all)")
	cmd.Flags().StringP("filter", "f", "", "Filter by keyword")
	cmd.Flags().StringP("export", "e", "", "Export format: csv or json")
	cmd.Flags().BoolP("clear", "c", false, "Clear all history")

	return cmd
}

// viewHistory displays calculation history
func viewHistory(limit int, filter string) error {
	var calculations []database.Calculation
	var err error

	if filter != "" {
		calcFilter := database.CalculationFilter{
			Search: &filter,
		}
		calculations, err = historyManager.GetHistoryByFilter(calcFilter)
	} else {
		calculations, err = historyManager.GetHistory()
	}

	if err != nil {
		return err
	}

	// Apply limit
	if limit > 0 && limit < len(calculations) {
		calculations = calculations[:limit]
	}

	if len(calculations) == 0 {
		fmt.Println("No history found")
		return nil
	}

	fmt.Printf("\nCalculation History\n")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-5s %-12s %-40s %-20s %s\n", "ID", "Mode", "Expression", "Result", "Created")
	fmt.Println(strings.Repeat("-", 80))

	for _, calc := range calculations {
		fmt.Printf("%-5d %-12s %-40s %-20s %s\n",
			calc.ID, calc.Mode, calc.Expression, calc.Result, calc.CreatedAt.Format("2006-01-02 15:04:05"))
	}
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("Total: %d calculations\n\n", len(calculations))

	return nil
}

// clearHistory clears all calculation history
func clearHistory() error {
	return historyManager.ClearHistory()
}

// exportHistory exports history to a file
func exportHistory(format string) error {
	var exportFormat database.ExportFormat
	switch format {
	case "csv":
		exportFormat = database.ExportCSV
	case "json":
		exportFormat = database.ExportJSON
	default:
		return fmt.Errorf("unsupported export format: %s (use csv or json)", format)
	}

	filepath, err := historyManager.ExportHistoryToFile(exportFormat, dataDir)
	if err != nil {
		return err
	}

	fmt.Printf("History exported to: %s\n", filepath)
	return nil
}

// =============================================================================
// Settings Command
// =============================================================================

func newSettingsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "View and edit application settings",
		Long: `View and modify application configuration settings.

Examples:
  calc settings              # View all settings
  calc settings --set=mode=scientific  # Set mode to scientific
  calc settings --set=decimal-places=5 # Set decimal places`,
		Run: func(cmd *cobra.Command, args []string) {
			set, _ := cmd.Flags().GetStringSlice("set")
			get, _ := cmd.Flags().GetString("get")
			list, _ := cmd.Flags().GetBool("list")
			reset, _ := cmd.Flags().GetBool("reset")

			if reset {
				configManager.Reset()
				if err := configManager.Save(); err != nil {
					fmt.Printf("Error saving config: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("Settings reset to defaults")
				return
			}

			if list || (get == "" && len(set) == 0) {
				printSettings()
				return
			}

			if get != "" {
				printSetting(get)
				return
			}

			for _, setting := range set {
				if err := parseAndSetSetting(setting); err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
			}

			if err := configManager.Save(); err != nil {
				fmt.Printf("Error saving config: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Settings updated successfully")
			printSettings()
		},
	}
}

// parseAndSetSetting parses and sets a setting in the format "key=value"
func parseAndSetSetting(setting string) error {
	parts := strings.SplitN(setting, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid setting format: %s (use key=value)", setting)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	switch key {
	case "mode":
		return configManager.SetMode(value)
	case "angle-unit":
		return configManager.SetAngleUnit(value)
	case "decimal-places":
		var places int
		_, err := fmt.Sscanf(value, "%d", &places)
		if err != nil {
			return fmt.Errorf("invalid decimal places value: %s", value)
		}
		return configManager.SetDecimalPlaces(places)
	case "history-limit":
		var limit int
		_, err := fmt.Sscanf(value, "%d", &limit)
		if err != nil {
			return fmt.Errorf("invalid history limit value: %s", value)
		}
		return configManager.SetHistoryLimit(limit)
	case "auto-save":
		configManager.SetAutoSave(value == "true")
		return nil
	case "show-timestamps":
		configManager.SetShowTimestamps(value == "true")
		return nil
	default:
		return fmt.Errorf("unknown setting: %s", key)
	}
}

// printSettings prints all settings
func printSettings() {
	cfg := configManager.ExportConfig()

	fmt.Println("\nCurrent Settings:")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("  Mode:           %v\n", cfg["mode"])
	fmt.Printf("  Angle Unit:     %v\n", cfg["angle_unit"])
	fmt.Printf("  Decimal Places: %v\n", cfg["decimal_places"])
	fmt.Printf("  History Limit:  %v\n", cfg["history_limit"])
	fmt.Printf("  Auto Save:      %v\n", cfg["auto_save"])
	fmt.Printf("  Show Timestamps:%v\n", cfg["show_timestamps"])
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println()
}

// printSetting prints a specific setting
func printSetting(key string) {
	cfg := configManager.ExportConfig()

	value, ok := cfg[key]
	if !ok {
		fmt.Printf("Unknown setting: %s\n", key)
		return
	}

	fmt.Printf("%s: %v\n", key, value)
}

// =============================================================================
// Config Command
// =============================================================================

func newConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Manage configuration file",
		Long: `Manage the application configuration file.

Examples:
  calc config path      # Show config file path
  calc config edit      # Edit config file in editor
  calc config reset     # Reset config to defaults`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Printf("Config file: %s\n", configManager.GetConfigPath())
				return
			}

			switch args[0] {
			case "path":
				fmt.Printf("Config file: %s\n", configManager.GetConfigPath())
			case "edit":
				fmt.Printf("Edit config file: %s\n", configManager.GetConfigPath())
				// In a real implementation, you could use the user's editor
				fmt.Println("Open the file in your editor to modify settings")
			case "reset":
				configManager.Reset()
				if err := configManager.Save(); err != nil {
					fmt.Printf("Error resetting config: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("Config reset to defaults")
			default:
				fmt.Printf("Unknown config command: %s\n", args[0])
				fmt.Println("Usage: calc config [path|edit|reset]")
			}
		},
	}
}

// =============================================================================
// Version Command
// =============================================================================

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print the version of the calculator application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Calculator Application v%s\n", getVersion())
		},
	}
}

// getVersion returns the application version
func getVersion() string {
	return "0.1.0"
}

// Close closes the CLI and cleans up resources
func Close() {
	cleanupApplication()
}

// GetCalculator returns the current calculator instance
func GetCalculator() *calculator.Calculator {
	return currentCalc
}

// GetHistoryManager returns the history manager instance
func GetHistoryManager() *history.HistoryManager {
	return historyManager
}

// GetConfigManager returns the config manager instance
func GetConfigManager() *config.ConfigManager {
	return configManager
}

// GetDB returns the database instance
func GetDB() *database.Database {
	return db
}

// GetDataDir returns the current data directory
func GetDataDir() string {
	return dataDir
}

// SetDataDir sets the data directory
func SetDataDir(dir string) {
	dataDir = dir
}
