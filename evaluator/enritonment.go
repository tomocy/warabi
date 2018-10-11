package evaluator

import "github.com/tomocy/warabi/object"

var env = &environment{
	objs: map[string]object.Object{
		"true": &object.BooleanLiteral{
			Value: true,
		},
		"false": &object.BooleanLiteral{
			Value: false,
		},
	},
}

type environment struct {
	objs map[string]object.Object
}

func (e *environment) set(name string, obj object.Object) {
	e.objs[name] = obj
}

func (e environment) get(name string) (object.Object, bool) {
	obj, ok := e.objs[name]
	return obj, ok
}
