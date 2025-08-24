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

	var err object.Object = nil
	node.Pairs.Range(func(k ast.Expression, v ast.Expression) bool {
		// キー側の式を評価
		key := Eval(k, env)
		if isError(key) {
			err = key
			return false
		}
		// 値側の式を評価
		value := Eval(v, env)
		if isError(value) {
			err = value
			return false
		}
		// ハッシュを保存
		e := hash.Set(key, value)
		if err != nil {
			err = newError("%s", e.Error())
			return false
		}
		return true
	})

	if err != nil {
		return err
	}

	// for keyNode, valueNode := range node.Pairs {
	// 	key := Eval(keyNode, env)
	// 	if isError(key) {
	// 		return key
	// 	}

	// 	value := Eval(valueNode, env)
	// 	if isError(value) {
	// 		return value
	// 	}

	// 	err := hash.Set(key, value)
	// 	if err != nil {
	// 		return newError("%s", err.Error())
	// 	}
	// }

	return hash
}
