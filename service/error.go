package service

type ErrorNumber interface {
	ErrorNumber() int
}

type WithNumberError struct {
	Number  int
	Message string
}

func (err *WithNumberError) ErrorNumber() int {
	return err.Number
}

func (err *WithNumberError) Error() string {
	return err.Message
}
