package lexer

import (
	"fmt"
	"monkey/token"
	"regexp"

	"github.com/mattn/go-runewidth"
)

type LexerMode string

const (
	ROOT_MODE   LexerMode = "root"
	NORMAL_MODE LexerMode = "normal"
	STRING_MODE LexerMode = "string"
)

// 新しい字句解析器を返す
type Lexer struct {
	input       string
	position    int // current position in input (points to current char)
	last        int
	reHexa      *regexp.Regexp
	reOctal     *regexp.Regexp
	reBinary    *regexp.Regexp
	reImaginary *regexp.Regexp
	reFloat     *regexp.Regexp
	reInteger   *regexp.Regexp
	tokens      []*token.Token
	row         int
	col         int
}

// 新しい Lexerを返す
func GetTokens(input string) []*token.Token {
	token.Initialize()
	l := &Lexer{
		input:       input,
		position:    0,
		last:        len(input),
		reHexa:      regexp.MustCompile("^0x[0-9a-fA-F]+"),
		reOctal:     regexp.MustCompile("^0o[0-7]+"),
		reBinary:    regexp.MustCompile("^0b[01]+"),
		reImaginary: regexp.MustCompile(`^\d+(\.\d+)?i`),
		reFloat:     regexp.MustCompile(`^\d+\.\d+`),
		reInteger:   regexp.MustCompile(`^\d+`),
		tokens:      []*token.Token{},
		row:         1, // 行：1オリジン
		col:         1, // 列：1オリジン
	}
	l.tokenizeNormal(ROOT_MODE)
	l.addToken(token.EOF, "", l.row, l.col)
	return l.tokens
}

// これくらいないのかよ関数その１
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ホワイトスペースを読み飛ばす
func (l *Lexer) skipWhitespace() {
	for l.position < l.last {
		ch := l.input[l.position]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.position++
			l.col++
			continue
		} else if ch == '\n' {
			l.position++
			l.row++
			l.col = 1 // 1オリジン
		}
		break
	}
}

// l.positionからオペレータを探す
func (l *Lexer) getOperator(operators *[]token.TokenType) (string, int, int) {
	for _, v := range *operators {
		length := len(v)
		sliced := l.input[l.position:MinInt(l.position+length, l.last)]
		if sliced == string(v) {
			col := l.col
			row := l.row
			l.position += length
			l.col += length
			return sliced, row, col
		}
	}
	return "", 0, 0
}

// 識別子を探す
func (l *Lexer) getIdentifier() (string, int, int) {
	var i = l.position

	// １文字目をチェック（数字を許さない）
	if i < l.last {
		ch := l.input[i]
		if !('a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_') {
			return "", 0, 0
		}
	}
	// A-Za-z0-9_
	for i < l.last {
		ch := l.input[i]
		if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || '0' <= ch && ch <= '9' {
			i++
		} else {
			break
		}
	}
	if i != l.position {
		col := l.col
		row := l.row
		t := l.input[l.position:i]
		l.col += i - l.position // ここでtokenの最後まで進めちゃってるのが問題
		l.position = i
		return t, row, col
	}
	return "", 0, 0

}

// 数値を探す
func (l *Lexer) getNumber() (string, token.TokenType, int, int) {
	var i = l.position

	row := l.row
	col := l.col
	// 順序が重要な数値リテラル検出処理
	if m := l.reHexa.FindString(l.input[i:]); m != "" {
		l.position += len(m)
		l.col += len(m)
		return m, token.HEXA, row, col
	} else if m := l.reOctal.FindString(l.input[i:]); m != "" {
		l.position += len(m)
		l.col += len(m)
		return m, token.OCTAL, row, col
	} else if m := l.reBinary.FindString(l.input[i:]); m != "" {
		l.position += len(m)
		l.col += len(m)
		return m, token.BINARY, row, col
	} else if m := l.reImaginary.FindString(l.input[i:]); m != "" {
		l.position += len(m)
		l.col += len(m)
		return m, token.IMAGINARY, row, col
	} else if m := l.reFloat.FindString(l.input[i:]); m != "" {
		l.position += len(m)
		l.col += len(m)
		return m, token.FLOAT, row, col
	} else if m := l.reInteger.FindString(l.input[i:]); m != "" {
		l.position += len(m)
		l.col += len(m)
		return m, token.INTEGER, row, col
	}
	return "", "", 0, 0
}

// 未知のトークンを見つけた
func (l *Lexer) getUnknownToken() (string, int, int) {
	var i = l.position
	for i < l.last {
		ch := l.input[i]
		if ch != '\r' && ch != '\n' && ch != '\t' && ch != ' ' {
			i++
		} else {
			break
		}
	}
	if i != l.position {
		row := l.row
		col := l.col
		t := l.input[l.position:i]
		l.col += i - l.position
		l.position = i
		return t, row, col
	}
	return "", 0, 0
}

