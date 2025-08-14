

```

//
// 型はコロン後に後置で宣言
//  number/string/array/object/any
//
@imm stdlog = (key:string,val:number)=>{}

//
// @imm
//  イミュータブルな宣言。この変数は宣言後に１度だけ書き込める。
// @mut
//  ミュータブルな宣言。この変数は何度でも書き込める。
//  ループカウンタとかステータス管理とかに使うだろう。
//  関数宣言をmutにするとオーバーライドできる。
//
@imm member = (init:number)=>{
    @mut str = "a,b,c";
    @mut num = init;
    @imm alpha = ()=>{
        stdlog(str,num);
    }
    @mut beta = ()=>{
        stdlog("beta");
    }
}

//
// 多態性
//  多重継承ができる（非推奨）
//  変数や関数が衝突したときはエラー
//  下の例では、betaをオーバライドしてgammaとnameを追加している。
//  @fromをここに書きたくない
//
@imm user = (init:number,name:string)=>{
    @from member(init)
    @imm name = name+"@hoge"
    @imm gamma = ()=>{
        stdlog("GAMMA")
    }
    @imm beta = ()=>{
        stdlog("BETA")
    }
    // 処理を書けばコンストラクタになる
    //  colorはuserのメンバだが、colはメンバにならない。
    @mut color;
    {
        @imm col = math.randomColor();
        color=col;   // 
    }
}

//
// アヒル型
//  ダックタイピングはphpほど緩くはなくgoの様に宣言的。
//
@imm duck = ()=>{
    sound(count:number):string{
        return "quack!" * count;
    }
}
@imm cat = ()=>{
    sound(count:number):string{
        return "meow!" * count;
    }
}
@type sounds = {
    sound(count:number):string;
}

imm sound = (objA:sounds)=>{
    stdlog(objA.sounds());
}
sound(duck())
sound(cat())

//
//  型
//
@imm codes = @imm {
    foo:"Foo",
    bar:"Bar",
    baz:"Baz",
};

@type CodeKey = @typeof codes[foo]; // "foo"
@type CodeKey = @typeof @keyof codes;     // "foo"|"bar"|"baz"
@type CodeVal = @typeof codes[CodeKey]; // "Foo"|"Bar"|"Baz"
@type StrList = @typeof string[]

@imm a:string[] = [];


//
// undefined
//  goで言うゼロ値はない。
//  未初期化の変数はすべてundefined
//  存在しないインデックスにアクセスしてもundefined
//

//
// 条件分岐
//
@if isLegal(){
    stdlog("Legal");
}
@elif isIllegal(){
    stdlog("Illegal");
}
@else{
    stdlog("unknown");
}

//
// switch
//
@switch(hoge){
    @case "a":
        @break;
    @case "b":
        @break;
    @case "c":
        @break;
}

//
// ループ
//
@for(array as a){
    @if ( a=='1' ) @continue;
    @if ( a=='2' ) @break;
    stdlog(a);
}
a.each((v)=>{
    stdlog(v);
})
a.each((v,k)=>{
    stdlog(`${k}=${v}`);
})

//
// オブジェクト
//  オブジェクトに入れた順序は守られる(php風)
imm a = {};
a["a"] = 1;
a["b"] = 2;
a["c"] = 3;
a.each((v,k)=>{
    stdlog(`${k}=${v},`); // a=1,b=2,c=3,
})




//
// りてらる
//
0 整数
0.1 少数
0x00000000 16進数
0o777 8進数
0b010101010101 2進数
2+3i 複素数

//
//  stringを拡張できるか->できない（下の２つの処理が混同してしまう）
//  id = "foo"          // idに文字列を入れたい
//  id = @new id("bar") // idに新しいidオブジェクトを入れたい
//
//  
//

// thisはあるのか？
//  あるが使わなくてもいい

//
// 値の参照や設定はフックできる？
//  typescript風にできる
//  オブジェクトの変数でなくてもできる。
//  @getと@setはイミュータブル
//
@mut _id = "hoge"
@get id = ()=>{
    return _id;
}
@set id = (id:string)=>{
    _id = id;
}

//
// @constは、@imm @static のシンタックスシュガー
// @shareは、@mut @static のシンタックスシュガー
//
@imm Member = () => {

    @mut a = 1;
    a++;

    @imm b = a;

    @const HOGE=a;
    @const FOO = ()=>{
        @return "foo";
    }
    @share hoge=2;
}

//
// private / protected / publicはどうする？
// 決めてない。
//


//
// オブジェクトの生成（@newキーワード）
// どんな関数もnewできる。逆にnewできるのは関数のみ。
// pythonのように@newなしでは生成できない
// オーバーライドしながらnewできる。
// 
@imm member1 = @new member(5);
@imm member2 = @new member(5){
    @mut beta = ()=>{
        errlog("Beta")
    }
}
member1.alpha();
member2.beta();  // "Beta" (not "beta")
stdlog(member.@name); // member
stdlog(user.@parent.@name); // member

//
// 変数は宣言と同時に初期化（束縛）が必要
//  @imm @mut @get @put @const @share
//

//
// 予約語には必ず@がつく
// そのため ifとかreturnとかいう変数名もOK
//



```

