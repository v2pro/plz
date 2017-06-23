package any

type Any interface {
	LastError() error
	MustBeValid() Any
	ToBool() bool
	ToInt() int
	ToInt32() int32
	ToInt64() int64
	ToUint() uint
	ToUint32() uint32
	ToUint64() uint64
	ToFloat32() float32
	ToFloat64() float64
	ToString() string
	ToVal(val interface{})
	Get(path ...interface{}) Any
	// TODO: add Set
	Size() int
	Keys() []string
	GetInterface() interface{}
}
