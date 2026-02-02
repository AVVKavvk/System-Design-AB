package algo

import (
	"fmt"
	"strings"
)

type Node struct {
	children map[rune]*Node
	isEnd    bool
}

type Trie struct {
	root *Node
}

func NewTrieNode() *Trie {
	return &Trie{
		root: &Node{
			children: make(map[rune]*Node),
		},
	}
}

// Insert adds a word to the trie
func (t *Trie) Insert(word string) error {

	if word == "" || strings.Trim(word, " ") == "" {
		return fmt.Errorf("Word can't be empty")
	}
	// convert word to lowercase
	word = strings.ToLower(word)

	node := t.root

	for _, ch := range word {
		if _, exists := node.children[ch]; !exists {
			node.children[ch] = &Node{children: make(map[rune]*Node)}
		}
		node = node.children[ch]
	}

	node.isEnd = true

	return nil

}

func (t *Trie) IsExists(word string) bool {
	if word == "" || strings.Trim(word, " ") == "" {
		return false
	}
	// convert word to lowercase
	word = strings.ToLower(word)

	node := t.root

	for _, ch := range word {
		if _, exists := node.children[ch]; !exists {
			return false
		}
		node = node.children[ch]
	}

	return node.isEnd
}

// StartsWith returns true if there is any word in the trie that starts with the given prefix
func (t *Trie) StartsWith(prefix string) bool {
	if prefix == "" || strings.Trim(prefix, " ") == "" {
		return false
	}
	// convert word to lowercase
	prefix = strings.ToLower(prefix)

	node := t.root
	for _, char := range prefix {
		if _, ok := node.children[char]; !ok {
			return false
		}
		node = node.children[char]
	}
	return true
}

// PrintAllWords starts the recursion from the root
func (t *Trie) PrintAllWords() {
	t.printHelper(t.root, "")
}

// printHelper traverses the trie and builds words character by character
func (t *Trie) printHelper(node *Node, currentWord string) {
	if node.isEnd {
		fmt.Println(currentWord)
	}

	for char, nextNode := range node.children {
		// Concatenate the character to the string and go deeper
		t.printHelper(nextNode, currentWord+string(char))
	}
}

func (t *Trie) PrintStructure(node *Node, indent string) {
	for char, nextNode := range node.children {
		fmt.Printf("%s%c\n", indent, char)
		t.PrintStructure(nextNode, indent+"  ")
	}
}
