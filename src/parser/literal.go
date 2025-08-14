package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

// 二値
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: *p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// 整数値
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: *p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// 文字列
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: *p.curToken, Value: p.curToken.Literal}
}

// 配列
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: *p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

// ハッシュリテラル
func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: *p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBRACE) {
		p.nextToken()

		var key ast.Expression
		if p.curToken.Type == token.IDENT {
			key = &ast.StringLiteral{
				Token: token.Token{Type: token.STRING, Literal: p.curToken.Literal},
				Value: p.curToken.Literal,
			}
		} else {
			key = p.parseExpression(LOWEST)
		}

		// COLONを想定してNextする
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken() // ":"の次へ移動する
		// ここから式
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return hash
}

// 関数リテラル
// func (p *Parser) parseFunctionLiteral() ast.Expression {
// 	lit := &ast.FunctionLiteral{Token: *p.curToken}
// 	if !p.expectPeek(token.LPAREN) {
// 		return nil
// 	}
// 	lit.Parameters = p.parseFunctionParameters()
// 	if !p.expectPeek(token.LBRACE) {
// 		return nil
// 	}
// 	lit.Body = p.parseBlockStatement()
// 	return lit
// }

func (p *Parser) parseArrowFunctionLiteral() ast.Expression {

	lit := &ast.FunctionLiteral{Token: *p.curToken}

	// 関数は '(' から始まるというか、グループ化
	if p.curToken.Type != token.LPAREN {
		p.peekError(token.LPAREN)
		return nil
	}

	// 関数の引数の名前と型の定義
	params := []*ast.Identifier{}
	for !p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		if p.curToken.Type != token.IDENT {
			return nil
		}

		param := &ast.Identifier{Token: *p.curToken, Name: p.curToken.Literal}

		if p.peekTokenIs(token.COLON) {
			p.nextToken() // consume colon
			p.nextToken() // consume type
			param.Type = p.parseTypeAnnotation()
		}

		params = append(params, param)

		if p.peekTokenIs(token.COMMA) {
			p.nextToken()
		} else {
			break
		}
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// 戻り値の型があれば
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // ":"
		p.nextToken() // "typeの先頭"
		lit.ReturnType = p.parseTypeAnnotation()
	}

	// 矢印じゃないとダメ
	if !p.expectPeek(token.ARROW) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Parameters = params
	lit.Body = p.parseBlockStatement()

	return lit
}
