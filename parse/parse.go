package parse

import (
	"io"
	"github.com/v2pro/plz/countlog"
	"fmt"
	"errors"
)

func Parse(src *Source, lexer Lexer, precedence int) interface{} {
	token := lexer.TokenOf(src)
	if token == nil || token.PrefixPrecedence() == NotPrefix {
		src.ReportError(errors.New("can not parse"))
		return nil
	}
	left := token.PrefixParse(src)
	countlog.Trace("prefix", "token", token)
	for {
		if src.Error() != nil {
			return left
		}
		token = lexer.TokenOf(src)
		tokenPrecedence := token.InfixPrecedence()
		if token == nil || tokenPrecedence == NotInfix {
			return left
		}
		if precedence >= tokenPrecedence {
			return left
		}
		countlog.Trace("infix", "token", token)
		left = token.InfixParse(src, left)
	}
	return left
}

type Source struct {
	err     error
	reader  io.Reader
	current []byte
	buf     []byte
}

func NewSource(reader io.Reader, buf []byte) (*Source, error) {
	n, err := reader.Read(buf)
	if n == 0 {
		return nil, err
	}
	return &Source{
		reader:  reader,
		current: buf[:n],
		buf:     buf,
	}, nil
}

func NewSourceString(src string) *Source {
	return &Source{
		current: []byte(src),
	}
}

func (src *Source) SetBuffer(buf []byte) {
	src.buf = buf
}

func (src *Source) Current() []byte {
	return src.current
}

func (src *Source) ConsumeN(n int) {
	src.current = src.current[n:]
	if len(src.current) == 0 {
		src.Consume()
	}
}

func (src *Source) Consume() {
	if src.reader == nil {
		src.current = nil
		src.ReportError(io.EOF)
		return
	}
	n, err := src.reader.Read(src.buf)
	if err != nil {
		src.ReportError(err)
	}
	src.current = src.buf[:n]
}

func (src *Source) ReportError(err error) {
	if src.err == nil {
		src.err = err
	}
}

func (src *Source) Error() error {
	return src.err
}

const NotPrefix = 0
const NotInfix = 0

type Token interface {
	fmt.Stringer
	PrefixParse(src *Source) interface{}
	InfixParse(src *Source, left interface{}) interface{}
	PrefixPrecedence() int
	InfixPrecedence() int
}

type DummyToken struct {
}

func (token DummyToken) PrefixParse(src *Source) interface{} {
	panic("not prefix")
}

func (token DummyToken) PrefixPrecedence() int {
	return NotPrefix
}

func (token DummyToken) InfixParse(src *Source, left interface{}) interface{} {
	panic("not infix")
}

func (token DummyToken) InfixPrecedence() int {
	return NotInfix
}

type Lexer interface {
	TokenOf(src *Source) Token
}
