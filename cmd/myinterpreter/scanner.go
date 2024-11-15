package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type TokenType string

const (
	LeftParen    TokenType = "("
	RightParen   TokenType = ")"
	LeftBrace    TokenType = "{"
	RightBrace   TokenType = "}"
	Star         TokenType = "*"
	Comma        TokenType = ","
	Plus         TokenType = "+"
	Dot          TokenType = "."
	Minus        TokenType = "-"
	Semicolon    TokenType = ";"
	Equal        TokenType = "="
	EqualEqual   TokenType = "=="
	Bang         TokenType = "!"
	BangEqual    TokenType = "!="
	Less         TokenType = "<"
	LessEqual    TokenType = "<="
	Greater      TokenType = ">"
	GreaterEqual TokenType = ">="
	Slash        TokenType = "/"
	String       TokenType = "STR"
)

var tokenNames = map[TokenType]string{
	LeftParen:    "LEFT_PAREN",
	RightParen:   "RIGHT_PAREN",
	LeftBrace:    "LEFT_BRACE",
	RightBrace:   "RIGHT_BRACE",
	Star:         "STAR",
	Dot:          "DOT",
	Comma:        "COMMA",
	Plus:         "PLUS",
	Minus:        "MINUS",
	Semicolon:    "SEMICOLON",
	Equal:        "EQUAL",
	EqualEqual:   "EQUAL_EQUAL",
	Bang:         "BANG",
	BangEqual:    "BANG_EQUAL",
	Less:         "LESS",
	LessEqual:    "LESS_EQUAL",
	Greater:      "GREATER",
	GreaterEqual: "GREATER_EQUAL",
	Slash:        "SLASH",
	String:       "STRING",
}

type Token struct {
	tokenType TokenType
	line      int
	lexeme    string
	literal   string
}

func when[A any](cond bool, ok A, otherwise A) A {
	if cond {
		return ok
	} else {
		return otherwise
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", tokenNames[t.tokenType], t.lexeme, when(t.literal == "", "null", t.literal))
}

func generateStrToken(line int, literal string) Token {
	return Token{String, line, literal, strings.ReplaceAll(literal, `"`, "")}
}

func generateToken(tokenType TokenType, line int) Token {
	return Token{tokenType, line, string(tokenType), ""}
}

// TODO: handle more kind of errors
func reportError(line int, error string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, error)
}

var UnexpectedTokenError = errors.New("unexpected token")
var UnterminatedStringError = errors.New("unterminated string")

func getTokenByType(line []byte, lineNumber int, col int, target TokenType) (Token, error) {
	for i := 0; i < len(target); i++ {
		if col+i >= len(line) {
			return Token{}, UnexpectedTokenError
		}

		if line[col+i] != target[i] {
			return Token{}, UnexpectedTokenError
		}
	}

	return generateToken(target, lineNumber), nil
}

func matchNextChar(line []byte, col int, target byte) bool {
	if col+1 >= len(line) {
		return false
	}

	return line[col+1] == target
}

func countSkipLineComment(line []byte, col int) int {
	return len(line) - col
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func isComment(line []byte, col int) bool {
	return line[col] == '/' && matchNextChar(line, col, '/')
}

func getStringLiteral(line []byte, col int) (string, int, error) {
	builder := strings.Builder{}
	builder.WriteByte('"')
	i := col + 1
	for {
		if i >= len(line) || line[i] == '\n' {
			return "", i - col + 1, UnterminatedStringError
		}

		if line[i] == '"' {
			builder.WriteByte('"')
			return builder.String(), i - col + 1, nil
		}

		builder.WriteByte(line[i])
		i++
	}
}

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
		token, err := getTokenByType(line, lineNumber, col, EqualEqual)
		if err != nil {
			return generateToken(Equal, lineNumber), 1, nil
		}
		return token, len(token.lexeme), nil
	case '!':
		token, err := getTokenByType(line, lineNumber, col, BangEqual)
		if err != nil {
			return generateToken(Bang, lineNumber), 1, nil
		}
		return token, len(token.lexeme), nil
	case '<':
		token, err := getTokenByType(line, lineNumber, col, LessEqual)
		if err != nil {
			return generateToken(Less, lineNumber), 1, nil
		}
		return token, len(token.lexeme), nil
	case '>':
		token, err := getTokenByType(line, lineNumber, col, GreaterEqual)
		if err != nil {
			return generateToken(Greater, lineNumber), 1, nil
		}
		return token, len(token.lexeme), nil
	case '/':
		return generateToken(Slash, lineNumber), 1, nil
	case '"':
		str, count, err := getStringLiteral(line, col)
		if err != nil {
			return Token{}, count, err
		}
		return generateStrToken(lineNumber, str), count, nil
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
			// Handle line comments
			if isComment(line, col) {
				col += countSkipLineComment(line, col)
				continue
			}

			// Handle spaces
			if isSpace(line[col]) {
				col++
				continue
			}

			token, count, errToken := getToken(line, lineNumber, col)
			if errToken != nil {
				if errors.Is(errToken, UnexpectedTokenError) {
					reportError(lineNumber, fmt.Sprintf("Unexpected character: %s", string(line[col])))
					hasErrors = true
					col += count
					continue
				}

				if errors.Is(errToken, UnterminatedStringError) {
					reportError(lineNumber, "Unterminated string.")
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
