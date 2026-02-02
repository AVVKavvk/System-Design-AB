package algo

import (
	"bufio"
	"os"
	"strings"
)

var trie *Trie

func InitTrieWithAbuseWords() error {
	// TODO : In Production, read abuse words from DB/ File/ S3
	filePath := "abuse_word.txt"

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(strings.ToLower(scanner.Text()))
		if word != "" {
			// Add word to trie
			err = trie.Insert(word)
			if err != nil {
				return err
			}
		}
	}

	trie.PrintAllWords()
	// trie.PrintStructure(trie.root, "\n")

	return scanner.Err()
}

func isAbuseWord(word string, trie *Trie) bool {

	if word == "" || strings.Trim(word, " ") == "" {
		return false
	}
	return trie.IsExists(word)
}

func CheckAbuseAndGetNewMessage(message string) string {
	var result strings.Builder

	result.Grow(len(message))

	i := 0
	length := len(message)

	for i < length {
		// Skip leading spaces
		if message[i] == ' ' {
			result.WriteByte(' ')
			i++
			continue
		}

		j := i
		// Extract word using two pointers
		for j < length && message[j] != ' ' {
			j++
		}

		word := message[i:j]

		if isAbuseWord(word, trie) {
			result.WriteByte(word[0])
			result.WriteString(strings.Repeat("#", len(word)-1))
		} else {
			result.WriteString(word)
		}

		// Next word
		i = j

	}

	return result.String()

}

func init() {
	trie = NewTrieNode()
}
