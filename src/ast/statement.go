package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"regexp"
	"strings"
)

// イミュータブル変数の宣言ステートメント
type LetStatement struct {
	Token token.Token // the token.LET token
	Ident *Identifier // 識別子の名前
	Value Expression  // 式
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String(depth int) string {
	var out bytes.Buffer
	out.WriteString(Indent(depth))
	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Ident.String(depth + 1))
	out.WriteString(" = ")
	// 値
	if ls.Value != nil {
		out.WriteString(ls.Value.String(depth + 1))
	} else {
		out.WriteString("<?>")
	}
	out.WriteString(";\n")
	return out.String()
}

// 式ステートメント
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String(depth int) string {
	if es.Expression != nil {
		return es.Expression.String(depth)
	}
	return ""
}

// リターンステートメント
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String(depth int) string {
	var out bytes.Buffer
	if rs.ReturnValue != nil {
		out.WriteString(Indent(depth))
		out.WriteString("return ")
		out.WriteString(rs.ReturnValue.String(depth + 1))
		out.WriteString(";\n")
	}
	return out.String()
}

// ブロックステートメント
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String(depth int) string {
	var out bytes.Buffer
	out.WriteString("{\n")
	for _, s := range bs.Statements {
		out.WriteString(s.String(depth)) // +1しない
	}
	out.WriteString(Indent(depth-1) + "}")
	return out.String()
}

// コメントステートメント
// トークンリテラルに生コメントが入っている
type CommentStatement struct {
	Token    token.Token // '//' or '/*'
	Comments []string
}

func NewCommentStatement(t token.Token) *CommentStatement {

	var lines []string
	var prefix = 0
	if t.Type == token.BLOCK_COMMENT {
		re := regexp.MustCompile(`\s*\n\s*`)
		lines = re.Split(strings.TrimSpace(t.Literal), -1)
		if len(lines) > 0 {
			re2 := regexp.MustCompile(`^[#*=-]\s+`)
			if match := re2.FindString(lines[0]); match != "" {
				prefix = len(match)
			}
		}
	} else {
		re := regexp.MustCompile(`\s*\n`)
		lines = re.Split(t.Literal, -1)
		if len(lines) > 0 {
			re2 := regexp.MustCompile(`^\s+`)
			if match := re2.FindString(lines[0]); match != "" {
				prefix = len(match)
			}
		}
	}

	comments := []string{}
	if prefix > 0 {
		for _, line := range lines {
			comments = append(comments, line[prefix:])
		}
	} else {
		comments = lines
	}
	return &CommentStatement{Token: t, Comments: comments}
}

func (fl *CommentStatement) statementNode()       {}
func (fl *CommentStatement) TokenLiteral() string { return fl.Token.Literal }

func (fl *CommentStatement) String(depth int) string {
	indent := Indent(depth)

	var out bytes.Buffer
	out.WriteString("\n")
	if fl.Token.Type == token.BLOCK_COMMENT {
		out.WriteString(indent + "/*\n")
		for _, line := range fl.Comments {
			out.WriteString(indent + " * " + line + "\n")
		}
		out.WriteString(indent + " */\n")
	} else {
		for _, line := range fl.Comments {
			out.WriteString(indent + "// " + line + "\n")
		}
	}

	return out.String()
}

type AssignStatement struct {
	Token token.Token // '=' トークン
	Left  Expression
	Right Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String(depth int) string {
	return fmt.Sprintf("%s = %s;", as.Left.String(depth), as.Right.String(depth))
}

type DeriveStatement struct {
	Token token.Token
	Right Expression
}

func (es *DeriveStatement) statementNode()       {}
func (es *DeriveStatement) TokenLiteral() string { return es.Token.Literal }
func (es *DeriveStatement) String(depth int) string {
	return "..." + es.Right.String(depth)
}
