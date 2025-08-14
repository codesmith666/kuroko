package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	DOT
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NE:       EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
	token.ACCESS:   DOT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	errors     []string
	tokens     []*token.Token
	position   int
	curToken   *token.Token
	peekToken  *token.Token
	peek2Token *token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func NewParser(input string) *Parser {

	p := &Parser{
		errors:   []string{},
		tokens:   lexer.GetTokens(input),
		position: 0,
	}

	// そのトークンが式の先頭に現れる可能性があるならprefixに登録する
	p.prefixParseFns = map[token.TokenType]prefixParseFn{
		token.IDENT:    p.parseIdentifier,
		token.INTEGER:  p.parseIntegerLiteral,
		token.STRING:   p.parseStringLiteral,
		token.NOT:      p.parsePrefixExpression,
		token.MINUS:    p.parsePrefixExpression,
		token.INC:      p.parsePrefixExpression,
		token.PARSE:    p.parsePrefixExpression,
		token.TRUE:     p.parseBoolean,
		token.FALSE:    p.parseBoolean,
		token.LPAREN:   p.parseGroupedExpression, // 関数リテラルもここで
		token.IF:       p.parseIfExpression,
		token.LBRACKET: p.parseArrayLiteral,
		token.LBRACE:   p.parseHashLiteral,
	}

	p.infixParseFns = map[token.TokenType]infixParseFn{
		token.PLUS:     p.parseInfixExpression,
		token.MINUS:    p.parseInfixExpression,
		token.SLASH:    p.parseInfixExpression,
		token.ASTERISK: p.parseInfixExpression,
		token.EQ:       p.parseInfixExpression,
		token.NE:       p.parseInfixExpression,
		token.LT:       p.parseInfixExpression,
		token.GT:       p.parseInfixExpression,
		token.LPAREN:   p.parseCallExpression,
		token.LBRACKET: p.parseIndexExpression,
		token.ACCESS:   p.parseDotExpression,
	}
	// 最初のトークンを準備する
	p.nextToken()
	return p
}

// 指定された位置のトークンを取得する
// positionがp.tokensの位置を越えていたら最後のトークン(EOF)を返す
func (p *Parser) getToken(position int) *token.Token {
	if position < len(p.tokens) {
		return p.tokens[position]
	}
	return p.tokens[len(p.tokens)-1] // 最後を返し続ける
}

// 次のトークンを読み込んで準備する
func (p *Parser) nextToken() {
	p.curToken = p.getToken(p.position)
	p.peekToken = p.getToken(p.position + 1)
	p.peek2Token = p.getToken(p.position + 2)
	p.position++
}

// 現在のトークンタイプは？
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// 次のトークンタイプは？
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// 次の次のトークンタイプは？
func (p *Parser) peek2TokenIs(t token.TokenType) bool {

	return p.peek2Token.Type == t
}

// 次のトークンが想定通りならトークンを進める
// 想定通りでなければエラー
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) DumpTokens() {
	for _, t := range p.tokens {
		fmt.Printf("%q\n", t.String())
	}
}
