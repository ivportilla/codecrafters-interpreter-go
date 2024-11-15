package main

import (
	"bufio"
	"errors"
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
	Equal      TokenType = "="
	EqualEqual TokenType = "=="
	Bang       TokenType = "!"
	BangEqual  TokenType = "!="
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
	Equal:      "EQUAL",
	EqualEqual: "EQUAL_EQUAL",
	Bang:       "BANG",
	BangEqual:  "BANG_EQUAL",
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

func matchNext(line []byte, pos int, char byte) bool {
	if pos+1 >= len(line) {
		return false
	}

	return line[pos+1] == char
}

var UnexpectedTokenError = errors.New("unexpected token")

func getToken(line []byte, lineNumber int, col int) (Token, int, error) {
	switch line[col] {
	case '(':
		token := generateToken(LeftParen, lineNumber)
		return token, 1, nil
	case ')':
		token := generateToken(RightParen, lineNumber)
		return token, 1, nil
	case '{':
		token := generateToken(LeftBrace, lineNumber)
		return token, 1, nil
	case '}':
		token := generateToken(RightBrace, lineNumber)
		return token, 1, nil
	case '*':
		token := generateToken(Star, lineNumber)
		return token, 1, nil
	case '.':
		token := generateToken(Dot, lineNumber)
		return token, 1, nil
	case ',':
		token := generateToken(Comma, lineNumber)
		return token, 1, nil
	case '+':
		token := generateToken(Plus, lineNumber)
		return token, 1, nil
	case '-':
		token := generateToken(Minus, lineNumber)
		return token, 1, nil
	case ';':
		token := generateToken(Semicolon, lineNumber)
		return token, 1, nil
	case '=':
		var token Token
		if matchNext(line, col, '=') {
			token = generateToken(EqualEqual, lineNumber)
			return token, 2, nil
		}
		token = generateToken(Equal, lineNumber)
		return token, 1, nil
	case '!':
		var token Token
		if matchNext(line, col, '=') {
			token = generateToken(BangEqual, lineNumber)
			return token, 2, nil
		}
		token = generateToken(Bang, lineNumber)
		return token, 1, nil
	default:
		return Token{}, 1, UnexpectedTokenError
	}
}

func scan(reader *bufio.Reader) {
	hasErrors := false
	tokens := make([]Token, 0)
	for lineNumber := 1; ; {
		line, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("Error reading line: %v", err)
		}

		for col := 0; col < len(line); {
			token, count, errToken := getToken(line, lineNumber, col)
			if errToken != nil {
				if errors.Is(errToken, UnexpectedTokenError) {
					reportError(lineNumber, string(line[col]))
					hasErrors = true
					col += count
					continue
				}

				log.Fatalf("Unexpected error: %v", errToken)
			}

			fmt.Println(token.String())
			tokens = append(tokens, token)
			col += count
		}

		// Check if EOF
		if err == io.EOF {
			fmt.Println("EOF  null")
			break
		}

		// Next line
		lineNumber++
	}

	if hasErrors {
		os.Exit(65)
	}
}
