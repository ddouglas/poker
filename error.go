package poker

import "fmt"

var (
	ErrInternalServerErrorContactDeveloper = fmt.Errorf("Internal Server Error, please try again, if error persist, contact the developer")
)

type ErrorVisibility uint

const (
	PublicErrorVisibility ErrorVisibility = iota
	InternalErrorVisibility
)

type Error interface {
	error
	// Error() error
	Public() error
}

var _ Error = (*ValidationError)(nil)

type ValidationError struct {
	Visibility ErrorVisibility
	Message    error
}

func (e ValidationError) Public() error {
	if e.Visibility == PublicErrorVisibility {
		return e.Message
	}

	return fmt.Errorf("Internal Server Error")
}

func (e ValidationError) Error() string {
	if e.Message == nil {
		return ""
	}

	return e.Message.Error()
}
