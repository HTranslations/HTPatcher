package main

import (
	"encoding/json"
	"fmt"
	"htpatcher/internal/marshal"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rvdata2json <file.rvdata2> [output.json]")
		fmt.Println("If output is not specified, writes to <file>.json")
		os.Exit(1)
	}

	inputPath := os.Args[1]
	outputPath := ""

	if len(os.Args) >= 3 {
		outputPath = os.Args[2]
	} else {
		// Replace .rvdata2 extension with .json
		base := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))
		outputPath = base + ".json"
	}

	// Read the input file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Parse the marshal data
	parsed, err := marshal.Parse(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing marshal data: %v\n", err)
		os.Exit(1)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting to JSON: %v\n", err)
		os.Exit(1)
	}

	// Write the output file
	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Converted %s -> %s\n", inputPath, outputPath)
}
