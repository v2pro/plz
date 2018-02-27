package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/parse"
	"io"
)

func Test(t *testing.T) {
	t.Run("one plus one", test.Case(func(ctx *countlog.Context) {
		src := parse.NewSourceString(`1+1`)
		dst := expr.Parse(src)
		must.Equal(io.EOF, src.Error())
		must.Equal(2, dst)
	}))
	t.Run("one plus one minus one", test.Case(func(ctx *countlog.Context) {
		src := parse.NewSourceString(`1+1-1`)
		must.Equal(1, expr.Parse(src))
	}))
}

type exprLexer struct {
	value *valueToken
	plus  *plusToken
	minus *minusToken
}

var expr = newExprLexer()

func newExprLexer() *exprLexer {
	return &exprLexer{
		value: &valueToken{},
		plus:  &plusToken{},
		minus: &minusToken{},
	}
}

func (lexer *exprLexer) Parse(src *parse.Source) interface{} {
	return parse.Parse(src, lexer)
}

func (lexer *exprLexer) TokenOf(src *parse.Source) parse.Token {
	switch src.Current()[0] {
	case '+':
		return lexer.plus
	case '-':
		return lexer.minus
	default:
		return lexer.value
	}
}

type valueToken struct {
}

func (token *valueToken) ParsePrefix(src *parse.Source) interface{} {
	return parse.Int(src)
}

func (token *valueToken) ParseInfix(src *parse.Source, left interface{}) interface{} {
	return nil
}

func (token *valueToken) Precedence() int {
	return 0
}

func (token *valueToken) String() string {
	return "val"
}

type plusToken struct {
}

func (token *plusToken) ParsePrefix(src *parse.Source) interface{} {
	return nil
}

func (token *plusToken) ParseInfix(src *parse.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.ConsumeN(1)
	rightValue := expr.Parse(src).(int)
	return leftValue + rightValue
}

func (token *plusToken) Precedence() int {
	return 0
}

func (token *plusToken) String() string {
	return "+"
}

type minusToken struct {
}

func (token *minusToken) ParsePrefix(src *parse.Source) interface{} {
	return nil
}

func (token *minusToken) ParseInfix(src *parse.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.ConsumeN(1)
	rightValue := expr.Parse(src).(int)
	return leftValue - rightValue
}

func (token *minusToken) Precedence() int {
	return 0
}

func (token *minusToken) String() string {
	return "-"
}
