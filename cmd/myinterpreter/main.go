package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func tokenizeFile(filename string) ([]Token, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	reader := bufio.NewReader(file)
	if data, _ := reader.Peek(1); len(data) > 0 {
		tokens, err := scan(reader)
		return tokens, err
	} else {
		return []Token{generateEOFToken(0)}, nil
	}
}

func handleCommand(command string, params ...string) {
	switch command {
	case "tokenize":
		tokens, err := tokenizeFile(params[0])
		if err != nil {
			if errors.Is(err, TokenScanError) {
				for _, token := range tokens {
					fmt.Println(token.String())
				}
				os.Exit(65)
			}
			os.Exit(1)
		}

		for _, token := range tokens {
			fmt.Println(token.String())
		}
	case "parse":
		tokens, err := tokenizeFile(params[0])
		if err != nil {
			if errors.Is(err, TokenScanError) {
				os.Exit(65)
			}
			os.Exit(1)
		}
		parse(tokens)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh <COMMAND> <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	handleCommand(command, os.Args[2:]...)
}
