package object

/*
 * 戻り値
 */
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

/*
 * break
 */
type Break struct{}

func (rv *Break) Type() ObjectType { return BREAK_OBJ }
func (rv *Break) Inspect() string  { return "break" }

/*
 * continue
 */
type Continue struct{}

func (rv *Continue) Type() ObjectType { return CONTINUE_OBJ }
func (rv *Continue) Inspect() string  { return "continue" }
