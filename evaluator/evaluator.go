package evaluator

import (
	"fmt"
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
	if len(spec.Values) == 0 {
		spec = restoreZeroValues(*spec)
	}
	for i := 0; i < len(spec.Names); i++ {
		obj := evaluateExpression(spec.Values[i])
		object.Env.Set(spec.Names[i].Name, obj)
		objs = append(objs, obj)
	}

	return objs
}

var zeroValues = map[string]ast.Expr{
	"int": &ast.BasicLit{
		Kind:  token.INT,
		Value: "0",
	},
	"string": &ast.BasicLit{
		Kind:  token.STRING,
		Value: `""`,
	},
	"byte": &ast.BasicLit{
		Kind:  token.CHAR,
		Value: "'0'",
	},
	"rune": &ast.BasicLit{
		Kind:  token.CHAR,
		Value: "'0'",
	},
	"float32": &ast.BasicLit{
		Kind:  token.FLOAT,
		Value: "0.0",
	},
}

func restoreZeroValues(spec ast.ValueSpec) *ast.ValueSpec {
	spec.Values = make([]ast.Expr, len(spec.Names))
	ident, ok := spec.Type.(*ast.Ident)
	if !ok {
		return &spec
	}

	zeroValue, ok := zeroValues[ident.Name]
	if !ok {
		return &spec
	}

	for i := 0; i < len(spec.Names); i++ {
		spec.Values[i] = zeroValue
	}

	return &spec
}

func evaluateExpression(expr ast.Expr) object.Object {
	switch expr := expr.(type) {
	case *ast.ParenExpr:
		return evaluateParenOperation(expr)
	case *ast.BinaryExpr:
		return evaluateBinaryOperation(expr)
	case *ast.UnaryExpr:
		return evaluateUnaryOperation(expr)
	case *ast.Ident:
		return evaluateIdentifier(expr)
	case *ast.BasicLit:
		return evaluateBasicLiteral(expr)
	default:
		return nil
	}
}

func evaluateParenOperation(expr *ast.ParenExpr) object.Object {
	return evaluateExpression(expr.X)
}

func evaluateBinaryOperation(expr *ast.BinaryExpr) object.Object {
	leftObj := evaluateExpression(expr.X)
	rightObj := evaluateExpression(expr.Y)
	switch {
	case leftObj.Kind() == object.Integer && rightObj.Kind() == object.Integer:
		return evaluateBinaryOperationOfIntegerLiteral(
			leftObj.(*object.IntegerLiteral),
			expr.Op,
			rightObj.(*object.IntegerLiteral),
		)
	case leftObj.Kind() == object.String && rightObj.Kind() == object.String:
		return evaluateBinaryOperationOfStringLiteral(
			leftObj.(*object.StringLiteral),
			expr.Op,
			rightObj.(*object.StringLiteral),
		)
	case leftObj.Kind() == object.Character && rightObj.Kind() == object.Character:
		return evaluateBinaryOperationOfCharacterLiteral(
			leftObj.(*object.CharacterLiteral),
			expr.Op,
			rightObj.(*object.CharacterLiteral),
		)
	case leftObj.Kind() == object.FloatingPoint && rightObj.Kind() == object.FloatingPoint:
		return evaluateBinaryOperationOfFloatingPointLiteral(
			leftObj.(*object.FloatingPointLiteral),
			expr.Op,
			rightObj.(*object.FloatingPointLiteral),
		)
	case leftObj.Kind() == object.FloatingPoint && rightObj.Kind() == object.Integer:
		intObj := rightObj.(*object.IntegerLiteral)
		floatObj := &object.FloatingPointLiteral{
			Value: float32(intObj.Value),
		}
		return evaluateBinaryOperationOfFloatingPointLiteral(
			leftObj.(*object.FloatingPointLiteral),
			expr.Op,
			floatObj,
		)
	case leftObj.Kind() == object.Integer && rightObj.Kind() == object.FloatingPoint:
		intObj := leftObj.(*object.IntegerLiteral)
		floatObj := &object.FloatingPointLiteral{
			Value: float32(intObj.Value),
		}
		return evaluateBinaryOperationOfFloatingPointLiteral(
			floatObj,
			expr.Op,
			rightObj.(*object.FloatingPointLiteral),
		)
	default:
		return nil
	}
}

func evaluateBinaryOperationOfIntegerLiteral(
	leftObj *object.IntegerLiteral,
	operator token.Token,
	rightObj *object.IntegerLiteral,
) object.Object {
	switch operator {
	case token.ADD:
		return &object.IntegerLiteral{Value: leftObj.Value + rightObj.Value}
	case token.SUB:
		return &object.IntegerLiteral{Value: leftObj.Value - rightObj.Value}
	case token.MUL:
		return &object.IntegerLiteral{Value: leftObj.Value * rightObj.Value}
	case token.QUO:
		if rightObj.Value == 0 {
			return nil
		}
		return &object.IntegerLiteral{Value: leftObj.Value / rightObj.Value}
	case token.REM:
		return &object.IntegerLiteral{Value: leftObj.Value % rightObj.Value}
	case token.LSS:
		return convertToBooleanLiteral(leftObj.Value < rightObj.Value)
	case token.GTR:
		return convertToBooleanLiteral(leftObj.Value > rightObj.Value)
	case token.LEQ:
		return convertToBooleanLiteral(leftObj.Value <= rightObj.Value)
	case token.GEQ:
		return convertToBooleanLiteral(leftObj.Value >= rightObj.Value)
	default:
		return nil
	}
}

