package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strings"
)

type BuiltinFunction func(args ...Object) Object

type ObjectType string

/*
 * 実行時の型
 */
const (
	NULL_OBJ  ObjectType = "NULL"
	ERROR_OBJ ObjectType = "ERROR"

	INTEGER_OBJ ObjectType = "INTEGER"
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
	STRING_OBJ  ObjectType = "STRING"

	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"

	FUNCTION_OBJ ObjectType = "FUNCTION"
	BUILTIN_OBJ  ObjectType = "BUILTIN"

	ARRAY_OBJ ObjectType = "ARRAY"
	HASH_OBJ  ObjectType = "HASH"

	TYPE_OBJ  ObjectType = "TYPE"
	CLASS_OBJ ObjectType = "CLASS"
)

/*
 *	ハッシュキーの型
 */
type HashKey struct {
	Type  ObjectType
	Value uint64
}

/*
 * インタフェイス
 */
type Hashable interface {
	HashKey() HashKey
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

/*
 * 整数
 */
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

/*
 * 二値
 */
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

/*
 *　null
 */
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

/*
 * 戻り値
 */
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

/*
 * エラー
 */
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

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
		params = append(params, p.String(0))
	}

	if f.Name == "" {
		out.WriteString("$unnamed")
	} else {
		out.WriteString(f.Name)
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") => \n")
	out.WriteString(f.Body.String(0))

	return out.String()
}

/*
 * 文字列
 */
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

/*
 * 組み込み
 */
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

/*
 * 配列
 */
type Array struct {
	Elements []Object
}

func (ao *Array) Type() ObjectType { return ARRAY_OBJ }
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

/*
 * 連想配列
 */
type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

/*
 * クラス
 */
type Class struct {
	Hash
	Name     string
	Children map[string]struct{}
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string {
	var out bytes.Buffer

	out.WriteString(c.Name)
	out.WriteString(c.Hash.Inspect())
	if c.Children == nil {
		out.WriteString("from ()")
	} else {
		out.WriteString("from (")
		from := []string{}
		for k, _ := range c.Children {
			from = append(from, k)
		}
		out.WriteString(strings.Join(from, ","))
		out.WriteString(")")
	}
	return out.String()
}

/*
 * 型
 */
type Type struct {
	Name string
}

func (t *Type) Type() ObjectType { return TYPE_OBJ }
func (t *Type) Inspect() string  { return t.Name }
