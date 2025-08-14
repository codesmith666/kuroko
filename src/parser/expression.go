package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
)

// 式の検出
func (p *Parser) parseExpression(precedence int) ast.Expression {

	// まず式の先頭から始められるトークンか調べ左辺を初期化する
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	// ここは Prattマジックなのであとでしっかり見る
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// NOT とか MINUSとか前置演算子
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    *p.curToken,
		Operator: p.curToken.Literal, // NOTとかMINUSとか
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// 二項演算子
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    *p.curToken,
		Operator: p.curToken.Literal, // 四則演算とか
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// インデックス参照式
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: *p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

// IF式
func (p *Parser) parseIfExpression() ast.Expression {
	// curはif
	expression := &ast.IfExpression{Token: *p.curToken}

	// 次がifの式の開き括弧をだったらcurを開き括弧にする
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	// ifの式の閉じかっこを検出
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	p.nextToken() // curをtoken.RPARENから１つ先にすすめる

	// ブロックが始まるか式があるはず
	expression.Consequence = p.parseBlockStatement()

	// elseがなかったらここで終了
	if !p.peekTokenIs(token.ELSE) {
		return expression
	}
	p.nextToken() // curは token.ELSE になる
	p.nextToken() // curは token.LBRACEか式になるはず

	// ブロックが始まるか式があるはず
	expression.Alternative = p.parseBlockStatement()

	return expression
}

// グループ化された式
// 関数定義を検出するし
func (p *Parser) parseGroupedExpression() ast.Expression {

	// peekToken が IDENT かつ次に ':' または ',' があればアロー関数とみなす
	if (p.peekTokenIs(token.IDENT) && p.peek2TokenIs(token.COLON)) || p.peekTokenIs(token.RPAREN) {
		return p.parseArrowFunctionLiteral()
	}

	p.nextToken()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// リスト式？これなに？
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// 関数呼び出し
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: *p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

// DOTによる呼び出し
func (p *Parser) parseDotExpression(left ast.Expression) ast.Expression {
	t := *p.curToken // DOT
	p.nextToken()    // 名前へ

	if p.curToken.Type != token.IDENT {
		msg := fmt.Sprintf("expected property name after '.', got %s", p.curToken.Type)
		p.errors = append(p.errors, msg)
		return nil
	}

	ident := &ast.Identifier{
		Token: *p.curToken,
		Name:  p.curToken.Literal,
	}
	return &ast.DotExpression{
		Token: t,
		Left:  left,
		Right: ident,
	}
}
