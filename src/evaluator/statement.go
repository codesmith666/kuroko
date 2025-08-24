package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

/*
 * ブロックステートメントを実行
 */
func evalBlockStatement(
	block *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		// エラーじゃなかったら
		if result != nil {
			// ループを抜けて終了するパターン
			switch result.Type() {
			case object.ERROR_OBJ:
				return result
			case object.RETURN_VALUE_OBJ:
				return result
			case object.BREAK_OBJ:
				return result
			case object.CONTINUE_OBJ:
				return result
			}
		}
	}
	return result
}

/*
 * 派生
 */
func evalDeriveStatment(
	node *ast.DeriveStatement,
	env *object.Environment,
) object.Object {

	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	if right == nil {
		return newError("parse operator(...) requires a hash,got=nil on line %d col %d.", node.Token.Row, node.Token.Col)
	}

	switch rightValue := right.(type) {
	case *object.Hash:
		env.DeriveFromHash(rightValue)
		return right
	case *object.Class:
		env.DeriveFromClass(rightValue)
		return right
	default:
		return newError("parse operator(...) requires a hash,got= %s", right.Type())
	}
}

/*
 * 変数束縛
 */
func evalLetStatement(node *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}
	env.Set(node.Ident.Name, val)
	return nil
}

/*
 * 変数への代入
 */
func evalAssignStatement(
	stmt *ast.AssignStatement,
	env *object.Environment,
) object.Object {
	right := Eval(stmt.Right, env)
	if isError(right) {
		return right
	}

	switch nameExpr := stmt.Left.(type) {

	// 変数再代入
	case *ast.Identifier:
		result := env.Get(nameExpr.Name)
		if result == nil {
			return newError("identifier not found: %s", nameExpr.Name)
		}
		env.Set(nameExpr.Name, right)
		return right

	// アクセス演算子
	case *ast.DotExpression:
		left := Eval(nameExpr.Left, env)
		if isError(left) {
			return left
		}
		index := nameExpr.Right

		switch leftObj := left.(type) {
		case *object.Hash:
			key := &object.String{Value: index.Name}
			leftObj.Set(key, right)
			return right
		default:
			return newError("assignment target not assignable: %s", left.Type())
		}
	// 添字代入
	case *ast.IndexExpression:
		left := Eval(nameExpr.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(nameExpr.Index, env)
		if isError(index) {
			return index
		}

		switch leftObj := left.(type) {
		case *object.Array:
			idx, ok := index.(*object.Integer)
			if !ok {
				return newError("array index is not integer: %s", index.Type())
			}
			if idx.Value < 0 || idx.Value >= int64(len(leftObj.Elements)) {
				return newError("index out of range")
			}
			leftObj.Elements[idx.Value] = right
			return right

		case *object.Hash:
			err := leftObj.Set(index, right)
			if err != nil {
				return newError("%s", err.Error())
			}
			return right

		default:
			return newError("assignment target not assignable: %s", left.Type())
		}
	default:
		return newError("invalid assignment target")
	}
}

func evalLoopStatement(
	node *ast.LoopStatement,
	env *object.Environment,
) object.Object {
	val := Eval(node.Bind.Value, env)
	if isError(val) {
		return val
	}
	key := node.Bind.Ident.Name
	exEnv := object.NewEnclosedEnvironment(env)
	iter := object.NewHash()
	index := int64(0)
	kk := &object.String{Value: "k"}
	kv := &object.String{Value: "v"}
	ki := &object.String{Value: "i"}

	switch hash := val.(type) {
	case *object.Hash:
		var ret object.Object = nil
		hash.Range(func(k *object.Object, v *object.Object) bool {
			iter.Set(kk, *k)
			iter.Set(kv, *v)
			iter.Set(ki, &object.Integer{Value: index})
			exEnv.Set(key, iter)
			evaluated := Eval(node.Block, exEnv)
			switch evaluated.(type) {
			case *object.Break:
				return false
			case *object.Continue:
				return true
			case *object.ReturnValue:
				ret = evaluated
				return false
			}
			index++
			return true
		})
		return ret
	}
	return nil
}
