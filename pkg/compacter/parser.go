package compacter

import (
	"bufio"
	"io"

	"github.com/fairyhunter13/go-lexer"
	"github.com/fairyhunter13/pool"
)

const (
	AllToken lexer.TokenType = iota
)

type Parser struct {
	Scanner *bufio.Scanner

	// state of the parser
	AfterEqual bool
	Quote      rune
}

func NewParser(rd io.Reader) *Parser {
	return &Parser{
		Scanner: bufio.NewScanner(rd),
	}
}

func (p *Parser) Scan() bool { return p.Scanner.Scan() }
func (p *Parser) Err() error { return p.Scanner.Err() }
func (p *Parser) GetLine(l *lexer.L) (fn lexer.StateFunc) {
	char := l.Peek()
	for char != lexer.EOFRune {
		if !p.AfterEqual && char == '=' {
			p.AfterEqual = true
		}

		if p.AfterEqual && (char == '\'' || char == '"') {
			fn = p.AllQuote
			break
		}

		char = l.NextPeek()
	}

	l.Emit(AllToken)
	return
}

func (p *Parser) AllQuote(l *lexer.L) (fn lexer.StateFunc) {
	p.Quote = l.Next()
	char := l.Peek()
	for {
		if char == lexer.EOFRune {
			if p.Scanner.Scan() {
				l.Append(p.Scanner.Text())
				char = l.Peek()
				continue
			}
			break
		}

		if char == '\\' {
			l.Next()
			char = l.NextPeek()
			continue
		}

		if char == p.Quote {
			l.Next()
			break
		}

		char = l.NextPeek()
	}

	l.Emit(AllToken)
	return
}

func (p *Parser) Text() (line string) {
	line = p.Scanner.Text()
	lex := lexer.New(line, p.GetLine)
	lex.Start()

	builder := pool.GetStrBuilder()
	defer pool.Put(builder)
	var (
		token *lexer.Token
		done  bool
	)
	for {
		token, done = lex.NextToken()
		if done {
			break
		}

		if token.Type != AllToken {
			continue
		}

		builder.WriteString(token.Value)
	}

	line = builder.String()
	return
}
