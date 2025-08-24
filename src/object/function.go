package object

import (
	"bytes"
	"monkey/ast"
	"strings"
)

/*
 * 関数
 */
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
	Name       string
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("------------\n")

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	if f.Name == "" {
		out.WriteString("$unnamed")
	} else {
		out.WriteString(f.Name)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") => \n")
	out.WriteString(f.Body.String())

	return out.String()
}

/*
 * 組み込み
 */
type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
