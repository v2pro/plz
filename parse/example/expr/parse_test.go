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

func (lexer *exprLexer) TokenOf(src *parse.Source) parse.Token {
	switch src.Current()[0] {
	case '+':
		return lexer.plus
	case '-':
		return lexer.minus
	case '*':
		return lexer.multiply
	default:
		return lexer.value
	}
}

type valueToken struct {
	parse.DummyToken
}

func (token *valueToken) PrefixParse(src *parse.Source) interface{} {
	return parse.Int(src)
}

func (token *valueToken) PrefixPrecedence() int {
	return 1
}

func (token *valueToken) String() string {
	return "val"
}

type plusToken struct {
	parse.DummyToken
}

func (token *plusToken) InfixParse(src *parse.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.ConsumeN(1)
	rightValue := expr.Parse(src, token.InfixPrecedence()).(int)
	return leftValue + rightValue
}

func (token *plusToken) InfixPrecedence() int {
	return precedenceSum
}

func (token *plusToken) String() string {
	return "+"
}

type minusToken struct {
	parse.DummyToken
}

func (token *minusToken) InfixParse(src *parse.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.ConsumeN(1)
	rightValue := expr.Parse(src, token.InfixPrecedence()).(int)
	return leftValue - rightValue
}

func (token *minusToken) InfixPrecedence() int {
	return precedenceSum
}

func (token *minusToken) String() string {
	return "-"
}

type multiplyToken struct {
	parse.DummyToken
}

func (token *multiplyToken) InfixParse(src *parse.Source, left interface{}) interface{} {
	leftValue := left.(int)
	src.ConsumeN(1)
	rightValue := expr.Parse(src, token.InfixPrecedence()).(int)
	return leftValue * rightValue
}

func (token *multiplyToken) InfixPrecedence() int {
	return precedenceProduct
}

func (token *multiplyToken) String() string {
	return "*"
}
