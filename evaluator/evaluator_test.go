package evaluator

import (
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
				&object.BooleanLiteral{
					Value: true,
				},
				&object.BooleanLiteral{
					Value: false,
				},
				&object.BooleanLiteral{
					Value: true,
				},
				&object.BooleanLiteral{
					Value: false,
				},
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
			`var a, b = true, false`,
			[]object.Object{
				&object.BooleanLiteral{
					Value: true,
				},
				&object.BooleanLiteral{
					Value: false,
				},
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
