package lexer

import (
	"fmt"
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {

	input := `
	hoge foo bar
	true false if else return share const imm mut

	= + - ! * / | & ~ || && ++ -- += -= ? < > <= >= == != <=> =>
	... .. . // ほげほげ
	/* がぶがぶ */ , ; : ( ) { } [ ] '
	//
	/**/

	=+-!* /|&~||&&++--+=-=?<><= >===!=<=>=>
	..... .,;:(){}[]'

	$hoge なのなの

	123 123.456 0x0a 0o777 0b010101 123i 123.456i
	"hogehoge=+*/()がぶがぶreturn" + foo
	"hogeがぶhoge${foo+"hoge"+bar}return${baz}${qux}" + quux
	"
hoge
gab
" + hoge
	foo + "bar${"foo
bar
baz
がぶ"}baz${"です
デス
desu"}"
`
	tests := []struct {
		Type    token.TokenType
		Literal string
		Row     int
		Col     int
	}{
		// identifier
		{token.IDENT, "hoge", 2, 2},
		{token.IDENT, "foo", 2, 7},
		{token.IDENT, "bar", 2, 11},
		// reserved
		{token.TRUE, "true", 3, 2},
		{token.FALSE, "false", 3, 7},
		{token.IF, "if", 3, 13},
		{token.ELSE, "else", 3, 16},
		{token.RETURN, "return", 3, 21},
		{token.SHARE, "share", 3, 28},
		{token.CONST, "const", 3, 34},
		{token.IMM, "imm", 3, 40},
		{token.MUT, "mut", 3, 44},
		// operator
		{token.ASSIGN, "=", 5, 2},
		{token.PLUS, "+", 5, 4},
		{token.MINUS, "-", 5, 6},
		{token.NOT, "!", 5, 8},
		{token.ASTERISK, "*", 5, 10},
		{token.SLASH, "/", 5, 12},
		{token.BIT_OR, "|", 5, 14},
		{token.BIT_AND, "&", 5, 16},
		{token.BIT_NOT, "~", 5, 18},
		{token.OR, "||", 5, 20},
		{token.AND, "&&", 5, 23},
		{token.INC, "++", 5, 26},
		{token.DEC, "--", 5, 29},
		{token.PLUS_ASSIGN, "+=", 5, 32},
		{token.MINUS_ASSIGN, "-=", 5, 35},
		{token.QUESTION, "?", 5, 38},
		{token.LT, "<", 5, 40},
		{token.GT, ">", 5, 42},
		{token.LE, "<=", 5, 44},
		{token.GE, ">=", 5, 47},
		{token.EQ, "==", 5, 50},
		{token.NE, "!=", 5, 53},
		{token.UFO, "<=>", 5, 56},
		{token.ARROW, "=>", 5, 60},
		{token.PARSE, "...", 6, 2},
		{token.RANGE, "..", 6, 6},
		{token.ACCESS, ".", 6, 9},
		{token.LINE_COMMENT, " ほげほげ", 6, 11},
		{token.BLOCK_COMMENT, " がぶがぶ ", 7, 2},
		{token.COMMA, ",", 7, 17},
		{token.SEMICOLON, ";", 7, 19},
		{token.COLON, ":", 7, 21},
		{token.LPAREN, "(", 7, 23},
		{token.RPAREN, ")", 7, 25},
		{token.LBRACE, "{", 7, 27},
		{token.RBRACE, "}", 7, 29},
		{token.LBRACKET, "[", 7, 31},
		{token.RBRACKET, "]", 7, 33},
		{token.SINGLE_QUOTE, "'", 7, 35},
		{token.LINE_COMMENT, "", 8, 2},
		{token.BLOCK_COMMENT, "", 9, 2},

		// operator2
		{token.ASSIGN, "=", 11, 2},
		{token.PLUS, "+", 11, 3},
		{token.MINUS, "-", 11, 4},
		{token.NOT, "!", 11, 5},
		{token.ASTERISK, "*", 11, 6},
		{token.SLASH, "/", 11, 8},
		{token.BIT_OR, "|", 11, 9},
		{token.BIT_AND, "&", 11, 10},
		{token.BIT_NOT, "~", 11, 11},
		{token.OR, "||", 11, 12},
		{token.AND, "&&", 11, 14},
		{token.INC, "++", 11, 16},
		{token.DEC, "--", 11, 18},
		{token.PLUS_ASSIGN, "+=", 11, 20},
		{token.MINUS_ASSIGN, "-=", 11, 22},
		{token.QUESTION, "?", 11, 24},
		{token.LT, "<", 11, 25},
		{token.GT, ">", 11, 26},
		{token.LE, "<=", 11, 27},
		{token.GE, ">=", 11, 30},
		{token.EQ, "==", 11, 32},
		{token.NE, "!=", 11, 34},
		{token.UFO, "<=>", 11, 36},
		{token.ARROW, "=>", 11, 39},
		{token.PARSE, "...", 12, 2},
		{token.RANGE, "..", 12, 5},
		{token.ACCESS, ".", 12, 8},
		{token.COMMA, ",", 12, 9},
		{token.SEMICOLON, ";", 12, 10},
		{token.COLON, ":", 12, 11},
		{token.LPAREN, "(", 12, 12},
		{token.RPAREN, ")", 12, 13},
		{token.LBRACE, "{", 12, 14},
		{token.RBRACE, "}", 12, 15},
		{token.LBRACKET, "[", 12, 16},
		{token.RBRACKET, "]", 12, 17},
		{token.SINGLE_QUOTE, "'", 12, 18},

		// unknown
		{token.ERR, "$hoge", 14, 2},
		{token.ERR, "なのなの", 14, 8},

		// number literal
		{token.INTEGER, "123", 16, 2},
		{token.FLOAT, "123.456", 16, 6},
		{token.HEXA, "0x0a", 16, 14},
		{token.OCTAL, "0o777", 16, 19},
		{token.BINARY, "0b010101", 16, 25},
		{token.IMAGINARY, "123i", 16, 34},
		{token.IMAGINARY, "123.456i", 16, 39},

		// string literal
		{token.STRING, "hogehoge=+*/()がぶがぶreturn", 17, 3},
		{token.PLUS, "+", 17, 33},
		{token.IDENT, "foo", 17, 35},

		{token.STRING, "hogeがぶhoge", 18, 3},
		{token.PLUS, "+", 18, 15},
		{token.IDENT, "foo", 18, 17},
		{token.PLUS, "+", 18, 20},
		{token.STRING, "hoge", 18, 22},
		{token.PLUS, "+", 18, 27},
		{token.IDENT, "bar", 18, 28},
		{token.PLUS, "+", 18, 31},
		{token.STRING, "return", 18, 32},
		{token.PLUS, "+", 18, 38},
		{token.IDENT, "baz", 18, 40},
		{token.PLUS, "+", 18, 43},
		{token.STRING, "", 18, 44},
		{token.PLUS, "+", 18, 44},
		{token.IDENT, "qux", 18, 46},
		{token.PLUS, "+", 18, 49},
		{token.STRING, "", 18, 50},
		{token.PLUS, "+", 18, 52},
		{token.IDENT, "quux", 18, 54},

		{token.STRING, "\nhoge\ngab\n", 19, 3},
		{token.PLUS, "+", 22, 3},
		{token.IDENT, "hoge", 22, 5},

		{token.IDENT, "foo", 23, 2},
		{token.PLUS, "+", 23, 6},
		{token.STRING, "bar", 23, 9},
		{token.PLUS, "+", 23, 12},
		{token.STRING, "foo\nbar\nbaz\nがぶ", 23, 15},
		{token.PLUS, "+", 26, 6},
		{token.STRING, "baz", 26, 7},
		{token.PLUS, "+", 26, 10},
		{token.STRING, "です\nデス\ndesu", 26, 13},
		{token.PLUS, "+", 28, 6},
		{token.STRING, "", 28, 7},

		{token.EOF, "", 29, 1},
	}

	tokens := GetTokens(input)
	for i, test := range tests {
		tok := tokens[i]
		fmt.Printf("%s\n", tok.String())

		if tok.Type != test.Type {
			t.Errorf("expected token type %s, got=%s", test.Type, tok.Type)
		}
		if tok.Literal != test.Literal {
			t.Errorf("expected token literal %s, got=%s", test.Literal, tok.Literal)
		}
		if tok.Row != test.Row {
			t.Errorf("expected token row %d, got=%d", test.Row, tok.Row)
		}
		if tok.Col != test.Col {
			t.Errorf("expected token col %d, got=%d", test.Col, tok.Col)
		}
	}

}
