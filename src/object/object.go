package object

type ObjectType string

/*
 * 実行時の型
 */
const (
	NULL_OBJ         ObjectType = "NULL"
	ERROR_OBJ        ObjectType = "ERROR"
	INTEGER_OBJ      ObjectType = "INTEGER"
	FLOAT_OBJ        ObjectType = "FLOAT"
	COMPLEX_OBJ      ObjectType = "COMPLEX"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	STRING_OBJ       ObjectType = "STRING"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
	HASH_OBJ         ObjectType = "HASH"
	TYPE_OBJ         ObjectType = "TYPE"
	CLASS_OBJ        ObjectType = "CLASS"
	BREAK_OBJ        ObjectType = "BREAK"
	CONTINUE_OBJ     ObjectType = "CONTINUE"
)

var (
	NULL      = &Null{}
	TRUE      = &Boolean{Value: true}
	FALSE     = &Boolean{Value: false}
	UNDEFINED = &Undefined{}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

/*
 * エラー
 */
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
