package pratt

import (
	"io"
	"context"
)

type Token interface {
	ParsePrefix(ctx context.Context, src *Source) interface{}
	ParseInfix(ctx context.Context, src *Source, left interface{}) interface{}
	Precedence() int
}

type Source struct {
	reader  io.Reader
	current []byte
	next    []byte
}

func NewSource(reader io.Reader) (*Source, error) {
	return nil, nil
}

func NewSourceString(src string) *Source {
	return &Source{
		reader:  nil,
		current: []byte(src),
		next:    nil,
	}
}

func (src *Source) Current() []byte {
	return src.current
}

func (src *Source) Next() []byte {
	return src.next
}

func (src *Source) Consume(n int) error {
	return nil
}

type Lexer interface {
	CurrentToken(src *Source) Token
}

func Parse(ctx context.Context, src *Source, lexer Lexer) interface{} {
	return nil
}
