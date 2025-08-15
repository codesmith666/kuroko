package evaluator

import (
	"monkey/ast"
	"monkey/object"
	"strings"
)

var (
	NULL_OBJCT   = &object.Null{}
	TRUE_OBJECT  = &object.Boolean{Value: true}
	FALSE_OBJECT = &object.Boolean{Value: false}
)

func evalExpressions(
	exps []ast.Expression,
	env *object.Environment,
) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

/*
 * 単項演算子
 */
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE_OBJECT:
		return FALSE_OBJECT
	case FALSE_OBJECT:
		return TRUE_OBJECT
	case NULL_OBJCT:
		return TRUE_OBJECT
	default:
		return FALSE_OBJECT
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

/*
 * 二項演算子
 */
func evalInfixExpression(
	node *ast.InfixExpression,
	env *object.Environment,
) object.Object {
	operator := node.Operator

	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	right := Eval(node.Right, env)
	if isError(right) {
		return right
	}

	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return evalBoolLiteral(left == right)
	case operator == "!=":
		return evalBoolLiteral(left != right)
	case operator == "instanceof":
		return evalInstanceOfExpression(left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "<":
		return evalBoolLiteral(leftVal < rightVal)
	case ">":
		return evalBoolLiteral(leftVal > rightVal)
	case "==":
		return evalBoolLiteral(leftVal == rightVal)
	case "!=":
		return evalBoolLiteral(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalStringInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	if operator != "+" {
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}

	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	return &object.String{Value: leftVal + rightVal}
}

/*
 * if式
 */
func evalIfExpression(
	ie *ast.IfExpression,
	env *object.Environment,
) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	isTruthy := func(obj object.Object) bool {
		switch obj {
		case NULL_OBJCT:
			return false
		case TRUE_OBJECT:
			return true
		case FALSE_OBJECT:
			return false
		default:
			return true
		}
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL_OBJCT
	}
}

/*
 * インデックスアクセス
 */
func evalIndexExpression(node *ast.IndexExpression, env *object.Environment) object.Object {

	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}
	index := Eval(node.Index, env)
	if isError(index) {
		return index
	}

	switch {
	// 配列のインデックスアクセス
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		arrayObject := left.(*object.Array)
		idx := index.(*object.Integer).Value
		max := int64(len(arrayObject.Elements) - 1)

		if idx < 0 || idx > max {
			return NULL_OBJCT
		}

		return arrayObject.Elements[idx]
	// ハッシュのインデックスアクセス
	case left.Type() == object.HASH_OBJ:
		hashObject := left.(*object.Hash)

		key, ok := index.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", index.Type())
		}

		pair, ok := hashObject.Pairs[key.HashKey()]
		if !ok {
			return NULL_OBJCT
		}
		return pair.Value
	default:
		return newError("index operator not supported: %s", left.Type())
	}
}

/*
 * ハッシュアクセス
 */
func evalDotExpression(node *ast.DotExpression, env *object.Environment) object.Object {
	// 左辺はハッシュ
	// 右辺はenvからGetできない識別子なので名前を取得する。
	left := Eval(node.Left, env)
	right := &object.String{Value: node.Right.Name}

	// ハッシュかどうかチェック
	hashObj, ok := left.(*object.Hash)
	if !ok {
		return newError("not a hash: %s", left.Type())
	}

	// ハッシュオブジェクトから名前で値を取得
	key := &object.String{Value: right.Value}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL_OBJCT
	}
	return pair.Value
}

/*
 * 関数呼び出し
 */
func evalCallExpression(
	ce *ast.CallExpression,
	env *object.Environment,

) object.Object {

	function := Eval(ce.Function, env)
	if isError(function) {
		return function
	}

	args := evalExpressions(ce.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	switch fn := function.(type) {

	case *object.Function:
		// 関数の実行環境を拡張する
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)
		for paramIdx, param := range fn.Parameters {
			extendedEnv.Set(param.Name, args[paramIdx])
		}
		evaluated := Eval(fn.Body, extendedEnv)
		// 戻り値を取得する
		if returnValue, ok := evaluated.(*object.ReturnValue); ok {
			result := returnValue.Value
			if class := result.(*object.Class); class != nil {
				class.Name = fn.Name
			}
			return result
		}
		return evaluated

	case *object.Builtin:
		return fn.Fn(args...)

	default:
		return newError("not a function: %s", fn.Type())
	}
}

func evalInstanceOfExpression(left object.Object, right object.Object) object.Object {

	switch l := left.(type) {
	case *object.Class:
		// fmt.Printf("%T", right)
		switch r := right.(type) {
		case *object.Function:
			//	名前が一致したら
			if l.Name == r.Name {
				return evalBoolLiteral(true)
			}
			// 子クラスを検索
			if _, ok := l.Children[r.Name]; ok {
				return evalBoolLiteral(true)
			}
			return evalBoolLiteral(false)
		default:
			return newError("right operand of instanceof must be a primitive, got %s", right.Type())
		}
	default:
		switch r := right.(type) {
		case *object.Type:
			leftType := strings.ToLower(string(left.Type()))
			rightType := strings.ToLower(r.Name)
			return evalBoolLiteral(leftType == rightType)
		default:
			return newError("right operand of instanceof must be a primitive, got %s", right.Type())
		}
	}

}
