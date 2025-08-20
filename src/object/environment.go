package object

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

type Environment struct {
	// class map[string]Object
	class *Class
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{class: NewClass(), outer: nil}
}

// Get
//
// thisが要求されたときは、 e.class を返す。
// これは Hashと互換性があるので Hashのようにアクセスで、
// つまり this として機能する。
// このスコープにnameが無かったら外側のスコープを探しに行く。
func (e *Environment) Get(name string) Object {
	if name == "this" {
		return e.class
	}
	val, err := e.class.Hash.Get(&String{Value: name})
	if err != nil && e.outer != nil {
		return e.outer.Get(name)
	}
	return val
}

func (e *Environment) Set(name string, val Object) Object {
	key := &String{Value: name}
	e.class.Set(key, val)
	return val
}

// ハッシュで環境を派生させる
// ... ステートメントで実行される。
func (e *Environment) DeriveFromHash(from *Hash) {
	// 子の環境をすべて自分のものとして取得する
	// ただし名前でアクセスできるもののみ
	// Hashの順序は守られる
	e.class.Hash.Merge(from)
}

// クラス情報で環境を派生させる
// ... ステートメントで実行される。
func (e *Environment) DeriveFromClass(from *Class) {
	// ハッシュ部分を派生
	e.class.Derive(from)
}
