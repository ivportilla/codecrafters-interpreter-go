package main

import "fmt"

type TokenType string

const (
	LeftParen  TokenType = "("
	RightParen TokenType = ")"
)

type Token struct {
	tokenType TokenType
	line      int
	lexeme    string
}

func (t Token) String() string {
	switch t.tokenType {
	case LeftParen:
		return fmt.Sprintf("LEFT_PAREN %s null", t.tokenType)
	case RightParen:
		return fmt.Sprintf("RIGHT_PAREN %s null", t.tokenType)
	default:
		return ""
	}
}

func scan(fileContents []byte) {
	tokens := make([]Token, 0)
	for i, char := range fileContents {
		switch char {
		case '(':
			token := Token{tokenType: LeftParen, line: i, lexeme: string(LeftParen)}
			tokens = append(tokens, token)
			fmt.Println(token.String())
		case ')':
			token := Token{tokenType: RightParen, line: i, lexeme: string(RightParen)}
			tokens = append(tokens, token)
			fmt.Println(token.String())
		}
	}
	fmt.Println("EOF  null")
}
