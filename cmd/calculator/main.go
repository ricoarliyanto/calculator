package main

import (
	"fmt"
	"os"

	"github.com/ricoarliyanto/calculator/internal/ui"
)

func main() {
	// Execute the CLI
	if err := ui.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