func evaluateBinaryOperationOfStringLiteral(
	leftObj *object.StringLiteral,
	operator token.Token,
	rightObj *object.StringLiteral,
) object.Object {
	switch operator {
	case token.ADD:
		return &object.StringLiteral{Value: leftObj.Value + rightObj.Value}
	case token.LSS:
		return convertToBooleanLiteral(leftObj.Value < rightObj.Value)
	case token.GTR:
		return convertToBooleanLiteral(leftObj.Value > rightObj.Value)
	case token.LEQ:
		return convertToBooleanLiteral(leftObj.Value <= rightObj.Value)
	case token.GEQ:
		return convertToBooleanLiteral(leftObj.Value >= rightObj.Value)
	default:
		return nil
	}
}

func evaluateBinaryOperationOfCharacterLiteral(
	leftObj *object.CharacterLiteral,
	operator token.Token,
	rightObj *object.CharacterLiteral,
) object.Object {
	switch operator {
	case token.ADD:
		return &object.CharacterLiteral{Value: leftObj.Value + rightObj.Value}
	case token.SUB:
		return &object.CharacterLiteral{Value: leftObj.Value - rightObj.Value}
	case token.MUL:
		return &object.CharacterLiteral{Value: leftObj.Value * rightObj.Value}
	case token.QUO:
		if rightObj.Value == 0 {
			return nil
		}
		return &object.CharacterLiteral{Value: leftObj.Value / rightObj.Value}
	case token.REM:
		return &object.CharacterLiteral{Value: leftObj.Value % rightObj.Value}
	case token.LSS:
		return convertToBooleanLiteral(leftObj.Value < rightObj.Value)
	case token.GTR:
		return convertToBooleanLiteral(leftObj.Value > rightObj.Value)
	case token.LEQ:
		return convertToBooleanLiteral(leftObj.Value <= rightObj.Value)
	case token.GEQ:
		return convertToBooleanLiteral(leftObj.Value >= rightObj.Value)
	default:
		return nil
	}
}

func evaluateBinaryOperationOfFloatingPointLiteral(
	leftObj *object.FloatingPointLiteral,
	operator token.Token,
	rightObj *object.FloatingPointLiteral,
) object.Object {
	switch operator {
	case token.ADD:
		return &object.FloatingPointLiteral{Value: leftObj.Value + rightObj.Value}
	case token.SUB:
		return &object.FloatingPointLiteral{Value: leftObj.Value - rightObj.Value}
	case token.MUL:
		return &object.FloatingPointLiteral{Value: leftObj.Value * rightObj.Value}
	case token.QUO:
		if rightObj.Value == 0 {
			return nil
		}
		return &object.FloatingPointLiteral{Value: leftObj.Value / rightObj.Value}
	case token.LSS:
		return convertToBooleanLiteral(leftObj.Value < rightObj.Value)
	case token.GTR:
		return convertToBooleanLiteral(leftObj.Value > rightObj.Value)
	case token.LEQ:
		return convertToBooleanLiteral(leftObj.Value <= rightObj.Value)
	case token.GEQ:
		return convertToBooleanLiteral(leftObj.Value >= rightObj.Value)
	default:
		return nil
	}
}

func evaluateUnaryOperation(expr *ast.UnaryExpr) object.Object {
	switch expr.Op {
	case token.SUB:
		return evaluateMinusOperation(expr)
	case token.NOT:
		return evaluateNotOperation(expr)
	default:
		return nil
	}
}

func evaluateMinusOperation(expr *ast.UnaryExpr) object.Object {
	obj := evaluateExpression(expr.X)
	intLiteral, ok := obj.(*object.IntegerLiteral)
	if !ok {
		return nil
	}

	intLiteral.Value *= -1
	return intLiteral
}

func evaluateNotOperation(expr *ast.UnaryExpr) object.Object {
	obj := evaluateExpression(expr.X)
	boolLiteral, ok := obj.(*object.BooleanLiteral)
	if !ok {
		return nil
	}

	if boolLiteral == object.True {
		return object.False
	}

	return object.True
}

func evaluateIdentifier(expr *ast.Ident) object.Object {
	obj, ok := object.Env.Get(expr.Name)
	if !ok {
		return nil
	}

	return obj
}

func evaluateBasicLiteral(expr *ast.BasicLit) object.Object {
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

func evaluateIntegerLiteral(expr *ast.BasicLit) object.Object {
	value, err := strconv.Atoi(expr.Value)
	if err != nil {
		return nil
	}
	return &object.IntegerLiteral{
		Value: value,
	}
}

func evaluateStringLiteral(expr *ast.BasicLit) object.Object {
	return &object.StringLiteral{
		Value: expr.Value[1 : len(expr.Value)-1],
	}
}

func evaluateCharacterLiteral(expr *ast.BasicLit) object.Object {
	return &object.CharacterLiteral{
		Value: []rune(expr.Value[1 : len(expr.Value)-1])[0],
	}
}

func evaluateFloatingPointLiteral(expr *ast.BasicLit) object.Object {
	value, err := strconv.ParseFloat(expr.Value, 32)
	if err != nil {
		return nil
	}
	return &object.FloatingPointLiteral{
		Value: float32(value),
	}
}

func convertToBooleanLiteral(b bool) object.Object {
	obj, _ := object.Env.Get(fmt.Sprintf("%t", b))
	return obj
}
