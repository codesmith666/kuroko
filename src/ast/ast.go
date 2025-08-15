package ast

import (
	"bytes"
	"fmt"
	"reflect"
)

// The base Node interface
type Node interface {
	TokenLiteral() string
	String(depth int) string
}

// All statement nodes implement this
type Statement interface {
	Node
	statementNode()
}

// All expression nodes implement this
type Expression interface {
	Node
	expressionNode()
}

func Indent(depth int) string {
	var out bytes.Buffer
	for i := 0; i < depth; i++ {
		out.WriteString("  ")
	}
	return out.String()
}

func Type(node Node) string {
	return reflect.TypeOf(node).String()
}

/*
 * astツリーを表示する
 */
func PrintAST(node Node, indent string) {

	if node == nil {
		return
	}

	t := reflect.TypeOf(node)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	fmt.Println(indent + "* " + t.Name())

	switch n := node.(type) {
	case *Program:
		for _, stmt := range n.Statements {
			PrintAST(stmt, indent+"  ")
		}

	case *LetStatement:
		PrintAST(n.Ident, indent+"  ")
		PrintAST(n.Value, indent+"  ")

	case *ExpressionStatement:
		fmt.Printf("%s %s\n", indent+"  ", n.Token.Type)
		PrintAST(n.Expression, indent+"  ")

	case *PrefixExpression:
		fmt.Printf("%sOperator: %s\n", indent+"  ", n.Operator)
		PrintAST(n.Right, indent+"  ")

	case *InfixExpression:
		fmt.Printf("%s  [left]\n", indent)
		PrintAST(n.Left, indent+"  ")
		fmt.Printf("%s  [Operator]\n", indent)
		fmt.Printf("%s  %s\n", indent, n.Operator)
		fmt.Printf("%s  [right]\n", indent)
		PrintAST(n.Right, indent+"  ")

	case *IntegerLiteral:
		fmt.Printf("%s  %d\n", indent, n.Value)

	case *StringLiteral:
		fmt.Printf("%s  %s\n", indent, n.Value)

	case *CallExpression:
		fmt.Printf("%s  [function]\n", indent)
		PrintAST(n.Function, indent+"  ")
		fmt.Printf("%s  [arguments]\n", indent)
		for _, a := range n.Arguments {
			PrintAST(a, indent+"  ")
		}

	case *FunctionLiteral:
		fmt.Printf("%s  [parameters]\n", indent)
		if len(n.Parameters) == 0 {
			fmt.Printf("%s  None\n", indent)
		} else {
			for _, p := range n.Parameters {
				PrintAST(p, indent+"  ")
			}
		}
		fmt.Printf("%s  [return type]\n", indent)
		if n.ReturnType != nil {
			fmt.Printf("%s  %s\n", indent, n.ReturnType.String())
		} else {
			fmt.Printf("%s  ?\n", indent)
		}

		fmt.Printf("%s  [block]\n", indent)
		for _, s := range n.Body.Statements {
			PrintAST(s, indent+"  ")
		}

	case *Identifier:
		var t string
		if n.Type != nil {
			t = n.Type.String()
		} else {
			t = "?"
		}
		fmt.Printf("%s  %s: %s\n", indent, n.Name, t)

	case *TypeLiteral:
		fmt.Printf("%s  %s\n", indent, n.Value)

	case *HashLiteral:
		for k, v := range n.Pairs {
			fmt.Printf("%s  [key]\n", indent)
			PrintAST(k, indent+"  ")
			fmt.Printf("%s  [value]\n", indent)
			PrintAST(v, indent+"  ")
		}

	case *ReturnStatement:
		if n.ReturnValue != nil {
			PrintAST(n.ReturnValue, indent+"  ")
		} else {
			fmt.Printf("%s  None\n", indent)
		}
	case *CommentStatement:
		for _, c := range n.Comments {
			fmt.Printf("%s  %s\n", indent, c)
		}
	case *IndexExpression:
		fmt.Printf("%s  [left]\n", indent)
		PrintAST(n.Left, indent+"  ")
		fmt.Printf("%s  [index]\n", indent)
		PrintAST(n.Index, indent+"  ")
	case *DotExpression:
		fmt.Printf("%s  [left]\n", indent)
		PrintAST(n.Left, indent+"  ")
		fmt.Printf("%s  [right]\n", indent)
		PrintAST(n.Right, indent+"  ")
	case *AssignStatement:
		fmt.Printf("%s  [left]\n", indent)
		PrintAST(n.Left, indent+"  ")
		fmt.Printf("%s  [right]\n", indent)
		PrintAST(n.Right, indent+"  ")

	default:
		fmt.Printf("%sUnknown node: %T\n", indent, n)
	}
}
