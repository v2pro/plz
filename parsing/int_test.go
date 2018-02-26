package parsing

import (
	"testing"
	"github.com/v2pro/plz/test"
	"github.com/v2pro/plz/countlog"
	"github.com/v2pro/plz/test/must"
	"io"
	"strings"
)

func Test_ConsumeUint64_from_string(t *testing.T) {
	type TestCase struct {
		Input  string
		Output uint64
	}
	testCases := []TestCase{{
		"1", 1,
	}, {
		"12", 12,
	}, {
		"123", 123,
	}}
	for _, tmp := range testCases {
		testCase := tmp
		t.Run(testCase.Input, test.Case(func(ctxObj *countlog.Context) {
			ctx := WithErrorReporter(ctxObj)
			src := NewSourceString(testCase.Input)
			must.Equal(testCase.Output, src.ConsumeUint64(ctx))
		}))
	}
	t.Run("Overflow", test.Case(func(ctxObj *countlog.Context) {
		ctx := WithErrorReporter(ctxObj)
		src := NewSourceString("18446744073709551615")
		must.Equal(uint64(18446744073709551615), src.ConsumeUint64(ctx))
		must.Equal(io.EOF, GetReportedError(ctx))
		ctx = WithErrorReporter(ctxObj)
		src = NewSourceString("18446744073709551616")
		must.Equal(uint64(0), src.ConsumeUint64(ctx))
		must.Pass(io.EOF != GetReportedError(ctx))
	}))
}

func Test_ConsumeUint64_from_reader(t *testing.T) {
	type TestCase struct {
		Input    string
		Output   uint64
		Selected bool
	}
	testCases := []TestCase{{
		Input: "1", Output: 1,
	}, {
		Input: "12", Output: 12,
	}, {
		Input: "123", Output: 123,
	}, {
		Input: "1234", Output: 1234, Selected: true,
	}, {
		Input: "12345", Output: 12345,
	}}
	for _, testCase := range testCases {
		if testCase.Selected {
			testCases = []TestCase{testCase}
			break
		}
	}
	for _, tmp := range testCases {
		testCase := tmp
		t.Run(testCase.Input, test.Case(func(ctxObj *countlog.Context) {
			ctx := WithErrorReporter(ctxObj)
			src := must.Call(NewSource, strings.NewReader(testCase.Input), 2)[0].(*Source)
			must.Equal(testCase.Output, src.ConsumeUint64(ctx))
		}))
	}
}
