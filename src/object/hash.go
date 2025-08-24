package object

import (
	"bytes"
	"fmt"
	"monkey/lib"
	"strings"
)

/*
 *  ハッシュ処理の結果オブジェクト
 */
type HashError struct {
	error   string
	message string
}

func (hr *HashError) clone(format string, args ...interface{}) *HashError {
	return &HashError{error: hr.error, message: fmt.Sprintf(format, args...)}
}
func (hr *HashError) Is(he *HashError) bool {
	return hr.error == he.error
}

func (hr *HashError) Error() string {
	return hr.message
}

var (
	NotFound   *HashError = &HashError{error: "NotFound"}
	InvalidKey *HashError = &HashError{error: "InvalidKey"}
)

/*
 *	ハッシュキーの型
 */
type HashKey struct {
	Type  ObjectType
	Value uint64
}

/*
 * インタフェイス
 */
type Hashable interface {
	HashKey() HashKey
}

/*
 * Key-Valueのペア
 */
type HashPair struct {
	key   Object
	value Object
}

func (hp *HashPair) Key() Object {
	return hp.key
}
func (hp *HashPair) Value() Object {
	return hp.value
}

/*
 * 順序付きマップ
 */
type Hash struct {
	pairs *lib.OrderedMap[HashKey, *HashPair]
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}

	h.pairs.Range(func(k HashKey, v *HashPair) bool {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			v.key.Inspect(), v.value.Inspect()))
		return true
	})

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// 新規
func NewHash() *Hash {
	return &Hash{
		pairs: lib.New[HashKey, *HashPair](),
	}
}

// 保存
func (h *Hash) Set(key Object, value Object) *HashError {

	if k, ok := key.(Hashable); ok {
		h.pairs.Set(k.HashKey(), &HashPair{key: key, value: value})
	}
	return InvalidKey.clone("key is not hashable: %#v", key)
}

// 取得
func (h *Hash) Get(key Object) (Object, *HashError) {
	if k, ok := key.(Hashable); ok {
		if hp, ok := h.pairs.Get(k.HashKey()); ok {
			return hp.value, nil
		}
		return nil, NotFound.clone("key not found: %#v", key)
	}
	return nil, InvalidKey.clone("key is not hashable: %#v", key)
}

// 削除
func (h *Hash) Delete(key Object) *HashError {
	if k, ok := key.(Hashable); ok {
		if ok := h.pairs.Delete(k.HashKey()); ok {
			return nil
		}
		return NotFound.clone("key not found: %#v", key)
	}
	return InvalidKey.clone("key is not hashable: %#v", key)

}

// 別のハッシュをマージする
func (h *Hash) Merge(source *Hash) *Hash {

	source.pairs.Range(func(k HashKey, v *HashPair) bool {
		keyObj := v.Key()
		if keyObj.Type() == STRING_OBJ {
			h.Set(keyObj, v.Value())
		}
		return true
	})
	return h
}

// ループ
func (h *Hash) Range(fn func(key *Object, val *Object) bool) {
	h.pairs.Range(func(k HashKey, p *HashPair) bool {
		return fn(&p.key, &p.value)
	})
}
