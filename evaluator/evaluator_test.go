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
			"var a, b = 10 + 5, 5 - 10",
			[]object.Object{
				&object.IntegerLiteral{
					Value: 15,
				},
				&object.IntegerLiteral{
					Value: -5,
				},
			},
		},
		{
			"var a, b, c = 5 * -5, 5 / 5, 10 % 5",
			[]object.Object{
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
			`var a = "go"`,
			[]object.Object{
				&object.StringLiteral{
					Value: `"go"`,
				},
			},
		},
		{
			`var a = "hello, " + "world"`,
			[]object.Object{
				&object.StringLiteral{
					Value: `"hello, world"`,
				},
			},
		},
	}

	for _, test := range tests {
		gots := Evaluate(test.source)
		if len(gots) != len(test.wants) {
			t.Fatalf("unexpected object length: got %d, expected %d\n", len(gots), len(test.wants))
		}
		for i := 0; i < len(test.wants); i++ {
			if !reflect.DeepEqual(gots[i], test.wants[i]) {
				t.Errorf("unexpected object: got %#v, expected %#v\n", gots[i], test.wants[i])
			}
		}
	}
}
