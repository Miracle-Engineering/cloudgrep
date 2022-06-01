package util

type UserError struct {
	msg string
}

func (e UserError) Error() string {
	return e.msg
}

func NewUserError(m string) UserError {
	return UserError{msg: m}
}
