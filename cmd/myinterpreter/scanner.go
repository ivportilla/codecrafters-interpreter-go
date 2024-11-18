package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"
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
	Number       TokenType = "NUM"
	Identifier   TokenType = "ID"
	Keyword      TokenType = "KEYWORD"
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
	Number:       "NUMBER",
	Identifier:   "IDENTIFIER",
	Keyword:      "KEYWORD",
}

var keywords = map[string]interface{}{
	"and":    struct{}{},
	"class":  struct{}{},
	"else":   struct{}{},
	"false":  struct{}{},
	"for":    struct{}{},
	"fun":    struct{}{},
	"if":     struct{}{},
	"nil":    struct{}{},
	"or":     struct{}{},
	"print":  struct{}{},
	"return": struct{}{},
	"super":  struct{}{},
	"this":   struct{}{},
	"true":   struct{}{},
	"var":    struct{}{},
	"while":  struct{}{},
}

type Token struct {
	tokenType TokenType
	line      int
	lexeme    string
	literal   any
}

func when[A any](cond bool, ok A, otherwise A) A {
	if cond {
		return ok
	} else {
		return otherwise
	}
}

func (t Token) String() string {
	switch t.tokenType {
	case Number:
		formatted := strconv.FormatFloat(t.literal.(float64), 'f', -1, 64)
		if !strings.Contains(formatted, ".") {
			formatted += ".0"
		}
		return fmt.Sprintf("%s %s %s", tokenNames[t.tokenType], t.lexeme, formatted)
	case Identifier:
		return fmt.Sprintf("%s %s %s", tokenNames[t.tokenType], t.lexeme, when(t.literal == nil, "null", t.literal))
	case Keyword:
		return fmt.Sprintf("%s %s %s", strings.ToUpper(t.lexeme), t.lexeme, when(t.literal == nil, "null", t.literal))
	default:
		return fmt.Sprintf("%s %s %s", tokenNames[t.tokenType], t.lexeme, when(t.literal == nil, "null", t.literal))
	}
}

func generateStrToken(line int, literal string) Token {
	return Token{String, line, literal, strings.ReplaceAll(literal, `"`, "")}
}

func generateNumberToken(line int, literal float64, lexeme string) Token {
	return Token{Number, line, lexeme, literal}
}

func generateIdentifierToken(line int, lexeme string) Token {
	return Token{Identifier, line, lexeme, nil}
}

func generateKeywordToken(line int, lexeme string) Token {
	return Token{Keyword, line, lexeme, nil}
}

func generateToken(tokenType TokenType, line int) Token {
	return Token{tokenType, line, string(tokenType), nil}
}

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

	for i := col + 1; ; i++ {
		if i >= len(line) || line[i] == '\n' {
			return "", i - col + 1, UnterminatedStringError
		}

		if line[i] == '"' {
			builder.WriteByte('"')
			return builder.String(), i - col + 1, nil
		}

		builder.WriteByte(line[i])
	}
}

func getNumberLiteral(line []byte, col int) (float64, string, int, error) {
	rawResult := ""
	i := col
	func() {
		for ; i < len(line); i++ {
			switch {
			case unicode.IsDigit(rune(line[i])):
				rawResult += string(line[i])
			case line[i] == '.' && !strings.Contains(rawResult, "."):
				rawResult += string(line[i])
			default:
				return
			}
		}
	}()

	result, err := strconv.ParseFloat(rawResult, 64)
	if err != nil {
		return 0.0, rawResult, i - col, errors.New("error converting target to number")
	}
	return result, rawResult, i - col, nil
}

func getIdentifier(line []byte, col int) (string, int) {
	i := col
	builder := strings.Builder{}
	func() {
		for ; i < len(line); i++ {
			switch {
			case unicode.IsDigit(rune(line[i])):
				builder.WriteByte(line[i])
			case line[i] == '_':
				builder.WriteByte(line[i])
			case unicode.IsLetter(rune(line[i])):
				builder.WriteByte(line[i])
			default:
				return
			}
		}
	}()

	return builder.String(), i - col
}

func getToken(line []byte, lineNumber int, col int) (Token, int, error) {
	switch {
	case line[col] == '(':
		token := generateToken(LeftParen, lineNumber)
		return token, 1, nil
	case line[col] == ')':
		token := generateToken(RightParen, lineNumber)
		return token, 1, nil
	case line[col] == '{':
		token := generateToken(LeftBrace, lineNumber)
		return token, 1, nil
	case line[col] == '}':
		token := generateToken(RightBrace, lineNumber)
		return token, 1, nil
	case line[col] == '*':
		token := generateToken(Star, lineNumber)
		return token, 1, nil
	case line[col] == '.':
		token := generateToken(Dot, lineNumber)
		return token, 1, nil
	case line[col] == ',':
		token := generateToken(Comma, lineNumber)
		return token, 1, nil
	case line[col] == '+':
		token := generateToken(Plus, lineNumber)
		return token, 1, nil
	case line[col] == '-':
		token := generateToken(Minus, lineNumber)
		return token, 1, nil
	case line[col] == ';':
		token := generateToken(Semicolon, lineNumber)
		return token, 1, nil
	case line[col] == '=':
		token, err := getTokenByType(line, lineNumber, col, EqualEqual)
		if err != nil {
			return generateToken(Equal, lineNumber), 1, nil
		}
		return token, len(token.lexeme), nil
	case line[col] == '!':
		token, err := getTokenByType(line, lineNumber, col, BangEqual)
		if err != nil {
			return generateToken(Bang, lineNumber), 1, nil
		}
		return token, len(token.lexeme), nil
	case line[col] == '<':
		token, err := getTokenByType(line, lineNumber, col, LessEqual)
		if err != nil {
			return generateToken(Less, lineNumber), 1, nil
		}
		return token, len(token.lexeme), nil
	case line[col] == '>':
		token, err := getTokenByType(line, lineNumber, col, GreaterEqual)
		if err != nil {
			return generateToken(Greater, lineNumber), 1, nil
		}
		return token, len(token.lexeme), nil
	case line[col] == '/':
		return generateToken(Slash, lineNumber), 1, nil
	case line[col] == '"':
		str, count, err := getStringLiteral(line, col)
		if err != nil {
			return Token{}, count, err
		}
		return generateStrToken(lineNumber, str), count, nil
	case unicode.IsDigit(rune(line[col])):
		number, lexeme, count, err := getNumberLiteral(line, col)
		if err != nil {
			return Token{}, count, err
		}
		return generateNumberToken(lineNumber, number, lexeme), count, nil
	case unicode.IsLetter(rune(line[col])) || line[col] == '_':
		target, count := getIdentifier(line, col)
		if _, isKeyword := keywords[target]; isKeyword {
			return generateKeywordToken(lineNumber, target), count, nil
		}

		return generateIdentifierToken(lineNumber, target), count, nil
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
