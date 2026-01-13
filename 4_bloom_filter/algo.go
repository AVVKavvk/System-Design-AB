package main

import "fmt"

func AlgoDryRun(sizeInKB int) {

	totalBits := sizeInKB * 8 * 1024 // size * 8 * 1024 bits

	var bits []uint64

	logicalColumn := 64 // Since we used 64 bits in each row uint64  - 32 for uint32
	logicalRow := totalBits / logicalColumn

	bits = make([]uint64, logicalRow)

	fmt.Printf("Total bits: %d\n", totalBits)
	fmt.Printf("Rows: %d, Columns: %d\n\n", logicalRow, logicalColumn)

	for row, word := range bits {
		fmt.Printf("Row %3d: %064b\n", row, word)
	}

	name := "vipin"
	fmt.Println(name)
	// Assume hash came from hash function is 121212
	hash := 121212
	index := hash % totalBits // 6524

	rowIndex := uint(index / logicalColumn) // 101
	colIndex := index % logicalColumn       // 60

	// 6524 bits means 101 row and 60 column  => 101*64 + 60 = 6524

	fmt.Printf("Row: %d, Column: %d\n", rowIndex, colIndex)

	bits[rowIndex] = bits[rowIndex] | (1 << colIndex)

	for row, word := range bits {
		fmt.Printf("Row %3d: %064b\n", row, word)
	}
}
