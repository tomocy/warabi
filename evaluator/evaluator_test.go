package evaluator

import (
	"go/ast"
	"go/token"
	"reflect"
	"testing"

	"github.com/tomocy/warabi/object"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		source string
		wants  []object.Object
	}{
		{
			"var a int",
			[]object.Object{
				&object.IntegerLiteral{
					Value: 0,
				},
			},
		},
		{
			"var a int = 10",
			[]object.Object{
				&object.IntegerLiteral{
					Value: 10,
				},
			},
		},
		{
			"var a, b, c, d, e = 10 + 5, 5 - 10, 5 * -5, 5 / 5, 10 % 5",
			[]object.Object{
				&object.IntegerLiteral{
					Value: 15,
				},
				&object.IntegerLiteral{
					Value: -5,
				},
				&object.IntegerLiteral{
					Value: -25,
				},
				&object.IntegerLiteral{
					Value: 1,
				},
				&object.IntegerLiteral{
					Value: 0,
				},
			},
		},
		{
			"var a = 5 * (1 + 3)",
			[]object.Object{
				&object.IntegerLiteral{
					Value: 20,
				},
			},
		},
		{
			"var a, b, c, d = 5 < 10, 1 > 1, 0 <= 0, 98 >= 99",
			[]object.Object{
				object.True,
				object.False,
				object.True,
				object.False,
			},
		},
		{
			`var a = "go"`,
			[]object.Object{
				&object.StringLiteral{
					Value: "go",
				},
			},
		},
		{
			`var a = "hello, " + "world"`,
			[]object.Object{
				&object.StringLiteral{
					Value: "hello, world",
				},
			},
		},
		{
			`var a, b, c, d = "a" < "a", "a" < "b", "a" <= "a", "a" >= "z"`,
			[]object.Object{
				object.False,
				object.True,
				object.True,
				object.False,
			},
		},
		{
			`var a = 'a'`,
			[]object.Object{
				&object.CharacterLiteral{
					Value: 'a',
				},
			},
		},
		{
			`var a, b, c, d, e = 'a' + 'a', 'b' - 'b', 'c' * 'c', 'd' / 'd', 'e' % 'e'`,
			[]object.Object{
				&object.CharacterLiteral{
					Value: 'a' + 'a',
				},
				&object.CharacterLiteral{
					Value: 'b' - 'b',
				},
				&object.CharacterLiteral{
					Value: 'c' * 'c',
				},
				&object.CharacterLiteral{
					Value: 'd' / 'd',
				},
				&object.CharacterLiteral{
					Value: 'e' % 'e',
				},
			},
		},
		{
			`var a, b, c, d = 'a' < 'a', 'a' < 'b', 'a' <= 'a', 'a' >= 'z'`,
			[]object.Object{
				object.False,
				object.True,
				object.True,
				object.False,
			},
		},
		{
			`var a, b, c, d = 5.5 + 2, 4 - 2.0, 6 * 0.0, 3 / 2.0`,
			[]object.Object{
				&object.FloatingPointLiteral{
					Value: 7.5,
				},
				&object.FloatingPointLiteral{
					Value: 2.0,
				},
				&object.FloatingPointLiteral{
					Value: 0.0,
				},
				&object.FloatingPointLiteral{
					Value: 1.5,
				},
			},
		},
		{
			"var a, b, c, d = 5.0 < 10, 1 > 1.1, 0.0 <= 0, 99.0 >= 99.1",
			[]object.Object{
				object.True,
				object.False,
				object.True,
				object.False,
			},
		},
		{
			"var a, b, c, d = true, false, !false, !true",
			[]object.Object{
				object.True,
				object.False,
				object.True,
				object.False,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.source, func(t *testing.T) {
			gots := Evaluate(test.source)
			if len(gots) != len(test.wants) {
				t.Fatalf("unexpected object length: got %d, expected %d\n", len(gots), len(test.wants))
			}
			for i := 0; i < len(test.wants); i++ {
				if !reflect.DeepEqual(gots[i], test.wants[i]) {
					t.Errorf("unexpected object: got %#v, expected %#v\n", gots[i], test.wants[i])
				}
			}
		})
	}
}

