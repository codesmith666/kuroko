package evaluator

import (
	"monkey/ast"
	"monkey/object"
	"strings"
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
	case object.TRUE:
		return object.FALSE
	case object.FALSE:
		return object.TRUE
	case object.NULL:
		return object.TRUE
	default:
		return object.FALSE
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
	case left.Type() == object.FLOAT_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator, left, right)
	case left.Type() == object.COMPLEX_OBJ && right.Type() == object.COMPLEX_OBJ:
		return evalComplexInfixExpression(operator, left, right)

	// 型昇格
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalFloatInfixExpression(operator,
			&object.Float{Value: float64(left.(*object.Integer).Value)}, right)

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalFloatInfixExpression(operator,
			left, &object.Float{Value: float64(right.(*object.Integer).Value)})

	case left.Type() == object.FLOAT_OBJ && right.Type() == object.COMPLEX_OBJ:
		return evalComplexInfixExpression(operator,
			&object.Complex{Value: complex(left.(*object.Float).Value, 0)}, right)

	case left.Type() == object.COMPLEX_OBJ && right.Type() == object.FLOAT_OBJ:
		return evalComplexInfixExpression(operator,
			left, &object.Complex{Value: complex(right.(*object.Float).Value, 0)})

	case left.Type() == object.INTEGER_OBJ && right.Type() == object.COMPLEX_OBJ:
		return evalComplexInfixExpression(operator,
			&object.Complex{Value: complex(float64(left.(*object.Integer).Value), 0)}, right)

	case left.Type() == object.COMPLEX_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalComplexInfixExpression(operator,
			left, &object.Complex{Value: complex(float64(right.(*object.Integer).Value), 0)})
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

func evalFloatInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Float).Value
	rightVal := right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: leftVal + rightVal}
	case "-":
		return &object.Float{Value: leftVal - rightVal}
	case "*":
		return &object.Float{Value: leftVal * rightVal}
	case "/":
		return &object.Float{Value: leftVal / rightVal}
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

func evalComplexInfixExpression(
	operator string,
	left, right object.Object,
) object.Object {
	leftVal := left.(*object.Complex).Value
	rightVal := right.(*object.Complex).Value

	switch operator {
	case "+":
		return &object.Complex{Value: leftVal + rightVal}
	case "-":
		return &object.Complex{Value: leftVal - rightVal}
	case "*":
		return &object.Complex{Value: leftVal * rightVal}
	case "/":
		return &object.Complex{Value: leftVal / rightVal}
	// case "<":
	// 	return evalBoolLiteral(leftVal < rightVal)
	// case ">":
	// 	return evalBoolLiteral(leftVal > rightVal)
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
		case object.NULL:
			return false
		case object.TRUE:
			return true
		case object.FALSE:
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
		return object.NULL
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
			return object.NULL
		}

		return arrayObject.Elements[idx]
	// ハッシュのインデックスアクセス
	case left.Type() == object.HASH_OBJ:
		hashObject := left.(*object.Hash)

		val, err := hashObject.Get(index)
		switch true {
		case err == nil:
			return val
		case err.Is(object.InvalidKey):
			return newError("%s", err.Error())
		case err.Is(object.NotFound):
			return object.UNDEFINED
		default:
			return newError("evalIndexExpression:Unreachable")
		}
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
	val, err := hashObj.Get(&object.String{Value: right.Value})
	switch true {
	case err == nil:
		return val
	case err.Is(object.InvalidKey):
		return newError("%s", err.Error())
	case err.Is(object.NotFound):
		return object.UNDEFINED
	default:
		return newError("evalDotExpression:Unreachable")
	}
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
				if !class.SetClassName(fn.Name) {
					return newError("class name already initialized.")
				}
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
			//	名前が一致したら真なのでTRUEオブジェクトを返す
			if l.InstanceOf(r.Name) {
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
