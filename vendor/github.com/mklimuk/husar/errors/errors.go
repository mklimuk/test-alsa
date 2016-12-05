package errors

import e "errors"

//ErrorType represents a type of error
type ErrorType int

// predefined error types
const (
	NotFound ErrorType = iota
	BadRequest
	Unrecognized
)

//Error is a generic Husar error
type Error struct {
	error
	Type ErrorType
	Ctx  map[string]string
}

//NewError is Error constructor
func NewError(msg string, t ErrorType) error {
	return Error{
		error: e.New(msg),
		Type:  t,
	}
}

//NewWithCtx provides a way to introduce a context to the custom error
func NewWithCtx(msg string, t ErrorType, ctx map[string]string) error {
	return Error{
		error: e.New(msg),
		Type:  t,
		Ctx:   ctx,
	}
}

//GetCtx returns error context
func GetCtx(e error) *map[string]string {
	var herr Error
	var err bool
	if herr, err = e.(Error); err == false {
		return nil
	}
	return &herr.Ctx
}

//GetType can be used in switch statements whenever several error types are possible
func GetType(e error) ErrorType {
	var herr Error
	var err bool
	if herr, err = e.(Error); err == false {
		return Unrecognized
	}
	return herr.Type
}

//IsType checks whether an error is a Husar error of given type
func IsType(e error, t ErrorType) bool {
	if herr, err := e.(Error); err == false {
		return false
	} else if herr.Type == t {
		return true
	}
	return false
}
