package evaluator

import (
	"fmt"
	"monkey/object"
)

/*
 * エラーオブジェクトを返す
 */
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

/*
 * エラーオブジェクトかどうか調べる
 */
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
