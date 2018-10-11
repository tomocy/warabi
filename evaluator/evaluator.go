package evaluator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/tomocy/warabi/object"
)

const packageStatement = "package main\n"

var fileSet = token.NewFileSet()

func Evaluate(src string) []object.Object {
	file, err := parser.ParseFile(fileSet, "main.go", packageStatement+src, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	return evaluateDeclarations(file.Decls)
}

func evaluateDeclarations(decls []ast.Decl) []object.Object {
	var objs []object.Object
	for _, decl := range decls {
		objs = append(objs, evaluateDeclaration(decl)...)
	}

	return objs
}

func evaluateDeclaration(decl ast.Decl) []object.Object {
	switch decl := decl.(type) {
	case *ast.GenDecl:
		return evaluateGenericsDeclaration(decl)
	default:
		return nil
	}
}

func evaluateGenericsDeclaration(decl *ast.GenDecl) []object.Object {
	var objs []object.Object
	for _, spec := range decl.Specs {
		objs = append(objs, evaluateSpecification(spec)...)
	}

	return objs
}

func evaluateSpecification(spec ast.Spec) []object.Object {
	switch spec := spec.(type) {
	case *ast.ValueSpec:
		return evaluateValueSpecification(spec)
	default:
		return nil
	}
}

func evaluateValueSpecification(spec *ast.ValueSpec) []object.Object {
	var objs []object.Object
	for i := 0; i < len(spec.Names); i++ {
		objs = append(objs, evaluateExpression(spec.Values[i])...)
	}

	return objs
}

func evaluateExpression(expr ast.Expr) []object.Object {
	switch expr := expr.(type) {
	case *ast.ParenExpr:
		return evaluateParenExpression(expr)
	case *ast.BinaryExpr:
		return evaluateBinaryExpression(expr)
	case *ast.UnaryExpr:
		return evaluateUnaryExpression(expr)
	case *ast.BasicLit:
		return evaluateBasicLiteral(expr)
	default:
		return nil
	}
}

func evaluateParenExpression(expr *ast.ParenExpr) []object.Object {
	return evaluateExpression(expr.X)
}

func evaluateBinaryExpression(expr *ast.BinaryExpr) []object.Object {
	leftObjs := evaluateExpression(expr.X)
	if len(leftObjs) != 1 {
		return nil
	}
	leftObj := leftObjs[0]

	rightObjs := evaluateExpression(expr.Y)
	if len(rightObjs) != 1 {
		return nil
	}
	rightObj := rightObjs[0]

	switch {
	case leftObj.Kind() == object.Integer && rightObj.Kind() == object.Integer:
		return evaluateBinaryExpressionOfIntegerLiteral(leftObj.(*object.IntegerLiteral), expr.Op, rightObj.(*object.IntegerLiteral))
	case leftObj.Kind() == object.String && rightObj.Kind() == object.String:
		return evaluateBinaryExpressionOfStringLiteral(leftObj.(*object.StringLiteral), expr.Op, rightObj.(*object.StringLiteral))
	case leftObj.Kind() == object.Character && rightObj.Kind() == object.Character:
		return evaluateBinaryExpressionOfCharacterLiteral(leftObj.(*object.CharacterLiteral), expr.Op, rightObj.(*object.CharacterLiteral))
	case leftObj.Kind() == object.FloatingPoint && rightObj.Kind() == object.FloatingPoint:
		return evaluateBinaryExpressionOfFloatingPointLiteral(leftObj.(*object.FloatingPointLiteral), expr.Op, rightObj.(*object.FloatingPointLiteral))
	case leftObj.Kind() == object.FloatingPoint && rightObj.Kind() == object.Integer:
		intObj := rightObj.(*object.IntegerLiteral)
		floatObj := &object.FloatingPointLiteral{
			Value: float32(intObj.Value),
		}
		return evaluateBinaryExpressionOfFloatingPointLiteral(leftObj.(*object.FloatingPointLiteral), expr.Op, floatObj)
	case leftObj.Kind() == object.Integer && rightObj.Kind() == object.FloatingPoint:
		intObj := leftObj.(*object.IntegerLiteral)
		floatObj := &object.FloatingPointLiteral{
			Value: float32(intObj.Value),
		}
		return evaluateBinaryExpressionOfFloatingPointLiteral(floatObj, expr.Op, rightObj.(*object.FloatingPointLiteral))
	default:
		return nil
	}
}

func evaluateBinaryExpressionOfIntegerLiteral(leftObj *object.IntegerLiteral, operator token.Token, rightObj *object.IntegerLiteral) []object.Object {
	switch operator {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		return evaluateArithmeticOperation(leftObj, operator, rightObj)
	default:
		return nil
	}
}

func evaluateArithmeticOperation(leftObj *object.IntegerLiteral, operator token.Token, rightObj *object.IntegerLiteral) []object.Object {
	switch operator {
	case token.ADD:
		leftObj.Value += rightObj.Value
	case token.SUB:
		leftObj.Value -= rightObj.Value
	case token.MUL:
		leftObj.Value *= rightObj.Value
	case token.QUO:
		if rightObj.Value == 0 {
			return nil
		}
		leftObj.Value /= rightObj.Value
	case token.REM:
		leftObj.Value %= rightObj.Value
	default:
		return nil
	}

	return []object.Object{
		leftObj,
	}
}

func evaluateBinaryExpressionOfStringLiteral(leftObj *object.StringLiteral, operator token.Token, rightObj *object.StringLiteral) []object.Object {
	if operator != token.ADD {
		return nil
	}
	leftObj.Value += rightObj.Value
	return []object.Object{
		leftObj,
	}
}

func evaluateBinaryExpressionOfCharacterLiteral(leftObj *object.CharacterLiteral, operator token.Token, rightObj *object.CharacterLiteral) []object.Object {
	switch operator {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		return evaluateArithmeticOperationOfCharacterLiteral(leftObj, operator, rightObj)
	default:
		return nil
	}
}

func evaluateArithmeticOperationOfCharacterLiteral(leftObj *object.CharacterLiteral, operator token.Token, rightObj *object.CharacterLiteral) []object.Object {
	switch operator {
	case token.ADD:
		leftObj.Value += rightObj.Value
	case token.SUB:
		leftObj.Value -= rightObj.Value
	case token.MUL:
		leftObj.Value *= rightObj.Value
	case token.QUO:
		if rightObj.Value == 0 {
			return nil
		}
		leftObj.Value /= rightObj.Value
	case token.REM:
		leftObj.Value %= rightObj.Value
	default:
		return nil
	}

	return []object.Object{
		leftObj,
	}
}

func evaluateBinaryExpressionOfFloatingPointLiteral(leftObj *object.FloatingPointLiteral, operator token.Token, rightObj *object.FloatingPointLiteral) []object.Object {
	switch operator {
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		return evaluateArithmeticOperationOfFloatingPointLiteral(leftObj, operator, rightObj)
	default:
		return nil
	}
}

func evaluateArithmeticOperationOfFloatingPointLiteral(leftObj *object.FloatingPointLiteral, operator token.Token, rightObj *object.FloatingPointLiteral) []object.Object {
	switch operator {
	case token.ADD:
		leftObj.Value += rightObj.Value
	case token.SUB:
		leftObj.Value -= rightObj.Value
	case token.MUL:
		leftObj.Value *= rightObj.Value
	case token.QUO:
		if rightObj.Value == 0.0 {
			return nil
		}
		leftObj.Value /= rightObj.Value
	default:
		return nil
	}

	return []object.Object{
		leftObj,
	}
}

func evaluateUnaryExpression(expr *ast.UnaryExpr) []object.Object {
	switch expr.Op {
	case token.SUB:
		return evaluateMinusOperation(expr)
	default:
		return nil
	}
}

func evaluateMinusOperation(expr *ast.UnaryExpr) []object.Object {
	objs := evaluateExpression(expr.X)
	if len(objs) != 1 {
		return nil
	}

	obj, ok := objs[0].(*object.IntegerLiteral)
	if !ok {
		return nil
	}

	obj.Value *= -1
	return []object.Object{
		obj,
	}
}

func evaluateBasicLiteral(expr *ast.BasicLit) []object.Object {
	switch expr.Kind {
	case token.INT:
		return evaluateIntegerLiteral(expr)
	case token.STRING:
		return evaluateStringLiteral(expr)
	case token.CHAR:
		return evaluateCharacterLiteral(expr)
	case token.FLOAT:
		return evaluateFloatingPointLiteral(expr)
	default:
		return nil
	}
}

func evaluateIntegerLiteral(expr *ast.BasicLit) []object.Object {
	value, err := strconv.Atoi(expr.Value)
	if err != nil {
		return nil
	}
	return []object.Object{
		&object.IntegerLiteral{
			Value: value,
		},
	}
}

func evaluateStringLiteral(expr *ast.BasicLit) []object.Object {
	return []object.Object{
		&object.StringLiteral{
			Value: expr.Value[1 : len(expr.Value)-1],
		},
	}
}

func evaluateCharacterLiteral(expr *ast.BasicLit) []object.Object {
	return []object.Object{
		&object.CharacterLiteral{
			Value: []rune(expr.Value[1 : len(expr.Value)-1])[0],
		},
	}
}

func evaluateFloatingPointLiteral(expr *ast.BasicLit) []object.Object {
	value, err := strconv.ParseFloat(expr.Value, 32)
	if err != nil {
		return nil
	}
	return []object.Object{
		&object.FloatingPointLiteral{
			Value: float32(value),
		},
	}
}
