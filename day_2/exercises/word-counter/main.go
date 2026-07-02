package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("--- WORD FREQUENCY COUNTER ---")
	
	sentence := "Go is fast and Go is fun because Go is awesome"

	cleanSentence := strings.ToLower(sentence)
	words := strings.Fields(cleanSentence) 

	frequency := make(map[string]int)

	for _, word := range words {
		frequency[word]++ 
	}

	for word, count := range frequency {
		fmt.Printf("Word: '%s' | Count: %d\n", word, count)
	}

}