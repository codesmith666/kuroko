package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

// 二値
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String(depth int) string {
	return bl.Token.Literal
}

// 整数
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String(depth int) string {
	var out bytes.Buffer
	out.WriteString(" ")
	out.WriteString(il.Token.Literal)
	out.WriteString(" ")
	return out.String()
}

// 文字列
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String(depth int) string {
	return "\"" + sl.Token.Literal + "\""
}

// 配列
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String(depth int) string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String(depth+1))
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// マップというかハッシュというかのリテラル
type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String(depth int) string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String(depth+1)+":"+value.String(depth+1))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

// 関数リテラル（呼び出しではない）
type FunctionLiteral struct {
	Token      token.Token // The 'fn' token
	Parameters []*Identifier
	Body       *BlockStatement
	ReturnType *TypeNode
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String(depth int) string {
	var out bytes.Buffer
	var params = []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.Name+":"+p.Type.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if fl.ReturnType != nil {
		out.WriteString(":")
		out.WriteString(fl.ReturnType.String())
	}
	out.WriteString(" => ")
	out.WriteString(fl.Body.String(depth)) // +1しない
	return out.String()
}

// 型リテラル
type TypeLiteral struct {
	Token token.Token // 'string', 'integer', etc.
	Value string
}

func (tl *TypeLiteral) expressionNode()         {}
func (tl *TypeLiteral) TokenLiteral() string    { return tl.Token.Literal }
func (tl *TypeLiteral) String(depth int) string { return tl.Value }
