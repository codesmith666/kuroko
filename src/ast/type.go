package ast

import (
	"monkey/token"
	"strings"
)

type TypeKind string

const (
	TypeSimple   TypeKind = "Simple"
	TypeArray    TypeKind = "Array"
	TypeMap      TypeKind = "Map" // インデックスシグネチャ
	TypeObject   TypeKind = "Object"
	TypeFunction TypeKind = "Function"
)

type ObjectProperty struct {
	Name string
	Type *TypeNode
}

type TypeNode struct {
	Token       token.Token // ':' トークン
	Kind        TypeKind
	Name        string    // Kind == Simple → 型名（例: number, string）
	ElementType *TypeNode // Kind == Array → 要素の型
	KeyType     *TypeNode // Kind == Map → keyの型
	ValueType   *TypeNode // Kind == Map → valueの型
	Properties  []*ObjectProperty
	Parameters  []*Identifier
	ReturnType  *TypeNode
}

func (tn *TypeNode) TokenLiteral() string {
	return tn.Token.Literal
}

// 型は開業しない
func (tn *TypeNode) String() string {
	debug := func() string {
		switch tn.Kind {
		case TypeSimple:
			return tn.Name
		case TypeArray:
			return tn.ElementType.String() + "[]"
		case TypeMap:
			return "{ [" + tn.KeyType.String() + "]: " + tn.ValueType.String() + " }"
		case TypeObject:
			props := []string{}
			for _, p := range tn.Properties {
				props = append(props, p.Name+": "+p.Type.String())
			}
			return "{ " + strings.Join(props, ", ") + " }"
		case TypeFunction:
			params := []string{}
			for _, p := range tn.Parameters {
				params = append(params, p.Token.Literal+":"+p.Type.String())
			}
			return "(" + strings.Join(params, ", ") + ") => " + tn.ReturnType.String()
		default:
			return "unknown"
		}
	}
	return debug()
}
