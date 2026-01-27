package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apache/beam/sdks/v2/go/pkg/beam"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
)

var (
	input      = flag.String("input", "./data/input.txt", "Input file")
	outputFile = flag.String("output", "./output/result.json", "Output JSON file")
)

var wordRE = regexp.MustCompile(`[a-zA-Z]+('[a-z])?`)

type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// Step 1: Log lines as they're read
func logLine(line string, emit func(string)) {
	fmt.Printf("STEP 1 [Read Line]: %q\n", line)
	emit(line)
}

// Step 2: Extract and log words
func extractWords(line string, emit func(string)) {
	fmt.Printf("STEP 2 [Extract Words] Input line: %q\n", line)
	words := wordRE.FindAllString(line, -1)
	fmt.Printf("        Found %d words: %v\n", len(words), words)

	for _, word := range words {
		lowerWord := strings.ToLower(word)
		fmt.Printf("        Emitting word: %q\n", lowerWord)
		emit(lowerWord)
	}
}

// Step 3: Log individual words before counting
func logWord(word string, emit func(string)) {
	fmt.Printf("STEP 3 [Word Stream]: %q\n", word)
	emit(word)
}

// Step 4: Log word counts after aggregation
func logCount(word string, count int, emit func(string, int)) {
	fmt.Printf("STEP 4 [Word Count]: word=%q, count=%d\n", word, count)
	emit(word, count)
}

// Step 5: Format and log JSON output
func formatJSON(word string, count int, emit func(string)) {
	wc := WordCount{Word: word, Count: count}
	jsonBytes, _ := json.Marshal(wc)
	jsonStr := string(jsonBytes)
	fmt.Printf("STEP 5 [Format JSON]: word=%q, count=%d => %s\n", word, count, jsonStr)
	emit(jsonStr)
}

// Step 6: Log final output before writing
func logFinalOutput(jsonLine string, emit func(string)) {
	fmt.Printf("STEP 6 [Write Output]: %s\n", jsonLine)
	emit(jsonLine)
}

func main() {
	flag.Parse()
	beam.Init()

	ctx := context.Background()

	// Create pipeline
	p := beam.NewPipeline()
	s := p.Root()

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("PIPELINE EXECUTION - DETAILED LOGGING")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	// STEP 1: Read input
	fmt.Println(">>> STAGE 1: Reading input file...")
	lines := textio.Read(s, *input)
	linesLogged := beam.ParDo(s, logLine, lines)

	// STEP 2: Extract words
	fmt.Println("\n>>> STAGE 2: Extracting words from lines...")
	words := beam.ParDo(s, extractWords, linesLogged)

	// STEP 3: Log individual words
	fmt.Println("\n>>> STAGE 3: Individual words in stream...")
	wordsLogged := beam.ParDo(s, logWord, words)

	// STEP 4: Count words
	fmt.Println("\n>>> STAGE 4: Counting word occurrences...")
	counted := stats.Count(s, wordsLogged)
	countedLogged := beam.ParDo(s, logCount, counted)

	// STEP 5: Format as JSON
	fmt.Println("\n>>> STAGE 5: Formatting to JSON...")
	formatted := beam.ParDo(s, formatJSON, countedLogged)

	// STEP 6: Log final output
	fmt.Println("\n>>> STAGE 6: Writing to output file...")
	finalOutput := beam.ParDo(s, logFinalOutput, formatted)

	// Write to file
	textio.Write(s, *outputFile, finalOutput)

	// Execute pipeline
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("STARTING BEAM PIPELINE EXECUTION")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	if err := beamx.Run(ctx, p); err != nil {
		log.Fatalf("Failed to execute job: %v", err)
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("PIPELINE COMPLETED SUCCESSFULLY")
	fmt.Println(strings.Repeat("=", 80) + "\n")

	log.Println("✓ Output written to:", *outputFile)

	// Convert JSON lines to single JSON array
	convertToJSONArray(*outputFile)
}

func convertToJSONArray(outputPattern string) {
	files, err := filepath.Glob(outputPattern + "*")
	if err != nil {
		log.Printf("Warning: Could not glob output files: %v", err)
		return
	}

	if len(files) == 0 {
		log.Printf("Warning: No output files found matching pattern: %s", outputPattern)
		return
	}

	var wordCounts []WordCount
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			var wc WordCount
			if err := json.Unmarshal([]byte(line), &wc); err == nil {
				wordCounts = append(wordCounts, wc)
			}
		}
	}

	if len(wordCounts) == 0 {
		return
	}

	finalJSON, _ := json.MarshalIndent(wordCounts, "", "  ")
	finalFile := strings.TrimSuffix(outputPattern, filepath.Ext(outputPattern)) + "_final.json"
	os.WriteFile(finalFile, finalJSON, 0644)

	fmt.Println("\n>>> POST-PROCESSING: Combining shards into single JSON file...")
	log.Printf("✓ Final JSON array written to: %s", finalFile)
	log.Printf("✓ Total word count entries: %d", len(wordCounts))

	fmt.Println("\n=== Summary of Results ===")
	for _, wc := range wordCounts {
		fmt.Printf("  • %s: %d\n", wc.Word, wc.Count)
	}
}
