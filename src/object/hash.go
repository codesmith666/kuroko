package object

import (
	"bytes"
	"fmt"
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
	order []HashKey
	pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Items() {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.key.Inspect(), pair.value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// 新規
func NewHash() *Hash {
	return &Hash{
		pairs: make(map[HashKey]HashPair),
		order: []HashKey{},
	}
}

// 保存
func (h *Hash) Set(key Object, value Object) *HashError {
	if k, ok := key.(Hashable); ok {
		hash := k.HashKey()
		if _, exists := h.pairs[hash]; !exists {
			h.order = append(h.order, hash)
		}
		h.pairs[hash] = HashPair{key: key, value: value}
		return nil
	}
	return InvalidKey.clone("key is not hashable: %#v", key)
}

// 取得
func (h *Hash) Get(key Object) (Object, *HashError) {
	if k, ok := key.(Hashable); ok {
		if hashPair, ok := h.pairs[k.HashKey()]; ok {
			return hashPair.value, nil
		} else {
			return nil, NotFound.clone("key not found: %#v", key)
		}
	}
	return nil, InvalidKey.clone("key is not hashable: %#v", key)
}

// 削除
func (h *Hash) Delete(key Object) *HashError {
	if k, ok := key.(Hashable); ok {
		hash := k.HashKey()
		if _, exists := h.pairs[hash]; !exists {
			return NotFound.clone("key not found: %#v", key)
		}
		delete(h.pairs, hash)

		// keys からも削除（新しい一覧を作って保存する）
		newKeys := make([]HashKey, 0, len(h.order)-1)
		for _, k := range h.order {
			if k != hash {
				newKeys = append(newKeys, hash)
			}
		}
		h.order = newKeys
		return nil
	}
	return InvalidKey.clone("key is not hashable: %#v", key)

}

// 順番の正しいHashPairの一覧を返す（for用）
func (h *Hash) Items() []HashPair {
	result := []HashPair{}
	for _, k := range h.order {
		result = append(result, h.pairs[k])
	}
	return result
}

// 別のハッシュをマージする
func (h *Hash) Merge(source *Hash) *Hash {
	for _, hashPair := range source.Items() {
		keyObj := hashPair.Key()
		if keyObj.Type() == STRING_OBJ {
			h.Set(keyObj, hashPair.Value())
		}
	}
	return h
}
