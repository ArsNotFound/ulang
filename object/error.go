package object

type Error struct {
	Message string
}

func (e *Error) Bool() bool {
	return false
}

func (e *Error) Clone() Object {
	return &Error{Message: e.Message}
}

func (e *Error) Type() Type {
	return ErrorType
}

func (e *Error) String() string {
	return e.Message
}

func (e *Error) Inspect() string {
	return "ERROR:" + e.Message
}
