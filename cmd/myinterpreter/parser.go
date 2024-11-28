package main

import (
	"fmt"
	"log"
)

type Expr interface {
	Print() string
}
type Boolean struct {
	Value bool
}
type NumberLit struct {
	Value float64
}
type StringLit struct {
	Value string
}
type Nil struct{}

func NewNil() Expr {
	return &Nil{}
}

func NewBoolean(value bool) Expr {
	return &Boolean{value}
}

func NewNumberLit(value float64) Expr { return &NumberLit{value} }

func NewStringLit(value string) Expr { return &StringLit{value} }

func NewLiteral(token Token) (Expr, error) {
	switch token.tokenType {
	case Keyword:
		switch token.lexeme {
		case "true":
			return NewBoolean(true), nil
		case "false":
			return NewBoolean(false), nil
		case "nil":
			return NewNil(), nil
		default:
			return nil, fmt.Errorf("unsupported keyword type: %s", token.lexeme)
		}
	case Number:
		return NewNumberLit(token.literal.(float64)), nil
	case String:
		return NewStringLit(token.literal.(string)), nil
	default:
		return nil, fmt.Errorf("unsupported token type: %s", token.lexeme)
	}
}

func (boolExpr *Boolean) Print() string {
	return when(boolExpr.Value, "true", "false")
}

func (nilExpr *Nil) Print() string {
	return "nil"
}

func (numberExpr *NumberLit) Print() string { return formatFloatNumber(numberExpr.Value) }

func (stringExpr *StringLit) Print() string { return stringExpr.Value }

func printAST(expr Expr) string {
	return expr.Print()
}

type Parser struct {
	tokens  []Token
	current int
}

func (p *Parser) currentToken() Token {
	return p.tokens[p.current]
}

func (p *Parser) MatchLiteral() (Expr, error) {
	return NewLiteral(p.currentToken())
}

func (p *Parser) MatchExpr() (Expr, error) {
	return p.MatchLiteral()
}

func parse(tokens []Token) {
	parser := Parser{tokens: tokens, current: 0}
	expr, err := parser.MatchExpr()
	if err != nil {
		log.Fatal(err)
	}
	astText := printAST(expr)
	fmt.Println(astText)
}
