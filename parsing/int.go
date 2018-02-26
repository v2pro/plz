package parsing

import (
	"strconv"
	"context"
	"errors"
)

var intDigits []int8

const invalidCharForNumber = int8(-1)
const uint64SafeToMultiple10 = uint64(0xffffffffffffffff)/10 - 1

func init() {
	intDigits = make([]int8, 256)
	for i := 0; i < len(intDigits); i++ {
		intDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		intDigits[i] = i - int8('0')
	}
}

func (src *Source) ConsumeInt(ctx context.Context) int {
	if strconv.IntSize == 32 {
		return int(src.ConsumeInt32(ctx))
	}
	return int(src.ConsumeInt64(ctx))
}

func (src *Source) ConsumeInt32(ctx context.Context) int32 {
	panic("not implemented")
}

func (src *Source) ConsumeInt64(ctx context.Context) int64 {
	return 0
}

func (src *Source) ConsumeUint64(ctx context.Context) uint64 {
	value := uint64(0)
	for GetReportedError(ctx) == nil {
		buf := src.Current()
		for _, c := range buf {
			ind := intDigits[c]
			if ind == invalidCharForNumber {
				return value
			}
			if value > uint64SafeToMultiple10 {
				value2 := (value << 3) + (value << 1) + uint64(ind)
				if value2 < value {
					ReportError(ctx, errors.New("ConsumeUint64: overflow"))
					return 0
				}
				value = value2
				continue
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
		src.ConsumeCurrent(ctx)
	}
	return value
}
