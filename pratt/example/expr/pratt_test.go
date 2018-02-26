package test

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/pratt"
	"github.com/v2pro/plz/test/must"
)

func Test(t *testing.T) {
	t.Run("one plus one", test.Case(func(ctx *countlog.Context) {
		src := pratt.NewSourceString(`1+1`)
		dst := pratt.Parse(ctx, src, &exprLexer{})
		must.Equal(2, dst)
	}))
}

type exprLexer struct {
}

func (lexer *exprLexer) CurrentToken(src *pratt.Source) pratt.Token {
	return nil
}

type valueToken struct {
}