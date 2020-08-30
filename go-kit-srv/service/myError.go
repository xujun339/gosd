package service

type MyError struct {
	Code int
	Msg string
}

func NewMyError(code int, msg string) *MyError {
	return &MyError{Code: code, Msg: msg}
}

func (e *MyError)Error() string {
	return e.Msg
}


