package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
	"strings"
)

// 前置演算子
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String(depth int) string {
	var out bytes.Buffer
	if pe.Token.Type == token.PARSE {
		out.WriteString(Indent(depth))
	}
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String(depth + 1))
	out.WriteString(")")
	if pe.Token.Type == token.PARSE {
		out.WriteString(";\n")
	}
	return out.String()
}

// 二項演算子
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String(depth int) string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String(depth + 1))
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String(depth + 1))
	out.WriteString(")")

	return out.String()
}

// if式
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String(depth int) string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(ie.Condition.String(depth + 1))
	out.WriteString("){")
	out.WriteString(ie.Consequence.String(depth + 1))
	out.WriteString("}")

	if ie.Alternative != nil {
		out.WriteString(" else{")
		out.WriteString(ie.Alternative.String(depth + 1))
		out.WriteString("}")
	}
	out.WriteString("\n")

	return out.String()
}

// 配列インデックス式
type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String(depth int) string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String(depth + 1))
	out.WriteString("[")
	out.WriteString(ie.Index.String(depth + 1))
	out.WriteString("])")

	return out.String()
}

// アクセス演算子
type DotExpression struct {
	Token token.Token // The '.' token
	Left  Expression  // person
	Right *Identifier // name
}

func (de *DotExpression) expressionNode()      {}
func (de *DotExpression) TokenLiteral() string { return de.Token.Literal }
func (de *DotExpression) String(depth int) string {
	return fmt.Sprintf("%s.%s", de.Left.String(depth), de.Right.String(depth))
}

// 関数呼び出し
type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String(depth int) string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String(depth+1))
	}

	out.WriteString(ce.Function.String(depth + 1))
	//out.WriteString(Type(ce))
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
