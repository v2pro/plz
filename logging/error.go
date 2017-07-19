package logging

type errorWrapper struct {
	cause error
	msg string
	kv []interface{}
}

func (err *errorWrapper) Cause() error {
	return err.cause
}

func (err *errorWrapper) Error() string {
	return err.msg
}

func (err *errorWrapper) KV() []interface{} {
	return err.kv
}