package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type TokenType string

const (
	LeftParen  TokenType = "("
	RightParen TokenType = ")"
	LeftBrace  TokenType = "{"
	RightBrace TokenType = "}"
	Star       TokenType = "*"
	Comma      TokenType = ","
	Plus       TokenType = "+"
	Dot        TokenType = "."
	Minus      TokenType = "-"
	Semicolon  TokenType = ";"
)

var tokenNames = map[TokenType]string{
	LeftParen:  "LEFT_PAREN",
	RightParen: "RIGHT_PAREN",
	LeftBrace:  "LEFT_BRACE",
	RightBrace: "RIGHT_BRACE",
	Star:       "STAR",
	Dot:        "DOT",
	Comma:      "COMMA",
	Plus:       "PLUS",
	Minus:      "MINUS",
	Semicolon:  "SEMICOLON",
}

type Token struct {
	tokenType TokenType
	line      int
	lexeme    string
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s null", tokenNames[t.tokenType], t.lexeme)
}

func generateToken(tokenType TokenType, line int) Token {
	return Token{tokenType, line, string(tokenType)}
}

// TODO: handle more kind of errors
func reportError(line int, input string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", line, input)
}

func scan(reader *bufio.Reader) {
	hasErrors := false
	tokens := make([]Token, 0)
	for i := 1; ; {
		bytes, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("Error reading line: %v", err)
		}

		for _, current := range bytes {
			switch current {
			case '(':
				token := generateToken(LeftParen, i)
				tokens = append(tokens, token)
			case ')':
				token := generateToken(RightParen, i)
				tokens = append(tokens, token)
			case '{':
				token := generateToken(LeftBrace, i)
				tokens = append(tokens, token)
			case '}':
				token := generateToken(RightBrace, i)
				tokens = append(tokens, token)
			case '*':
				token := generateToken(Star, i)
				tokens = append(tokens, token)
			case '.':
				token := generateToken(Dot, i)
				tokens = append(tokens, token)
			case ',':
				token := generateToken(Comma, i)
				tokens = append(tokens, token)
			case '+':
				token := generateToken(Plus, i)
				tokens = append(tokens, token)
			case '-':
				token := generateToken(Minus, i)
				tokens = append(tokens, token)
			case ';':
				token := generateToken(Semicolon, i)
				tokens = append(tokens, token)
			default:
				reportError(i, string(current))
				hasErrors = true
				continue
			}

			if len(tokens) > 0 {
				fmt.Println(tokens[len(tokens)-1].String())
			}
		}

		if err == io.EOF {
			fmt.Println("EOF  null")
			break
		}

		// Next line
		i++
	}

	if hasErrors {
		os.Exit(65)
	}
}
