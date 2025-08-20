package object

import (
	"bytes"
	"strings"
)

/*
 * クラス
 */
type Class struct {
	Hash
	name     string
	children map[string]struct{}
}

func NewClass() *Class {
	return &Class{
		Hash:     *NewHash(),
		name:     "$unnamed",
		children: make(map[string]struct{}),
	}

}

func (c *Class) InstanceOf(name string) bool {
	if c.name == name {
		return true
	}
	// 子クラスを検索
	if _, ok := c.children[name]; ok {
		return true
	}
	return false
}

func (c *Class) Derive(from *Class) {
	// パラメータを継承する
	c.Hash.Merge(&from.Hash)

	// 子クラス名をすべて引き継ぐ
	for childName, _ := range from.children {
		c.children[childName] = struct{}{}
	}
	c.children[from.name] = struct{}{}
}

func (c *Class) ClassName() string {
	return c.name
}

func (c *Class) SetClassName(className string) bool {
	if c.name == "$unnamed" {
		c.name = className
		return true
	}
	return false
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string {
	var out bytes.Buffer

	out.WriteString(c.name)
	out.WriteString(c.Hash.Inspect())
	if c.children == nil {
		out.WriteString("from ()")
	} else {
		out.WriteString("from (")
		from := []string{}
		for k, _ := range c.children {
			from = append(from, k)
		}
		out.WriteString(strings.Join(from, ","))
		out.WriteString(")")
	}
	return out.String()
}
