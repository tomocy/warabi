package evaluator

import "github.com/tomocy/warabi/object"

var env = &environment{
	objs: map[string]object.Object{
		"true":  object.True,
		"false": object.False,
	},
}

var builtins = map[string]bool{
	"true":  true,
	"false": true,
}

type environment struct {
	objs map[string]object.Object
}

func (e *environment) set(name string, obj object.Object) {
	if builtins[name] {
		return
	}
	e.objs[name] = obj
}

func (e environment) get(name string) (object.Object, bool) {
	obj, ok := e.objs[name]
	return obj, ok
}
