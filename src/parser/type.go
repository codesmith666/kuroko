package parser

import (
	"monkey/ast"
	"monkey/token"
)

func (p *Parser) parseTypeAnnotation() *ast.TypeNode {
	tok := p.curToken

	// 単純型: number, string, boolean
	if p.curToken.Type == token.IDENT || p.curToken.Type == token.TYPE {
		name := p.curToken.Literal
		// 型の後ろに[]がついていたら配列
		if p.peekTokenIs(token.LBRACKET) && p.peek2TokenIs(token.RBRACKET) {
			// 例: number[]
			p.nextToken() // consume '['
			p.nextToken() // consume ']'
			return &ast.TypeNode{
				Token: *tok,
				Kind:  ast.TypeArray,
				ElementType: &ast.TypeNode{
					Token: *tok,
					Kind:  ast.TypeSimple,
					Name:  name,
				},
			}
		}
		// そうでなかったら単純型なんだけど、
		// ここにオブジェクトタイプも入る
		return &ast.TypeNode{
			Token: *tok,
			Kind:  ast.TypeSimple,
			Name:  name,
		}
	}

	// 関数型
	if p.curToken.Type == token.LPAREN {
		tok := p.curToken
		p.nextToken()

		params := []*ast.Identifier{}

		for !p.curTokenIs(token.RPAREN) {
			if p.curToken.Type != token.IDENT {
				p.errors = append(p.errors, "expected identifier in function type params")
				return nil
			}

			t := *p.curToken
			if !p.expectPeek(token.COLON) {
				return nil
			}

			p.nextToken() // 型へ
			paramType := p.parseTypeAnnotation()

			params = append(params, &ast.Identifier{Token: t, Name: t.Literal, Type: paramType})

			if p.peekTokenIs(token.COMMA) {
				p.nextToken()
				p.nextToken()
			} else {
				break
			}
		}

		if !p.expectPeek(token.RPAREN) {
			return nil
		}

		if !p.expectPeek(token.ARROW) { // '=>'
			p.errors = append(p.errors, "expected => in function type")
			return nil
		}

		// 戻り値の型へ
		p.nextToken()
		returnType := p.parseTypeAnnotation()

		return &ast.TypeNode{
			Token:      *tok,
			Kind:       ast.TypeFunction,
			Parameters: params,
			ReturnType: returnType,
		}
	}

	// インデックスシグネチャの場合
	if p.curToken.Type == token.LBRACE && p.peekTokenIs(token.LBRACKET) {
		p.nextToken() // '{' をから '['にする
		p.nextToken() // '[' をから keyにする
		p.nextToken() // keyから':'
		p.nextToken() // ':'から型（stringなど）

		// キーのタイプ
		keyType := &ast.TypeNode{
			Token: *p.curToken,
			Kind:  ast.TypeSimple,
			Name:  p.curToken.Literal,
		}

		// RBRACKETにすすめる
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}

		// COLONにすすめる
		if !p.expectPeek(token.COLON) {
			return nil
		}

		// valueの型へ
		p.nextToken()
		valueType := p.parseTypeAnnotation()

		if !p.expectPeek(token.RBRACE) {
			return nil
		}

		return &ast.TypeNode{
			Token:     *tok,
			Kind:      ast.TypeMap,
			KeyType:   keyType,
			ValueType: valueType,
		}
	}

	if p.curToken.Type == token.LBRACE && p.peekTokenIs(token.IDENT) {

		// オブジェクト型:
		props := []*ast.ObjectProperty{}
		p.nextToken() // '{' から IDENTへ

		// }を検出するまで
		for !p.curTokenIs(token.RBRACE) {
			// プロパティ名（識別子）
			if p.curToken.Type != token.IDENT {
				p.errors = append(p.errors, "expected identifier in object type")
				return nil
			}
			propName := p.curToken.Literal

			if !p.expectPeek(token.COLON) {
				return nil
			}

			p.nextToken() // 型へ
			propType := p.parseTypeAnnotation()
			props = append(props, &ast.ObjectProperty{
				Name: propName,
				Type: propType,
			})

			if p.peekTokenIs(token.COMMA) {
				p.nextToken() // COMMNAにする
				p.nextToken() // IDENTにする
			} else {
				break
			}
		}

		if !p.expectPeek(token.RBRACE) {
			return nil
		}

		return &ast.TypeNode{
			Token:      *tok,
			Kind:       ast.TypeObject,
			Properties: props,
		}
	}
	p.peekError(token.LBRACE)
	return nil

}
