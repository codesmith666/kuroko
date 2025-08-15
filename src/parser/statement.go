package parser

import (
	"monkey/ast"
	"monkey/token"
)

// 現在のトークンを調べて処理を行う
// 現在の実装では IMMとRETURNと式しかないがいずれ拡張するだろう
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {

	case token.IMM:
		return p.parseLetStatement()
	case token.MUT:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.BLOCK_COMMENT, token.LINE_COMMENT:
		return p.parseCommentStatement()
	case token.PARSE:
		return p.parseEllipsisStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// コメントは値を返さないので式でなくて文とする
func (p *Parser) parseCommentStatement() *ast.CommentStatement {
	// parseProgramのループの最後で p.nextToken()してるので、
	// 処理が終わった時に最後のトークンを示していればよい。
	// = p.nextToken()は必要ない
	return ast.NewCommentStatement(*p.curToken)
}

// Letステートメント
// LetStatementにはimmとmutの二種類がある
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: *p.curToken}

	// 次が識別子でないとエラー
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// 識別子を取得
	stmt.Ident = &ast.Identifier{Token: *p.curToken, Name: p.curToken.Literal}

	// その次にCOLONがあれば型指定（省略可能）
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // curがCOLONに
		p.nextToken() // 型に入る
		stmt.Ident.Type = p.parseTypeAnnotation()
	}

	// ここで終了してもよい
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}

	// 終了しない場合、次はASSIGN(=)でなければエラー
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// トークンを'='から進めて式の先頭に
	p.nextToken()

	// 式を得る
	stmt.Value = p.parseExpression(LOWEST)

	// セミコロンがあれば読み飛ばす（なくてもいい）
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// 式ステートメント
// *ast.ExpressionStatement
func (p *Parser) parseExpressionStatement() ast.Statement {

	// curTokenは LETでもRETURNでもない式の先頭のトークンが入っている
	stmt := &ast.ExpressionStatement{Token: *p.curToken}

	// 式の解析を開始する
	//   式に対する代入文があるときと
	//   関数呼び出しなどのただの式の実行で処理をわける
	expr := p.parseExpression(LOWEST)
	if p.peekTokenIs(token.ASSIGN) {
		return p.parseAssignStatement(expr)
	} else {
		stmt.Expression = expr

		// セミコロンがあれば飛ばす（省略可能）
		if p.peekTokenIs(token.SEMICOLON) {
			p.nextToken()
		}
	}
	return stmt
}

// RETURNステートメント
// 値が必ずあることを前提に処理が書いてある
// なければ？
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	// リターンステートメントを準備
	stmt := &ast.ReturnStatement{Token: *p.curToken}

	// 式を取得する
	p.nextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	// セミコロンがあれば飛ばす（なくてもエラーにならない）
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: *p.curToken}
	block.Statements = []ast.Statement{}

	// ステートメントを一覧に追加
	// このとき return this を見つけたら
	// コンストラクタを示すフラグも返す。
	addStatements := func(stmt ast.Statement) {
		block.Statements = append(block.Statements, stmt)
	}

	// {～}をブロックとして取得するか、１分だけ取得するか
	if p.curTokenIs(token.LBRACE) {
		p.nextToken()
		for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
			stmt := p.parseStatement()
			addStatements(stmt)
			p.nextToken()
		}
		// この時点でcurはRBRACEだけとp.NexeToken()しない
	} else {
		stmt := p.parseExpressionStatement()
		addStatements(stmt)
	}
	// コンストラクタを示すフラグをセット
	return block
}

/*
 *	代入
 */
func (p *Parser) parseAssignStatement(left ast.Expression) *ast.AssignStatement {
	stmt := &ast.AssignStatement{Token: *p.peekToken, Left: left}

	// '=' へ進む
	p.nextToken()

	// '=' の次へ
	p.nextToken()

	stmt.Right = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseEllipsisStatement() ast.Statement {
	stmt := &ast.DeriveStatement{Token: *p.curToken}
	p.nextToken()
	stmt.Right = p.parseExpression(LOWEST)
	return stmt
}
