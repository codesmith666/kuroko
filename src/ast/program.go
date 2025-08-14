package ast

import (
	"bytes"
)

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String(depth int) string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String(depth + 1))
	}
	// re := regexp.MustCompile(`\n{2,}`)
	// return re.ReplaceAllString(out.String(), "\n")
	return out.String()
}
