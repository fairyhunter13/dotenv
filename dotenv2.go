package dotenv

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fairyhunter13/go-lexer"
	"github.com/fairyhunter13/pool"
)

// Option is the option struct.
type Option struct {
	Overload bool
	Paths    []string
}

// Default specifies the default values of the option struct.
func (o *Option) Default() {
	if len(o.Paths) <= 0 {
		o.Paths = append(o.Paths, ".env")
	}
}

// FnOption is a functional option for this package
type FnOption func(*Option)

// WithPaths specifies all paths of the environment file.
func WithPaths(paths ...string) FnOption {
	return func(o *Option) {
		o.Paths = paths
	}
}

// WithOverload overdrives the current set environment.
func WithOverload(overload bool) FnOption {
	return func(o *Option) {
		o.Overload = overload
	}
}

// Load2 is the version 2 of Load.
func Load2(opts ...FnOption) (err error) {
	optStruct := new(Option)
	for _, opt := range opts {
		opt(optStruct)
	}
	optStruct.Default()

	var isSuccess bool
	for _, path := range optStruct.Paths {
		err = loadFile2(path, optStruct.Overload)
		if err != nil {
			continue
		}

		isSuccess = true
	}

	if isSuccess {
		err = nil
	}
	return
}

// loadFile2 is the version 2 of loadFile.
func loadFile2(path string, overload bool) error {
	env, err := ReadFile2(path)
	if err != nil {
		return err
	}
	LoadMap(env, overload)
	return nil
}

// ReadFile2 is the version 2 of ReadFile.
func ReadFile2(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read2(f)
}

// LoadReader2 is the version 2 of LoadReader.
func LoadReader2(r io.Reader) error {
	env, err := Read2(r)
	if err != nil {
		return err
	}
	LoadMap(env, false)
	return nil
}

// Read2 is the version 2 of Read.
func Read2(rd io.Reader) (envMap map[string]string, err error) {
	return newParser(rd).Scan()
}

type parser struct {
	Scanner *bufio.Scanner

	// state of the parser
	AfterEqual bool
	Quote      rune
}

func newParser(rd io.Reader) *parser {
	return &parser{
		Scanner: bufio.NewScanner(rd),
	}
}

func (p *parser) Scan() (resultMap map[string]string, err error) {
	var (
		k, v, line string
		envMap     = map[string]string{}
	)
	for p.Scanner.Scan() {
		line = p.getMultiline(envMap)
		k, v, err = ParseString(line)
		if err == ErrInvalidln {
			err = fmt.Errorf("could not parse file: %v", err)
			return
		}

		if err != nil {
			continue
		}

		envMap[k] = v
	}

	if err = p.Scanner.Err(); err != nil {
		err = fmt.Errorf("error reading file: %v", err)
		return
	}

	resultMap = envMap
	return
}

const (
	AllToken lexer.TokenType = iota
)

func (p *parser) GetLine(l *lexer.L) (fn lexer.StateFunc) {
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

func (p *parser) AllQuote(l *lexer.L) (fn lexer.StateFunc) {
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

func (p *parser) getMultiline(envMap map[string]string) (line string) {
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
	if regexVar.MatchString(line) {
		line = regexVar.ReplaceAllStringFunc(line, func(s string) string {
			return envMap[strings.Trim(s, "${}")]
		})
	}
	return
}
