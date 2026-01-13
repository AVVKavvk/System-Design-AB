# Bloom Filter Implementation Guide

## Overview

A bloom filter is a space-efficient probabilistic data structure used to test whether an element is a member of a set.

### Key Characteristics

- Can definitively tell you that an element is **not** in the set
- Can tell you that an element **might be** in the set (with some probability of false positives)
- Never gives false negatives - if it says something isn't there, it definitely isn't
- Uses much less memory than storing the actual set

### How It Works

The bloom filter uses a bit array and multiple hash functions:

1. **Adding an element**: Hash it with each function and set those bit positions to 1
2. **Checking membership**: Hash the element again
   - If all positions are 1 → element might be in the set
   - If any position is 0 → element is definitely not in the set

### Common Use Cases

- Web browsers checking if a URL is in a list of malicious sites
- Databases avoiding expensive disk lookups for non-existent keys
- Distributed systems reducing network calls
- Spell checkers

### Trade-offs

The key trade-off is between memory usage and false positive rate. You can tune memory usage based on how many false positives you can tolerate.

---

## False Positive Rate Formula

For a bloom filter with:

- `m` = number of bits (8192 in our example)
- `n` = number of elements inserted
- `k` = number of hash functions (1 in our example)

The formula is:

```
FPR ≈ (1 - e^(-kn/m))^k
```

Where:

- `e` is Euler's number (~2.71828)
- `FPR` is the false positive rate

### Example Calculations (k=1, m=8192)

**After 100 items:**

```
FPR ≈ (1 - e^(-1×100/8192))^1
    ≈ (1 - e^(-0.0122))^1
    ≈ (1 - 0.9879)^1
    ≈ 0.0121 = 1.21%
```

**After 500 items:**

```
FPR ≈ (1 - e^(-500/8192))^1
    ≈ 0.0594 = 5.94%
```

**After 1000 items:**

```
FPR ≈ (1 - e^(-1000/8192))^1
    ≈ 0.1153 = 11.53%
```

---

## Implementation in Go

### Bit Array Structure

The implementation uses a 2D logical structure built on a 1D array of `uint64`:

```go
totalBits := sizeInKB * 8 * 1024  // Convert KB to bits
logicalColumn := 64                // Each uint64 holds 64 bits
logicalRow := totalBits / logicalColumn
bits := make([]uint64, logicalRow)
```

### Hash Index Mapping

To map a hash value to a specific bit:

1. Calculate the bit index: `index = hash % totalBits`
2. Find the row: `rowIndex = index / logicalColumn`
3. Find the column: `colIndex = index % logicalColumn`

**Example:**

- Hash value: 121212
- Total bits: 8192
- Bit index: 121212 % 8192 = 6524
- Row: 6524 / 64 = 101
- Column: 6524 % 64 = 60
- Verification: 101 × 64 + 60 = 6524 ✓

### Setting a Bit

To set the bit at position (row, col):

```go
bits[rowIndex] = bits[rowIndex] | (1 << colIndex)
```

This operation:

1. Creates a mask with bit set at `colIndex`: `1 << colIndex`
2. Performs bitwise OR to set that bit in the row

### Complete Example Code

```go
package main

import "fmt"

func AlgoDryRun(sizeInKB int) {
    totalBits := sizeInKB * 8 * 1024
    logicalColumn := 64
    logicalRow := totalBits / logicalColumn

    bits := make([]uint64, logicalRow)

    fmt.Printf("Total bits: %d\n", totalBits)
    fmt.Printf("Rows: %d, Columns: %d\n\n", logicalRow, logicalColumn)

    // Display initial state
    for row, word := range bits {
        fmt.Printf("Row %3d: %064b\n", row, word)
    }

    // Add element
    name := "vipin"
	fmt.Println(name)

	// Assume hash came from hash function is 121212
	hash := 121212
	index := hash % totalBits // 6524

	rowIndex := uint(index / logicalColumn) // 101
	colIndex := index % logicalColumn       // 60

	// 6524 bits means 101 row and 60 column  => 101*64 + 60 = 6524

	fmt.Printf("Row: %d, Column: %d\n", rowIndex, colIndex)

    // Set the bit
    bits[rowIndex] = bits[rowIndex] | (1 << colIndex)

    // Display updated state
    for row, word := range bits {
        fmt.Printf("Row %3d: %064b\n", row, word)
    }
}
```

### Visual Representation

For a 1KB bloom filter:

- Total bits: 8,192
- Rows: 128 (8192 / 64)
- Columns: 64 bits per row

When adding "vipin" with hash 121212:

```
Before: Row 101: 0000000000000000000000000000000000000000000000000000000000000000
After:  Row 101: 0001000000000000000000000000000000000000000000000000000000000000
   ↑ bit 60 set to 1
```

![](./images/1.png)
![](./images/2.png)
![](./images/3.png)
![](./images/4.png)
![](./images/5.png)

---

## Summary

This implementation demonstrates the core concept of bloom filters: using bit manipulation and hash functions to create a memory-efficient probabilistic data structure. The trade-off between space efficiency and accuracy makes bloom filters ideal for applications where false positives are acceptable but false negatives are not.