```

immutableはpublic、どうせ書き換えられないから。
mutableはprivate、書き換えるにはメソッドを経由しないといけない

//
// Number型にメソッドを持ちたい（けどMathはあるよな）
//
100.random(-2)      -> 0.00~99.99の乱数を生成
1.02123.ceil();     -> 1.1
1.02123.floor();    -> 1.0
1.02123.round();    -> 1.0
1.05.round(2);      -> 1.1
2.sqrt();           -> 1.41421356
-2.abs();           -> 2

5..10               -> [5,6,7,8,9,10]

imm Math = {
    const PI = 3.14159
    const ceil = (x:number)=>{
        if (Number.isInteger(x)) return x;
        return x > 0 ? parseInt(x) + 1 : parseInt(x);
    }
}

stdlog(Math.PI)
stdlog(Math.ceil(0.1))

//
// 連想配列
//
imm assoc = { a:1, b:2 }
と
imm assoc = { mut a=1;mut a=2; }
は等価

//
// オブジェクトを実行=コンストラクタ
//
imm Obj = (n1:number,n2:number)=>{
    mut a:number = n1;
    imm i:number = n2;
    mut foo:number=100;

    // オブジェクトの文字列表記を返す
    imm <string = ()=>{
        return `{a=${this.a},b=${this.b}}`;
    }
    // ↓やらない（関数を自前で書く）
    imm <number = ()=>{
        return a*i;
    }
    // 値の参照を監視
    imm <foo = ()=>{
        stdlog("<foo",v);
        return this.foo;
    }
    // 値の書き込みを監視（や変換）
    imm >foo = (v:number)=>{
        stdlog(">foo",v)
        return v;
    }
    // 普通の関数
    mut product = ()=>{
        return a*i;
    }
    return this
}
// オブジェクトを作成してダンプ
imm obj = Obj(1)        // $invoke()が呼び出される Obj.$invoke
stdlog(obj)             // オブジェクトの文字列表記を返す
stdlog(obj.product()) 
obj.foo;
obj.foo=0; 
// オブジェクトを作成しつつ監視関数を追加
imm obj2 = Obj(1){
    imm <a = ()=>{
        stdlog("<a",this.a)
    }
}
// 多態性
imm Obj2 = (n1:number,n2:number,n3:number)=>{
    ...Obj(n1,n2);                  //  ここで派生（関数読んでオブジェクトを生成しマージ）
    imm objProduct = product;       //  オーバライドの準備（オプション）
    mut c = n3;
    // mutならオーバーライドできる
    mut product(){
        return objProduct() * c;    // クロージャでcを参照できる
    }

    return this;
};
// Obj2のthisをオーバーライドしながら初期化
imm obj2 = Obj2(1,2,3){
    imm obj2Product = product;
    mut product(){
        return obj2Product()*4;
    }
}

// この表記でちゃんと型を解決できるのか
typeof Obj;
->
{
    a:number;
    i:number;
    foo:number;
    <string:()=>string;
    <foo:()=>number;
    >foo:(number)=>void;
    product:()=>number;
}

typeof Obj2
->
{
    a:number;
    i:number;
    foo:number;
    <string:()=>string;
    <foo:()=>number;
    >foo:(number)=>void;
    product:()=>number;
    c:number;
    objProduct:()=>number;
}


imm Hoge = (num:number)=>{

    const FOO = 10;             // static immutable number
    share bar = 10;             // static mutable number
    imm baz = Math.random(num); // immutable number;
    mut qux = "DEF";            // mutable string

    imm test = ()=>{
        mut i=0;
        stdlog("test"+i)
    }

    return this;
}

Hoge(1) 


```


