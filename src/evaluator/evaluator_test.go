package evaluator

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvaluator(t *testing.T) {

	var input = ""
	input = `
	// thisを返すのはコンストラクタ
	imm ClassA = ()=>{
		mut a1 = {
			hoge:"HOGE"
		};
		return this
	};
	imm ClassB = ()=>{
		mut b1 = 3;
		return this;
	}

	// AとBから派生したCというきもいオブジェクト
	// 多重継承は多くの言語で否定されているが・・・
	imm ClassC = ()=>{
		...ClassA()
		...ClassB()
		mut c1 = 5;

		imm sum = ()=>{
			return a1 + b1 + c1;
		}

		return this;
	}
	// 単なる関数
	// ClassBから派生したClassCをnewしてcに代入
	imm c=ClassC();

	puts(c.sum());
	`

	input = `
		mut a={b:1};
		a.b={c:2};
		puts(a.b.c);
`

	input = `
	imm ClassBase = ()=>{
		imm base=1234;
		imm f1 = ()=>{
			puts("f1");
		}
		//return this;
	}

	imm ClassA = ()=>{
		... ClassBase()
		mut a1 = {
			"hoge":"HOGE",
		};
		a1["gabu"] = "GABU";

		imm f2 = ()=>{
			f1();
			puts("constructor");
			puts(a1["gabu"]);
			puts(base);
		}
		return this;
	};
	imm a = ClassA();
	a.f2()

`

	env := object.NewEnvironment()
	p := parser.NewParser(input)
	a, ok := p.ParseProgram()
	ast.PrintAST(a, "")
	if !ok {
		p.OutputErrors()
	}

	e := Eval(a, env)
	if e != nil {
		var out bytes.Buffer
		out.WriteString("eval result: ")
		out.WriteString(e.Inspect())
		out.WriteString("\n")
		fmt.Printf("%s", out.String())
	}
}
