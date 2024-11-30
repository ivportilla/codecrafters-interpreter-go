package main

import (
	"errors"
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
type Grouping struct {
	Value Expr
}
type Unary struct {
	Operator   Token
	Expression Expr
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

func NewGrouping(expr Expr) Expr {
	return &Grouping{expr}
}

func NewUnary(op Token, exp Expr) Expr {
	return &Unary{op, exp}
}

func (boolExpr *Boolean) Print() string {
	return when(boolExpr.Value, "true", "false")
}

func (nilExpr *Nil) Print() string {
	return "nil"
}

func (numberExpr *NumberLit) Print() string { return formatFloatNumber(numberExpr.Value) }

func (stringExpr *StringLit) Print() string { return stringExpr.Value }

func (grouping *Grouping) Print() string { return "(group " + grouping.Value.Print() + ")" }

func (unary *Unary) Print() string {
	return fmt.Sprintf("(%s %s)", unary.Operator.lexeme, unary.Expression.Print())
}

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

func (p *Parser) check(tokenType TokenType) bool {
	return p.tokens[p.current].tokenType == tokenType
}

func (p *Parser) match(tokenType TokenType) bool {
	if p.check(tokenType) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) nextToken() Token {
	p.advance()
	return p.tokens[p.current]
}

func (p *Parser) previousToken() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previousToken()
}

func (p *Parser) consume(tokenType TokenType, errorMsg string) error {
	if !p.match(tokenType) {
		return errors.New(errorMsg)
	}
	return nil
}

func (p *Parser) isAtEnd() bool {
	return p.currentToken().tokenType == EOF
}

func (p *Parser) MatchUnary() (Expr, error) {
	if p.match(Bang) || p.match(Minus) {
		op := p.previousToken()
		expr, err := p.MatchUnary()
		if err != nil {
			return nil, err
		}
		res := NewUnary(op, expr)
		return res, nil
	} else {
		return p.MatchPrimary()
	}
}

func (p *Parser) MatchPrimary() (Expr, error) {
	if p.match(LeftParen) {
		expr, err := p.MatchExpr()
		if err != nil {
			return nil, err
		}
		err = p.consume(RightParen, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewGrouping(expr), nil
	} else {
		lit, err := NewLiteral(p.currentToken())
		p.advance()
		return lit, err
	}
}

func (p *Parser) MatchExpr() (Expr, error) {
	return p.MatchUnary()
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
