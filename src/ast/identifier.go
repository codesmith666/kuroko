package ast

import (
	"bytes"
	"monkey/token"
)

// Expressions
type Identifier struct {
	Token token.Token // the token.IDENT token
	Name  string      // 識別子の名前が入っている
	Type  *TypeNode   // 識別子の型が入っている（決まっていないときはnil）
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String(depth int) string {
	var out bytes.Buffer
	out.WriteString(i.Name)
	out.WriteString(":")
	if i.Type != nil {
		out.WriteString("<" + i.Type.String() + ">")
	} else {
		out.WriteString("<?>")
	}
	return out.String()
}
