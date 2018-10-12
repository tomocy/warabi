package object

var Env = &Environment{
	objs: map[string]Object{
		"true":  True,
		"false": False,
	},
}

var builtins = map[string]bool{
	"true":  true,
	"false": true,
}

type Environment struct {
	objs map[string]Object
}

func NewEnvironment() *Environment {
	return &Environment{
		objs: make(map[string]Object),
	}
}

func (e *Environment) Set(name string, obj Object) {
	if builtins[name] {
		return
	}
	e.objs[name] = obj
}

func (e Environment) Get(name string) (Object, bool) {
	obj, ok := e.objs[name]
	return obj, ok
}
