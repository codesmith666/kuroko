package parser

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"monkey/token"
	"reflect"
)

// プログラムのパースを開始する
// プログラムは、Statementの羅列である
func (p *Parser) ParseProgram() (*ast.Program, bool) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if !reflect.ValueOf(stmt).IsNil() {
			program.Statements = append(program.Statements, stmt)
		}
		// ブロックステートメントがRPARENで終了するのは、
		// ここで必ず読み飛しが１つ入るから？
		p.nextToken()
	}

	// エラー
	if len(p.errors) > 0 {
		return program, false
	}

	return program, true
}

func (p *Parser) OutputErrors() {
	var out bytes.Buffer
	out.WriteString("Woops! We ran into some monkey business here!\n")
	out.WriteString(" parser errors:\n")
	for _, msg := range p.errors {
		out.WriteString("\t" + msg + "\n")
	}
	fmt.Println(out.String())
}
