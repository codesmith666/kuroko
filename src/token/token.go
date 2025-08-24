package token

import (
	"fmt"
	"sort"
)

type TokenType string

const (
	// controller
	ERR TokenType = "ERR"
	EOF TokenType = "EOF"

	// Operators
	ASSIGN        TokenType = "="
	PLUS          TokenType = "+"
	MINUS         TokenType = "-"
	NOT           TokenType = "!"
	ASTERISK      TokenType = "*"
	SLASH         TokenType = "/"
	BIT_OR        TokenType = "|"
	BIT_AND       TokenType = "&"
	BIT_NOT       TokenType = "~"
	OR            TokenType = "||"
	AND           TokenType = "&&"
	INC           TokenType = "++"
	DEC           TokenType = "--"
	PLUS_ASSIGN   TokenType = "+="
	MINUS_ASSIGN  TokenType = "-="
	QUESTION      TokenType = "?"
	LT            TokenType = "<"
	GT            TokenType = ">"
	LE            TokenType = "<="
	GE            TokenType = ">="
	EQ            TokenType = "=="
	NE            TokenType = "!="
	UFO           TokenType = "<=>"
	ARROW         TokenType = "=>"
	PARSE         TokenType = "..."
	RANGE         TokenType = ".."
	ACCESS        TokenType = "."
	LINE_COMMENT  TokenType = "//"
	BLOCK_COMMENT TokenType = "/*"
	COMMENT_CLOSE TokenType = "*/"
	COMMA         TokenType = ","
	SEMICOLON     TokenType = ";"
	COLON         TokenType = ":"
	LPAREN        TokenType = "("
	RPAREN        TokenType = ")"
	LBRACE        TokenType = "{"
	RBRACE        TokenType = "}"
	LBRACKET      TokenType = "["
	RBRACKET      TokenType = "]"
	DOUBLE_QUOTE  TokenType = "\""
	SINGLE_QUOTE  TokenType = "'"
	INSTANCEOF    TokenType = "instanceof"

	INLINE_OPEN TokenType = "${"
	YEN_R       TokenType = "\\r"
	YEN_N       TokenType = "\\n"
	YEN_T       TokenType = "\\t"
	YEN_DQ      TokenType = "\\\""

	// Identifiers + literals + Type
	IDENT     TokenType = "IDENT"     // add, foobar, x, y, ...
	TYPE      TokenType = "TYPE"      // string,number,any
	INTEGER   TokenType = "INTEGER"   // 1343456
	FLOAT     TokenType = "FLOAT"     // 0.123
	IMAGINARY TokenType = "IMAGINARY" // 4i
	HEXA      TokenType = "HEXA"
	BINARY    TokenType = "BINARY"
	OCTAL     TokenType = "OCTAL"
	STRING    TokenType = "STRING"

	// 予約語
	TRUE   TokenType = "TRUE"
	FALSE  TokenType = "FALSE"
	IF     TokenType = "IF"
	ELSE   TokenType = "ELSE"
	LOOP   TokenType = "LOOP"
	RETURN TokenType = "RETURN"
	SHARE  TokenType = "SHARE"
	CONST  TokenType = "CONST"
	IMM    TokenType = "IMM"
	MUT    TokenType = "MUT"
)

// オペレータの配列
var NormalOperators = []TokenType{
	ASSIGN,
	PLUS,
	MINUS,
	NOT,
	ASTERISK,
	SLASH,
	BIT_OR,
	BIT_AND,
	BIT_NOT,
	OR,
	AND,
	INC,
	DEC,
	PLUS_ASSIGN,
	MINUS_ASSIGN,
	QUESTION,
	LT,
	GT,
	LE,
	GE,
	EQ,
	NE,
	UFO,
	ARROW,
	PARSE,
	RANGE,
	ACCESS,
	LINE_COMMENT,
	BLOCK_COMMENT,
	COMMENT_CLOSE,
	COMMA,
	SEMICOLON,
	COLON,
	LPAREN,
	RPAREN,
	LBRACE,
	RBRACE,
	LBRACKET,
	RBRACKET,
	DOUBLE_QUOTE,
	SINGLE_QUOTE,
	INSTANCEOF,
}

var StringOperators = []TokenType{
	INLINE_OPEN,
	DOUBLE_QUOTE,
	YEN_R,
	YEN_N,
	YEN_T,
	YEN_DQ,
}

var CommentOperators = []TokenType{
	COMMENT_CLOSE,
}

type Token struct {
	Type    TokenType
	Literal string
	Row     int
	Col     int
}

func (t *Token) String() string {
	return fmt.Sprintf("%s[row=%d:col=%d] %s", t.Literal, t.Row, t.Col, t.Type)
}

var Reserved = map[string]TokenType{
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"loop":   LOOP,
	"return": RETURN,
	"share":  SHARE,
	"const":  CONST,
	"imm":    IMM,
	"mut":    MUT,
}

var Types = map[string]bool{
	"string":  true,
	"array":   true,
	"object":  true,
	"number":  true,
	"boolean": true,
	"void":    true,
	"any":     true,
}

var initialized = false

func Initialize() {
	if initialized {
		return
	}
	initialized = true
	// オペレータを長い順にソート
	sort.Slice(NormalOperators, func(i, j int) bool {
		return len(NormalOperators[i]) > len(NormalOperators[j])
	})
	sort.Slice(StringOperators, func(i, j int) bool {
		return len(StringOperators[i]) > len(StringOperators[j])
	})
	sort.Slice(CommentOperators, func(i, j int) bool {
		return len(CommentOperators[i]) > len(CommentOperators[j])
	})

}
