package ast

import (
	"bytes"
	"monkey/lib"
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
func (bl *BooleanLiteral) String() string {
	return bl.Token.Literal
}

// 整数
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(" ")
	out.WriteString(il.Token.Literal)
	out.WriteString(" ")
	return out.String()
}

// 浮動小数点
type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// 複素数
type ComplexLiteral struct {
	Token token.Token
	Value complex128
}

func (cl *ComplexLiteral) expressionNode()      {}
func (cl *ComplexLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *ComplexLiteral) String() string       { return cl.Token.Literal }

// 文字列
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string {
	return "\"" + sl.Token.Literal + "\""
}

// 配列
type ArrayLiteral struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// マップというかハッシュというかのリテラル
type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs *lib.OrderedMap[Expression, Expression]
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	hl.Pairs.Range(func(e1, e2 Expression) bool {
		pairs = append(pairs, e1.String()+":"+e2.String())
		return true
	})

	// for key, value := range hl.Pairs {
	// 	pairs = append(pairs, key.String()+":"+value.String())
	// }
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
func (fl *FunctionLiteral) String() string {
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
	out.WriteString(fl.Body.String())
	return out.String()
}

// 型リテラル
type TypeLiteral struct {
	Token token.Token // 'string', 'integer', etc.
	Value string
}

func (tl *TypeLiteral) expressionNode()      {}
func (tl *TypeLiteral) TokenLiteral() string { return tl.Token.Literal }
func (tl *TypeLiteral) String() string       { return tl.Value }
