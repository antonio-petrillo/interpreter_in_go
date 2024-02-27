package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
)

var (
	TRUE  =  &object.Boolean{Value: true}
	FALSE =  &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.ReturnStatement:
		val :=  Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.Program:
		// return evalStatements(node.Statements) // to generic function call
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean: // This typo cost me a lot
		return nativeBoolBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		// return evalStatements(node.Statements) // to generic function call
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	}

	return nil
}

func evalIdentifier(identifier *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(identifier.Value)
	if !ok {
		return newError("identifier not found: " + identifier.Value)
	}

	return val
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment ) (result object.Object) {
	for _, statement := range block.Statements {
		result := Eval(statement, env)
		// if result != nil && result.Type() == object.RETURN_VALUE_OBJ {
		// 	return result
		// }
		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalProgram(program *ast.Program, env *object.Environment) (result object.Object) {
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		// if returnValue, ok := result.(*object.ReturnValue); ok {
		// 	return returnValue.Value
		// }

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		// return NULL
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())

	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	lValue := left.(*object.Integer).Value
	rValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: lValue + rValue}
	case "-":
		return &object.Integer{Value: lValue - rValue}
	case "*":
		return &object.Integer{Value: lValue * rValue}
	case "/":
		return &object.Integer{Value: lValue / rValue}
	case "<":
		return nativeBoolBooleanObject(lValue < rValue)
	case "<=":
		return nativeBoolBooleanObject(lValue <= rValue)
	case ">":
		return nativeBoolBooleanObject(lValue > rValue)
	case ">=":
		return nativeBoolBooleanObject(lValue >= rValue)
	case "==":
		return nativeBoolBooleanObject(lValue == rValue)
	case "!=":
		return nativeBoolBooleanObject(lValue != rValue)
	default:
		// return NULL
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		// return NULL
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalMinusPrefixOperatorExpression(obj object.Object) object.Object {
	if obj.Type() != object.INTEGER_OBJ {
		// return NULL
		return newError("unknown operator: -%s", obj.Type())
	}
	value := obj.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalBangOperatorExpression(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func nativeBoolBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

// func evalStatements(statements []ast.Statement) (result object.Object) {
// 	for _, statement := range statements {
// 		result = Eval(statement)

// 		if returnValue, ok := result.(*object.ReturnValue); ok {
// 			return returnValue.Value
// 		}
// 	}

// 	return result
// }

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.RETURN_VALUE_OBJ
	}
	return false
}
