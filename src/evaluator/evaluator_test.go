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
		// mutはオーバライド可能
		mut f2 = ()=>{
			puts("f2");
		}
		return this;
	}

	imm ClassA = ()=>{
		... ClassBase()
		imm parentF2 = f2;

		mut a1 = {
			"hoge":"HOGE",
		};
		a1["gabu"] = "GABU";

		imm f2 = ()=>{
			f1();
			parentF2();
			puts("constructor");
			puts(a1["gabu"]);
			puts(base);
		}
		return this;
	};
	imm a = ClassA();
	a.f2()
`
	input = `
	imm ClassA = ()=>{
		mut hoge=1;
		return this;
	}
	imm ClassB = ()=>{
		...ClassA()
		mut gabu=2;
		return this;
	}
	imm ClassC = ()=>{
		...ClassB()
		mut nano=3;
		return this;
	}
	imm ClassD = ()=>{
		mut desu=4;
		return this;
	}
	imm ClassE = ()=>{
		...ClassA()
		...ClassD()
		mut neow=5;
		return this;
	}
		
	imm a= ClassA();
	imm b= ClassB();
	imm c= ClassC();
	imm d= ClassD();
	imm e= ClassE();
	puts("hello" instanceof string)
	puts("------------- t/f/f/f/f")
	puts(a instanceof ClassA)
	puts(a instanceof ClassB)
	puts(a instanceof ClassC)
	puts(a instanceof ClassD)
	puts(a instanceof ClassE)

	puts("------------- t/t/f/f/f")
	puts(b instanceof ClassA)
	puts(b instanceof ClassB)
	puts(b instanceof ClassC)
	puts(b instanceof ClassD)
	puts(b instanceof ClassE)

	puts("------------- t/t/t/f/f")
	puts(c instanceof ClassA)
	puts(c instanceof ClassB)
	puts(c instanceof ClassC)
	puts(c instanceof ClassD)
	puts(c instanceof ClassE)

	puts("------------- f/f/f/t/f")
	puts(d instanceof ClassA)
	puts(d instanceof ClassB)
	puts(d instanceof ClassC)
	puts(d instanceof ClassD)
	puts(d instanceof ClassE)

	puts("------------- t/f/f/t/t")
	puts(e instanceof ClassA)
	puts(e instanceof ClassB)
	puts(e instanceof ClassC)
	puts(e instanceof ClassD)
	puts(e instanceof ClassE)
`

	input = `
		imm hash = {
			a:"hoge",
			b:"gabu",
			c:"nano",
			d:"desu",
		}


	 	loop(mut i=hash){
			if (i.v=="nano") break
			puts(i)
		}
		puts("----------")
	 	loop(mut i=hash){
			if (i.v=="nano") continue
			puts(i)
		}
		puts("----------")
		imm a = ()=>{
			puts("loop start")
			loop(mut i=hash){
				if (i.v=="nano") return i.v
				puts(i)
			}
			puts("loop end")
		}()
		puts(a)
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