func TestEvaluateFunctionLiteral(t *testing.T) {
	source := `
	func a(b string, c int, d string) (e string, f int) {
		var e = b + d
		var f = c * -1
	}`
	want := &object.FunctionLiteral{
		Params: []*ast.Field{
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "b",
					},
				},
				Type: &ast.Ident{
					Name: "string",
				},
			},
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "c",
					},
				},
				Type: &ast.Ident{
					Name: "int",
				},
			},
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "d",
					},
				},
				Type: &ast.Ident{
					Name: "string",
				},
			},
		},
		Results: []*ast.Field{
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "e",
					},
				},
				Type: &ast.Ident{
					Name: "string",
				},
			},
			&ast.Field{
				Names: []*ast.Ident{
					&ast.Ident{
						Name: "f",
					},
				},
				Type: &ast.Ident{
					Name: "int",
				},
			},
		},
		Body: []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "e",
								},
							},
							Values: []ast.Expr{
								&ast.BinaryExpr{
									X: &ast.Ident{
										Name: "b",
									},
									Op: token.ADD,
									Y: &ast.Ident{
										Name: "c",
									},
								},
							},
						},
					},
				},
			},
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{
								&ast.Ident{
									Name: "f",
								},
							},
							Values: []ast.Expr{
								&ast.BinaryExpr{
									X: &ast.Ident{
										Name: "c",
									},
									Op: token.MUL,
									Y: &ast.UnaryExpr{
										Op: token.SUB,
										X: &ast.BasicLit{
											Kind:  token.INT,
											Value: "1",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	gots := Evaluate(source)
	if len(gots) != 1 {
		t.Fatalf("unexpected object length: got %d, expected 1\n", len(gots))
	}
	got, ok := gots[0].(*object.FunctionLiteral)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *object.FunctionLiteral", gots[0])
	}

	checkFunctionLiteralParams(t, got, want)
	checkFunctionLiteralResults(t, got, want)
	checkFunctionLiteralBody(t, got, want)
}

func checkFunctionLiteralParams(t *testing.T, got *object.FunctionLiteral, want *object.FunctionLiteral) {
	for i := 0; i < len(want.Params); i++ {
		gotTypeIdent := got.Params[i].Type.(*ast.Ident)
		wantTypeIdent := want.Params[i].Type.(*ast.Ident)
		if gotTypeIdent.Name != wantTypeIdent.Name {
			t.Errorf("unexpected param type: got %s, expected %s\n", got.Params[i].Type, want.Params[i].Type)
		}
		for j := 0; j < len(want.Params[i].Names); j++ {
			if got.Params[i].Names[j].Name != want.Params[i].Names[j].Name {
				t.Errorf("unexpected param name: got %s, expected %s\n", want.Params[i].Names[j].Name, got.Params[i].Names[j].Name)
			}
		}
	}
}

func checkFunctionLiteralResults(t *testing.T, got *object.FunctionLiteral, want *object.FunctionLiteral) {
	for i := 0; i < len(want.Results); i++ {
		gotTypeIdent := got.Params[i].Type.(*ast.Ident)
		wantTypeIdent := want.Params[i].Type.(*ast.Ident)
		if gotTypeIdent.Name != wantTypeIdent.Name {
			t.Errorf("unexpected param type: got %s, expected %s\n", gotTypeIdent.Name, wantTypeIdent.Name)
		}
		for j := 0; j < len(want.Results[i].Names); j++ {
			if got.Results[i].Names[j].Name != want.Results[i].Names[j].Name {
				t.Errorf("unexpected param name: got %s, expected %s\n", want.Results[i].Names[j].Name, got.Results[i].Names[j].Name)
			}
		}
	}
}

func checkFunctionLiteralBody(t *testing.T, got *object.FunctionLiteral, want *object.FunctionLiteral) {
	if len(got.Body) != len(want.Body) {
		t.Errorf("unexpected body length: got %d, expected %d\n", len(got.Body), len(want.Body))
	}
	checkFunctionLiteralBodyStatement1(t, got.Body[0], want.Body[0])
	checkFunctionLiteralBodyStatement2(t, got.Body[1], want.Body[1])
}

func checkFunctionLiteralBodyStatement1(t *testing.T, got ast.Stmt, want ast.Stmt) {
	gotDeclStmt, ok := got.(*ast.DeclStmt)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.DeclStmt\n", got)
	}
	wantDeclStmt := want.(*ast.DeclStmt)
	gotGenDecl, ok := gotDeclStmt.Decl.(*ast.GenDecl)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.GenDecl\n", gotDeclStmt.Decl)
	}
	wantGenDecl := wantDeclStmt.Decl.(*ast.GenDecl)
	gotValueSpec, ok := gotGenDecl.Specs[0].(*ast.ValueSpec)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.ValueSpec\n", gotGenDecl.Specs[0])
	}
	wantValueSpec := wantGenDecl.Specs[0].(*ast.ValueSpec)
	for i := 0; i < len(wantValueSpec.Names); i++ {
		if gotValueSpec.Names[i].Name != wantValueSpec.Names[i].Name {
			t.Errorf("unexpected value spec name: got %v, expected %v\n", gotValueSpec.Names[i].Name, wantValueSpec.Names[i].Name)
		}
	}
	gotBinaryExpr, ok := gotValueSpec.Values[0].(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.BinaryExpr\n", gotValueSpec.Values[0])
	}
	wantBinaryExpr := wantValueSpec.Values[0].(*ast.BinaryExpr)
	gotXIdent, ok := gotBinaryExpr.X.(*ast.Ident)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.Ident\n", gotBinaryExpr.X)
	}
	wantXIdent := wantBinaryExpr.X.(*ast.Ident)
	if gotXIdent.Name != wantXIdent.Name {
		t.Errorf("unexpected ident name: got %s, expected %s\n", gotXIdent.Name, wantXIdent.Name)
	}
}