// UTF8で1文字取得する
func (l *Lexer) getRune() string {
	// １文字取り込んで進める（UTF8に対応する）
	var str string
	s := l.input[l.position]

	if s&0xf0 == 0xf0 {
		// 4バイトの先頭
		str += l.input[l.position : l.position+4]
		l.position += 4
		l.col += runewidth.StringWidth(str)
	} else if s&0xe0 == 0xe0 {
		// 3バイトの先頭
		str += l.input[l.position : l.position+3]
		l.position += 3
		l.col += runewidth.StringWidth(str)
	} else if s&0xc0 == 0xc0 {
		// 110xxxxx（2バイト文字）
		str += l.input[l.position : l.position+2]
		l.position += 2
		l.col += runewidth.StringWidth(str)
	} else if s&0x80 == 0 {
		// 0xxxxxxx(1バイト文字)
		str += l.input[l.position : l.position+1]
		l.position++
		if s == '\n' {
			l.row++
			l.col = 1
		} else {
			l.col++
		}
	} else {
		// UTF8の2～4バイト目が突然現れたら無視して進める
		l.position++
		l.col++
	}
	return str
}

// トークンを配列に追加する
func (l *Lexer) addToken(tokenType token.TokenType, tokenLiteral string, row int, col int) {
	l.tokens = append(l.tokens, &token.Token{Type: tokenType, Literal: tokenLiteral, Row: row, Col: col})
}

// 現在地から改行（か終端）が見つかるまでをコメントとして取得する
// getRuneしてもいいけどUTF-8では\nが不意に出てくることはない
func (l *Lexer) getLineComment() string {
	var i = l.position
	for i < l.last {
		ch := l.input[i]
		if ch == '\r' || ch == '\n' {
			break
		}
		i++
	}
	t := l.input[l.position:i]
	l.position = i
	return t

}

func (l *Lexer) getBlockComment() string {
	var str = ""
	for l.position < l.last {
		// オペレータを探す
		if ope, _, _ := l.getOperator(&token.CommentOperators); ope != "" {
			if ope == "*/" {
				return str
			}
		} else {
			str += l.getRune()
		}
	}
	return str
}

// 言語の構文で文字列をトークン化する
func (l *Lexer) tokenizeNormal(parentMode LexerMode) {

	for l.position < l.last {
		// ホワイトスペースをスキップ
		l.skipWhitespace()

		// オペレータを探す
		if ope, row, col := l.getOperator(&token.NormalOperators); ope != "" {
			switch ope {
			case "\"":
				l.tokenizeString(NORMAL_MODE)
			case "}":
				if parentMode == STRING_MODE {
					// 文字列内の構文解析か？
					return
				} else {
					l.addToken(token.TokenType(ope), ope, row, col)
				}
			// リテラルとしてコメントを取る
			// コメントの中の構文は終端のみなのでtokenizeではなくgetで実装
			case "//":
				c := l.getLineComment()
				length := len(l.tokens)
				if length > 0 && l.tokens[length-1].Type == token.LINE_COMMENT {
					l.tokens[length-1].Literal += ("\n" + c)
				} else {
					l.addToken(token.LINE_COMMENT, c, row, col)
				}
				fmt.Printf("%s\n", l.tokens[length-1].Literal)
			case "/*":
				l.addToken(token.BLOCK_COMMENT, l.getBlockComment(), row, col)
			default:
				l.addToken(token.TokenType(ope), ope, row, col)
			}
			continue
		}

		// 予約語と型と識別子を探す
		if ide, row, col := l.getIdentifier(); ide != "" {
			if t, ok := token.Reserved[ide]; ok {
				l.addToken(t, ide, row, col)
			} else if u := token.Types[ide]; u {
				l.addToken(token.TYPE, ide, row, col)
			} else {
				l.addToken(token.IDENT, ide, row, col)
			}
			continue
		}

		// 数値を探す
		if n, t, r, c := l.getNumber(); n != "" {
			l.addToken(t, n, r, c)
			continue
		}

		// 未知のトークンを検出する
		if unknown, row, col := l.getUnknownToken(); unknown != "" {
			l.addToken(token.ERR, unknown, row, col)
			continue
		}
	}
}

// 文字列の構文で文字列をトークン化する
func (l *Lexer) tokenizeString(parentMode LexerMode) {
	var str = ""

	var r int = l.row
	var c int = l.col
	for l.position < l.last {

		// オペレータを探す
		if ope, row, col := l.getOperator(&token.StringOperators); ope != "" {
			switch ope {
			case "\\r":
				str += "\r"
			case "\\n":
				str += "\n"
			case "\\t":
				str += "\t"
			case "\\\"":
				str += "\""
			case "${":
				l.addToken(token.STRING, str, r, c) // STRING
				l.addToken(token.PLUS, "+", row, col)
				l.tokenizeNormal(STRING_MODE)
				l.addToken(token.PLUS, "+", l.row, l.col-1) // }
				c = l.col
				r = l.row
				str = ""
				continue
			// 文字列の終了を検知したのでstrをtokenにする
			case "\"":
				fmt.Printf("%d::%d - %d::%d\n", r, c, l.row, l.col-1)
				l.addToken(token.STRING, str, r, c)
				return
			}
			continue
		} else {
			str += l.getRune()

		}
	}
}
