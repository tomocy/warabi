package object

var Env = &environment{
	objs: map[string]Object{
		"true":  True,
		"false": False,
	},
}

var builtins = map[string]bool{
	"true":  true,
	"false": true,
}

type environment struct {
	objs map[string]Object
}

func (e *environment) Set(name string, obj Object) {
	if builtins[name] {
		return
	}
	e.objs[name] = obj
}

func (e environment) Get(name string) (Object, bool) {
	obj, ok := e.objs[name]
	return obj, ok
}
