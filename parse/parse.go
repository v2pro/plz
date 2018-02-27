package parse

import (
	"io"
	"github.com/v2pro/plz/countlog"
)

func Parse(src *Source, lexer Lexer) interface{} {
	token := lexer.TokenOf(src)
	left := token.ParsePrefix(src)
	countlog.Trace("prefix", "token", token)
	if src.Error() != nil {
		return left
	}
	token = lexer.TokenOf(src)
	countlog.Trace("infix", "token", token)
	newLeft := token.ParseInfix(src, left)
	if newLeft != nil {
		left = newLeft
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

type Token interface {
	ParsePrefix(src *Source) interface{}
	ParseInfix(src *Source, left interface{}) interface{}
	Precedence() int
}

type Lexer interface {
	TokenOf(src *Source) Token
}
