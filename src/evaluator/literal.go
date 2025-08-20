package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

// nativeBoolToBooleanObject
func evalBoolLiteral(input bool) *object.Boolean {
	if input {
		return object.TRUE
	}
	return object.FALSE
}

func evalHashLiteral(
	node *ast.HashLiteral,
	env *object.Environment,
) object.Object {
	hash := object.NewHash()

	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}

		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		err := hash.Set(key, value)
		if err != nil {
			return newError("%s", err.Error())
		}
	}

	return hash
}
