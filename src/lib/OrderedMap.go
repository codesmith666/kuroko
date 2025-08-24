package lib

import "container/list"

type OrderedMap[K comparable, V any] struct {
	m   map[K]V
	pos map[K]*list.Element
	ll  *list.List // 各要素は K（キー）
}

func New[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		m:   make(map[K]V),             // map
		pos: make(map[K]*list.Element), //
		ll:  list.New()}
}

/*
 * 値を設定
 */
func (o *OrderedMap[K, V]) Set(k K, v V) {
	// キーが存在したら値だけ更新
	if _, ok := o.pos[k]; ok {
		o.m[k] = v
		return
	}
	// キーがないので新規に保存
	o.m[k] = v
	o.pos[k] = o.ll.PushBack(k)
}

/*
 * 値を取得
 */
func (o *OrderedMap[K, V]) Get(k K) (V, bool) {
	v, ok := o.m[k]
	return v, ok
}

/*
 * 値を削除
 */
func (o *OrderedMap[K, V]) Delete(k K) bool {
	if e, ok := o.pos[k]; ok {
		o.ll.Remove(e)
		delete(o.pos, k)
		delete(o.m, k)
		return true
	}
	return false
}

/*
 * 長さ（大きさ）
 */
func (o *OrderedMap[K, V]) Len() int {
	return len(o.m)
}

// Range は挿入順で反復。fn が false を返すと中断。
func (o *OrderedMap[K, V]) Range(fn func(K, V) bool) {
	for e := o.ll.Front(); e != nil; e = e.Next() {
		k := e.Value.(K)
		if !fn(k, o.m[k]) {
			return
		}
	}
}
