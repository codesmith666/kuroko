package parser

import (
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {

	var input = `
	if ( 1 ) { 2 }
	if ( 1 ) 2
	if ( 1 ) { 2 } else { 3 }
	if ( 1 )  2 else 3
	if ( 1 ) { 2; } else { 3; }
	if ( 1 )  2; else 3;
	if ( 1 ) { 2 } else 3
	if ( 1 ) 2 else { 3 }
	if ( 1 ) 2
	if ( 1 ) 2 else if ( 3 ){ 4 } else if (5) {6}
	if ( 1+2*3/4 )
`

	// imm b:number = 1;
	// imm a:number[] = [];
	// imm a:{[key:string]:number} = {};
	// imm a:{hoge:string,gabu:number} = {};
	/*
		imm a:(num:number,str:string)=>(a:number)=>number = (num:number,str:string)=>{
			mut b="hogehoge"+"gabugabu"
			return (a:number)=>{return a+16};
		}
		mut b = 1;

	*/

	input = `
		/*
		 * 関数が呼び出されるとオブジェクトを返す
		 * ~~~
		 *   hoge
		 * ~~~
		 */
		imm TestClassA = ()=>{
			imm a=0;
			return this
		}
		// クラスAから派生したクラスB
		// ...演算子でAの内容をBのオブジェクトにマッピングしている
		// ~~~
		//   hoge
		// ~~~
		imm TestClassB = ():Object=>{
			...TestClassA()
			return this
		}
		// クラス名は関数名になる
		imm b = TestClassB;
		imm b1 = b();
		imm b2 = TestClassB();
 		// "b"
		puts(b1.$name);
 		// "TestClassB"
		puts(b2.$name);
`
	input = `
	mut a
`

	p := NewParser(input)
	p.DumpTokens()
	a, _ := p.ParseProgram()

	fmt.Print("--------------------------------------\n")
	fmt.Printf("%s\n", input)
	fmt.Print("--------------------------------------\n")
	fmt.Printf("%s\n", a.String())
}
