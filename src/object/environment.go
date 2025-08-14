package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	// s := make(map[string]Object)
	p := make(map[HashKey]HashPair)
	s := &Hash{Pairs: p}
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	// store map[string]Object
	store *Hash
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {

	if name == "this" {
		return e.store, true
	}

	// obj := e.store[name]
	key := &String{Value: name}
	pair, ok := e.store.Pairs[key.HashKey()]
	obj := pair.Value

	// outerがあれば。
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	key := &String{Value: name}
	//e.store[name] = val
	e.store.Pairs[key.HashKey()] = HashPair{Key: key, Value: val}
	return val
}
