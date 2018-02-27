package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"github.com/v2pro/plz/parse"
	"io"
	"github.com/v2pro/plz/parse/read"
)

func Test(t *testing.T) {
	t.Run("1＋1", test.Case(func(ctx *countlog.Context) {
		src := parse.NewSourceString(`1+1`)
		dst := expr.Parse(src, 0)
		must.Equal(io.EOF, src.Error())
		must.Equal(2, dst)
	}))
	t.Run("1＋1－1", test.Case(func(ctx *countlog.Context) {
		src := parse.NewSourceString(`1+1-1`)
		must.Equal(1, expr.Parse(src, 0))
	}))
	t.Run("2×3＋1", test.Case(func(ctx *countlog.Context) {
		src := parse.NewSourceString(`2*3+1`)
		must.Equal(7, expr.Parse(src, 0))
	}))
}

const precedenceAssignment = 1
const precedenceConditional = 2
const precedenceSum = 3
const precedenceProduct = 4
const precedenceExponent = 5
const precedencePrefix = 6
const precedencePostfix = 7
const precedenceCall = 8

type exprLexer struct {
	value    *valueToken
	plus     *plusToken
	minus    *minusToken
	multiply *multiplyToken
}

var expr = newExprLexer()

func newExprLexer() *exprLexer {
	return &exprLexer{
		value:    &valueToken{},
		plus:     &plusToken{},
		minus:    &minusToken{},
		multiply: &multiplyToken{},
	}
}

func (lexer *exprLexer) Parse(src *parse.Source, precedence int) interface{} {
	return parse.Parse(src, lexer, precedence)
}

func (lexer *exprLexer) InfixToken(src *parse.Source) (parse.InfixToken, int) {
	switch src.Peek()[0] {
	case '+':
		return lexer.plus, precedenceSum
	case '-':
		return lexer.minus, precedenceSum
	case '*':
		return lexer.multiply, precedenceProduct
	default:
		return nil, 0
	}
}

func (lexer *exprLexer) PrefixToken(src *parse.Source) parse.PrefixToken {
	return lexer.value
}

type valueToken struct {
}

func (token *valueToken) PrefixParse(src *parse.Source) interface{} {
	return read.Int(src)
}

type plusToken struct {
}

func (token *plusToken) InfixParse(src *parse.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.ConsumeN(1)
	rightValue := expr.Parse(src, precedenceSum).(int)
	return leftValue + rightValue
}

type minusToken struct {
}

func (token *minusToken) InfixParse(src *parse.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.ConsumeN(1)
	rightValue := expr.Parse(src, precedenceSum).(int)
	return leftValue - rightValue
}

type multiplyToken struct {
}

func (token *multiplyToken) InfixParse(src *parse.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.ConsumeN(1)
	rightValue := expr.Parse(src, precedenceProduct).(int)
	return leftValue * rightValue
}