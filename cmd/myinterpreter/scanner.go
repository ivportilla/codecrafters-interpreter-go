package main

import "fmt"

type TokenType string

const (
	LeftParen  TokenType = "("
	RightParen TokenType = ")"
	LeftBrace  TokenType = "{"
	RightBrace TokenType = "}"
)

var tokenNames = map[TokenType]string{
	LeftParen:  "LEFT_PAREN",
	RightParen: "RIGHT_PAREN",
	LeftBrace:  "LEFT_BRACE",
	RightBrace: "RIGHT_BRACE",
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

func scan(fileContents []byte) {
	tokens := make([]Token, 0)
	for i, char := range fileContents {
		switch char {
		case '(':
			token := Token{tokenType: LeftParen, line: i, lexeme: string(LeftParen)}
			tokens = append(tokens, token)
		case ')':
			token := Token{tokenType: RightParen, line: i, lexeme: string(RightParen)}
			tokens = append(tokens, token)
		case '{':
			token := generateToken(LeftBrace, i)
			tokens = append(tokens, token)
		case '}':
			token := generateToken(RightBrace, i)
			tokens = append(tokens, token)
		}
		fmt.Println(tokens[len(tokens)-1].String())
	}
	fmt.Println("EOF  null")
}