func checkFunctionLiteralBodyStatement2(t *testing.T, got ast.Stmt, want ast.Stmt) {
	gotDeclStmt, ok := got.(*ast.DeclStmt)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.DeclStmt\n", got)
	}
	wantDeclStmt := want.(*ast.DeclStmt)
	gotGenDecl, ok := gotDeclStmt.Decl.(*ast.GenDecl)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.GenDecl\n", gotDeclStmt.Decl)
	}
	wantGenDecl := wantDeclStmt.Decl.(*ast.GenDecl)
	gotValueSpec, ok := gotGenDecl.Specs[0].(*ast.ValueSpec)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.ValueSpec\n", gotGenDecl.Specs[0])
	}
	wantValueSpec := wantGenDecl.Specs[0].(*ast.ValueSpec)
	for i := 0; i < len(wantValueSpec.Names); i++ {
		if gotValueSpec.Names[i].Name != wantValueSpec.Names[i].Name {
			t.Errorf("unexpected value spec name: got %v, expected %v\n", gotValueSpec.Names[i].Name, wantValueSpec.Names[i].Name)
		}
	}
	gotBinaryExpr, ok := gotValueSpec.Values[0].(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.BinaryExpr\n", gotValueSpec.Values[0])
	}
	wantBinaryExpr := wantValueSpec.Values[0].(*ast.BinaryExpr)
	gotXIdent, ok := gotBinaryExpr.X.(*ast.Ident)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.Ident\n", gotBinaryExpr.X)
	}
	wantXIdent := wantBinaryExpr.X.(*ast.Ident)
	if gotXIdent.Name != wantXIdent.Name {
		t.Errorf("unexpected ident name: got %s, expected %s\n", gotXIdent.Name, wantXIdent.Name)
	}
	if gotBinaryExpr.Op != wantBinaryExpr.Op {
		t.Errorf("unexpected operator: got %s, expected %s\n", gotBinaryExpr.Op, wantBinaryExpr.Op)
	}
	gotUnaryExpr, ok := gotBinaryExpr.Y.(*ast.UnaryExpr)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.UnaryExpr\n", gotBinaryExpr.Y)
	}
	wantUnaryExpr := wantBinaryExpr.Y.(*ast.UnaryExpr)
	if gotUnaryExpr.Op != wantUnaryExpr.Op {
		t.Errorf("unexpected operator: got %s, expected %s\n", gotUnaryExpr.Op, wantUnaryExpr.Op)
	}
	gotXBasicLit, ok := gotUnaryExpr.X.(*ast.BasicLit)
	if !ok {
		t.Fatalf("assertion failure: got %T, expected *ast.BasicLit\n", gotUnaryExpr.X)
	}
	wantXBasicLit := wantUnaryExpr.X.(*ast.BasicLit)
	if gotXBasicLit.Value != wantXBasicLit.Value {
		t.Errorf("unexpected operand: got %v, expected %v\n", gotXBasicLit.Value, wantXBasicLit.Value)
	}
}
