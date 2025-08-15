package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	// s := make(map[string]Object)
	p := make(map[HashKey]HashPair)
	s := &Class{Hash: Hash{Pairs: p}, Name: "$unnamed"}
	return &Environment{class: s, outer: nil}
}

type Environment struct {
	// class map[string]Object
	class *Class
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {

	// thisが要求されたときは、 e.store を返す
	// これｈHashと互換性があると同時に
	// thisをreturnする関数はコンストラクタのようにふるまう
	if name == "this" {
		return e.class, true
	}

	// obj := e.store[name]
	key := &String{Value: name}
	pair, ok := e.class.Pairs[key.HashKey()]
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
	e.class.Pairs[key.HashKey()] = HashPair{Key: key, Value: val}
	return val
}

// ハッシュで環境を派生させる
// ... ステートメントで実行される。
func (e *Environment) DeriveFromHash(from *Hash) {
	// 子の環境をすべて自分のものとして取得する
	// ただし名前でアクセスできるもののみ
	for _, pair := range from.Pairs {
		if pair.Key.Type() == STRING_OBJ {
			e.Set(pair.Key.(*String).Value, pair.Value)
		}
	}
}

// クラス情報で環境を派生させる
// ... ステートメントで実行される。
func (e *Environment) DeriveFromClass(from *Class) {
	// ハッシュ部分を派生
	e.DeriveFromHash(&from.Hash)

	// クラスなら子クラス名をすべて引き継ぐ
	// 領域確保
	if e.class.Children == nil {
		e.class.Children = make(map[string]struct{})
	}
	// 保存
	if from.Children != nil {
		for childName, _ := range from.Children {
			e.class.Children[childName] = struct{}{}
		}
	}
	//
	e.class.Children[from.Name] = struct{}{}

}
