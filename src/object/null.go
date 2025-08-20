package object

/*
 * null
 */
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

/*
 * Undefined
 */
type Undefined struct{}

func (n *Undefined) Type() ObjectType { return NULL_OBJ }
func (n *Undefined) Inspect() string  { return "undefined" }
