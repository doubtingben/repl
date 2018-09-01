package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	input := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(input)
	scanner.Split(cliSplitFunc)
	printHeader()
	for scanner.Scan() {
		switch text := strings.TrimSuffix(scanner.Text(), "\n"); text {
		case "help":
			fmt.Printf("here's your help: %s", text)
		case "noop":
			fmt.Printf("Here's noop: %s", text)
		default:
			fmt.Printf("wtf, %s", text)
		}
		fmt.Print("\ncli > ")
	}

}

func printHeader() {
	fmt.Println("Welcome to the cli")
	fmt.Print("cli > ")
}

func cliSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if string(data[len(data)-1]) == "\n" {
		return len(data), data, nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return
}
