package object

/*
 * 型
 */
type Type struct {
	Name string
}

func (t *Type) Type() ObjectType { return TYPE_OBJ }
func (t *Type) Inspect() string  { return t.Name }
