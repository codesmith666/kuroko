package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

func evalIdentifier(
	node *ast.Identifier,
	env *object.Environment,
) object.Object {
	if val, ok := env.Get(node.Name); ok {
		return val
	}

	if builtin, ok := builtins[node.Name]; ok {
		return builtin
	}

	return newError("identifier not found: %s", node.Name)
}
