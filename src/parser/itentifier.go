package parser

import (
	"monkey/ast"
)

// 識別子を抽象構文木に設定して返す
func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: *p.curToken, Name: p.curToken.Literal}
	return ident
}

// func (p *Parser) parseFunctionParameters() []*ast.Identifier {
// 	identifiers := []*ast.Identifier{}

// 	//	関数パラメータなし
// 	if p.peekTokenIs(token.RPAREN) {
// 		p.nextToken()
// 		return identifiers
// 	}

// 	// 最初の関数パラメータ
// 	p.nextToken()

// 	ident := &ast.Identifier{Token: *p.curToken, Name: p.curToken.Literal}
// 	identifiers = append(identifiers, ident)

// 	for p.peekTokenIs(token.COMMA) {
// 		p.nextToken()
// 		p.nextToken()
// 		ident := &ast.Identifier{Token: *p.curToken, Name: p.curToken.Literal}
// 		identifiers = append(identifiers, ident)
// 	}

// 	if !p.expectPeek(token.RPAREN) {
// 		return nil
// 	}

// 	return identifiers
// }
