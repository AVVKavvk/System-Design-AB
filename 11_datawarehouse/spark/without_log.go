package main

// import (
// 	"context"
// 	"encoding/json"
// 	"flag"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"regexp"
// 	"strings"

// 	"github.com/apache/beam/sdks/v2/go/pkg/beam"
// 	"github.com/apache/beam/sdks/v2/go/pkg/beam/io/textio"
// 	"github.com/apache/beam/sdks/v2/go/pkg/beam/transforms/stats"
// 	"github.com/apache/beam/sdks/v2/go/pkg/beam/x/beamx"
// )

// var (
// 	input      = flag.String("input", "./data/input.txt", "Input file")
// 	outputFile = flag.String("output", "./output/result.json", "Output JSON file")
// )

// var wordRE = regexp.MustCompile(`[a-zA-Z]+('[a-z])?`)

// type WordCount struct {
// 	Word  string `json:"word"`
// 	Count int    `json:"count"`
// }

// func extractWords(line string, emit func(string)) {
// 	for _, word := range wordRE.FindAllString(line, -1) {
// 		emit(strings.ToLower(word))
// 	}
// }

// func main() {
// 	flag.Parse()
// 	beam.Init()

// 	ctx := context.Background()

// 	// Create pipeline
// 	p := beam.NewPipeline()
// 	s := p.Root()

// 	// Read input
// 	lines := textio.Read(s, *input)

// 	// Extract words
// 	words := beam.ParDo(s, extractWords, lines)

// 	// Count words
// 	counted := stats.Count(s, words)

// 	// Format as JSON lines (one JSON object per line)
// 	formatted := beam.ParDo(s, func(word string, count int) string {
// 		wc := WordCount{Word: word, Count: count}
// 		jsonBytes, _ := json.Marshal(wc)
// 		return string(jsonBytes)
// 	}, counted)

// 	// Write to file
// 	textio.Write(s, *outputFile, formatted)

// 	// Execute pipeline
// 	log.Println("Starting pipeline execution...")
// 	if err := beamx.Run(ctx, p); err != nil {
// 		log.Fatalf("Failed to execute job: %v", err)
// 	}

// 	log.Println("Pipeline completed! Output written to:", *outputFile)

// 	// Convert JSON lines to single JSON array
// 	convertToJSONArray(*outputFile)
// }

// func convertToJSONArray(outputPattern string) {
// 	// Read all shard files and combine into single JSON array
// 	files, err := filepath.Glob(outputPattern + "*")
// 	if err != nil {
// 		log.Printf("Warning: Could not glob output files: %v", err)
// 		return
// 	}

// 	if len(files) == 0 {
// 		log.Printf("Warning: No output files found matching pattern: %s", outputPattern)
// 		return
// 	}

// 	var wordCounts []WordCount
// 	for _, file := range files {
// 		log.Printf("Reading shard: %s", file)
// 		data, err := os.ReadFile(file)
// 		if err != nil {
// 			log.Printf("Error reading file %s: %v", file, err)
// 			continue
// 		}

// 		lines := strings.Split(string(data), "\n")
// 		for _, line := range lines {
// 			if strings.TrimSpace(line) == "" {
// 				continue
// 			}
// 			var wc WordCount
// 			if err := json.Unmarshal([]byte(line), &wc); err == nil {
// 				wordCounts = append(wordCounts, wc)
// 			} else {
// 				log.Printf("Error unmarshaling line: %s, error: %v", line, err)
// 			}
// 		}
// 	}

// 	if len(wordCounts) == 0 {
// 		log.Printf("Warning: No word counts found")
// 		return
// 	}

// 	// Write as single JSON array
// 	finalJSON, err := json.MarshalIndent(wordCounts, "", "  ")
// 	if err != nil {
// 		log.Printf("Error marshaling JSON: %v", err)
// 		return
// 	}

// 	finalFile := strings.TrimSuffix(outputPattern, filepath.Ext(outputPattern)) + "_final.json"
// 	if err := os.WriteFile(finalFile, finalJSON, 0644); err != nil {
// 		log.Printf("Error writing final JSON: %v", err)
// 		return
// 	}

// 	log.Printf("✓ Final JSON array written to: %s", finalFile)
// 	log.Printf("✓ Total word count entries: %d", len(wordCounts))
// }
