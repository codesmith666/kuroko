package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"regexp"
	"strings"
)

/*
 * 変数の宣言ステートメント
 * イミュータブルとミュータブルがある
 */
//
type LetStatement struct {
	Token token.Token // the token.LET token
	Ident *Identifier // 識別子の名前
	Value Expression  // 式
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Ident.String())
	out.WriteString(" = ")
	// 値
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	} else {
		out.WriteString("<?>")
	}
	out.WriteString(";\n")
	return out.String()
}

/*
 * 式ステートメント
 * 値を返すけど結果を捨ててよい式（関数呼び出し）など
 */
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

/*
 * リターンステートメント
 */
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	if rs.ReturnValue != nil {
		out.WriteString("return ")
		out.WriteString(rs.ReturnValue.String())
		out.WriteString(";\n")
	}
	return out.String()
}

/*
 *	ブロックステートメント
 */
type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String() + ";") // +1しない
	}
	out.WriteString("}")
	return out.String()
}

/*
 * コメントステートメント
 * リテラルに生コメントが入っている
 */
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

func (fl *CommentStatement) String() string {

	var out bytes.Buffer
	out.WriteString("\n")
	if fl.Token.Type == token.BLOCK_COMMENT {
		out.WriteString("/*\n")
		for _, line := range fl.Comments {
			out.WriteString(" * " + line + "\n")
		}
		out.WriteString(" */\n")
	} else {
		for _, line := range fl.Comments {
			out.WriteString("// " + line + "\n")
		}
	}

	return out.String()
}

/*
 * 代入
 */
type AssignStatement struct {
	Token token.Token // '=' トークン
	Left  Expression
	Right Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	return fmt.Sprintf("%s = %s;", as.Left.String(), as.Right.String())
}

type DeriveStatement struct {
	Token token.Token
	Right Expression
}

func (es *DeriveStatement) statementNode()       {}
func (es *DeriveStatement) TokenLiteral() string { return es.Token.Literal }
func (es *DeriveStatement) String() string {
	return "..." + es.Right.String()
}

/*
 * 繰り返し構文
 */
type LoopStatement struct {
	Token token.Token     // 'for' トークン
	Bind  *LetStatement   // forの括弧の中
	Block *BlockStatement // ループする処理
}

func (ls *LoopStatement) statementNode()       {}
func (ls *LoopStatement) TokenLiteral() string { return ls.Token.Literal }
func (fs *LoopStatement) String() string {
	var out bytes.Buffer

	out.WriteString("loop(")
	out.WriteString(fs.Bind.String())
	out.WriteString(") ")
	out.WriteString(fs.Block.String())
	return out.String()
}

/*
 * Break
 */
type BreakStatement struct {
	Token token.Token
}

func (bs *BreakStatement) statementNode()       {}
func (bs *BreakStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BreakStatement) String() string       { return "break" }

/*
 * Continue
 */
type ContinueStatement struct {
	Token token.Token
}

func (cs *ContinueStatement) statementNode()       {}
func (cs *ContinueStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *ContinueStatement) String() string       { return "continue" }
