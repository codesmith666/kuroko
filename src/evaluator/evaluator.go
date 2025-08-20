package evaluator

import (
	"math/rand"
	"monkey/ast"
	"monkey/object"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	//
	// Statements
	//
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.DeriveStatement:
		return evalDeriveStatment(node, env)

	case *ast.AssignStatement:
		return evalAssignStatement(node, env)

	//
	// Literal
	//
	case *ast.TypeLiteral:
		return &object.Type{Name: node.Value}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}

	case *ast.ComplexLiteral:
		return &object.Complex{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.BooleanLiteral:
		return evalBoolLiteral(node.Value)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		name := "function-" + RandString(16)
		return &object.Function{
			Parameters: params,
			Env:        env,
			Body:       body,
			Name:       name,
		}
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	//
	// Expression
	//
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		return evalInfixExpression(node, env)

	case *ast.IndexExpression:
		return evalIndexExpression(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.CallExpression:
		return evalCallExpression(node, env)

	case *ast.DotExpression:
		return evalDotExpression(node, env)

	//
	// identifier
	//
	case *ast.Identifier:
		return evalIdentifier(node, env)
	}

	return nil
}
