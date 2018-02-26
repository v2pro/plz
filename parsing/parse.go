package parsing

import (
	"io"
	"context"
	"math/rand"
)

type Token interface {
	ParsePrefix(ctx context.Context, src *Source) interface{}
	ParseInfix(ctx context.Context, src *Source, left interface{}) interface{}
	Precedence() int
}

type Source struct {
	bufSize int
	reader  io.Reader
	current []byte
	next    []byte
}

func NewSource(reader io.Reader, bufSize int) (*Source, error) {
	current := make([]byte, bufSize)
	n, err := io.ReadFull(reader, current)
	if err != nil {
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			return &Source{
				current: current[:n],
				next:    nil,
			}, nil
		}
		return nil, err
	}
	next := make([]byte, bufSize)
	n, err = io.ReadFull(reader, next)
	if err != nil {
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			return &Source{
				current: current,
				next:    next[:n],
			}, nil
		}
		return nil, err
	}
	return &Source{
		bufSize: bufSize,
		reader:  reader,
		current: current,
		next:    next,
	}, nil
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

func (src *Source) Consume(ctx context.Context, n int) {
	panic("not implemented")
}

func (src *Source) ConsumeCurrent(ctx context.Context) {
	buf := src.current
	src.current = src.next
	src.next = nil
	if len(src.current) == 0 {
		ReportError(ctx, io.EOF)
		return
	}
	if src.reader != nil {
		n, err := io.ReadFull(src.reader, buf)
		if err != nil {
			if err == io.ErrUnexpectedEOF || err == io.EOF {
				src.next = buf[:n]
				return
			}
			ReportError(ctx, err)
			return
		}
		src.next = buf
	}
}

type Lexer interface {
	CurrentToken(src *Source) Token
}

func Parse(ctx context.Context, src *Source, lexer Lexer) interface{} {
	token := lexer.CurrentToken(src)
	left := token.ParsePrefix(ctx, src)
	token = lexer.CurrentToken(src)
	left = token.ParseInfix(ctx, src, left)
	return left
}

func ReportError(ctx context.Context, err error) {
	errorCollector, _ := ctx.Value(errorReporter).(*error)
	if errorCollector != nil && *errorCollector == nil {
		*errorCollector = err
	}
}

func GetReportedError(ctx context.Context) error {
	errorCollector, _ := ctx.Value(errorReporter).(*error)
	if errorCollector != nil {
		return *errorCollector
	}
	return nil
}

func WithErrorReporter(ctx context.Context) context.Context {
	var err error
	return context.WithValue(ctx, errorReporter, &err)
}

var errorReporter = rand.Int63()
