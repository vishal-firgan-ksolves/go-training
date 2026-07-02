package main

import (
	"fmt"
	"bufio"
	"os"
)

func OpenFile(path string){

}

func main(){
	fmt.Println("Started program.......")

	filePath := "./users.csv"

	file, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Error while opening file:", err)
		return
	}

	defer func() {
		file.Close()
		fmt.Println("OS: FILE CLOSED SAFELY")
	}()

	fmt.Printf("OS: FILE OPENED SUCCESSFULLY -> %s\n\n", filePath)

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
		fmt.Printf("Read Line %d: %s\n", lineCount, scanner.Text())

		if lineCount == 3 {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	fmt.Println("\nReaching the end of the main function.")
}