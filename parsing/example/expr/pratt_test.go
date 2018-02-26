package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"context"
	"github.com/v2pro/plz/parsing"
)

func Test(t *testing.T) {
	t.Run("one plus one", test.Case(func(ctx *countlog.Context) {
		src := parsing.NewSourceString(`1+1`)
		dst := parsing.Parse(ctx, src, &exprLexer{})
		must.Equal(2, dst)
	}))
}

type exprLexer struct {
}

var exprLexerInstance = parsing.Lexer(&exprLexer{})

func (lexer *exprLexer) CurrentToken(src *parsing.Source) parsing.Token {
	if src.Current()[0] == '+' {
		return plus
	}
	return value
}

type valueToken struct {
}

var value = parsing.Token(&valueToken{})

func (token *valueToken) ParsePrefix(ctx context.Context, src *parsing.Source) interface{} {
	return src.ConsumeInt()
}

func (token *valueToken) ParseInfix(ctx context.Context, src *parsing.Source, left interface{}) interface{} {
	return nil
}

func (token *valueToken) Precedence() int {
	return 0
}

type plusToken struct {
}

var plus = parsing.Token(&plusToken{})


func (token *plusToken) ParsePrefix(ctx context.Context, src *parsing.Source) interface{} {
	return nil
}

func (token *plusToken) ParseInfix(ctx context.Context, src *parsing.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.Consume(1)
	rightValue := parsing.Parse(ctx, src, exprLexerInstance).(int)
	return leftValue + rightValue
}

func (token *plusToken) Precedence() int {
	return 0
}