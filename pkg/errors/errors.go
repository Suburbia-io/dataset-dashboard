package errors

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

var Unwrap = errors.Unwrap
var Is = errors.Is
var As = errors.As

var (
	AlertCallback func(msg string)
)

var (
	Unexpected = NewErr("Unexpected")

	DBDuplicate      = NewErr("DBDuplicate")
	DBFKey           = NewErr("DBForeignKey")
	DBNullConstraint = NewErr("DBNullConstraint")
	DBNotFound       = NewErr("DBNotFound")
	DBUnknownColumn  = NewErr("DBUnknownColumn")

	HttpBadRequestArgs = NewErr("HttpBadRequestArgs")
	HttpNotAuthorized  = NewErr("HttpNotAuthorized")
	HttpForbidden      = NewErr("HttpForbidden")

	InvalidSessionToken = NewErr("InvalidSessionToken")
	AuthTokenNotFound   = NewErr("AuthTokenNotFound")

	EventInvalidName = NewErr("EventInvalidName")

	BrandInvalidLabel       = NewErr("BrandInvalidLabel")
	BrandInvalidDescription = NewErr("BrandInvalidDescription")

	CorporationInvalidSymbol   = NewErr("CorporationInvalidSymbol")
	CorporationInvalidName     = NewErr("CorporationInvalidName")
	CorporationInvalidISIN     = NewErr("CorporationInvalidISIN")
	CorporationInvalidCUSIP    = NewErr("CorporationInvalidCUSIP")

	CustomerInvalidName = NewErr("CustomerInvalidName")
	CustomerNotFound    = NewErr("CustomerNotFound")

	IllegalOffset      = NewErr("IllegalOffset")
	IllegalLimit       = NewErr("IllegalLimit")

	EmptyString       = NewErr("EmptyString")
	InvalidUUID       = NewErr("InvalidUUID")
	InvalidAuthToken  = NewErr("InvalidAuthToken")
	InvalidHumanName  = NewErr("InvalidHumanName")
	InvalidEmail      = NewErr("InvalidEmail")
	PasswordTooWeak   = NewErr("PasswordTooWeak")
	InvalidSlug       = NewErr("InvalidSlug")
	InvalidPath       = NewErr("InvalidPath")
	InvalidTagType    = NewErr("InvalidTagType")
	InvalidTag        = NewErr("InvalidTag")
	InvalidConfidence = NewErr("InvalidConfidence")
	UpsertNotAllowed  = NewErr("UpsertNotAllowed")
)

// -----------------------------------------------------------------------------

type Error struct {
	code    string
	msg     string
	wrapped error
}

func NewErr(code string) Error {
	return Error{code: code}
}

func (e Error) Error() string {
	if e.msg != "" {
		return fmt.Sprintf(strings.Join([]string{e.code, e.msg}, " - "))
	} else {
		return e.code
	}
}

func (e Error) Unwrap() error {
	return e.wrapped
}

func (e Error) Is(err error) bool {
	if err == nil {
		return false
	}

	switch ee := err.(type) {
	case Error:
		return e.code == ee.code
	case *Error:
		return e.code == ee.code
	default:
		return errors.Is(e.wrapped, err)
	}
}

func (e Error) As(iErr interface{}) bool {
	if err, ok := iErr.(*Error); ok {
		*err = e
		return true
	}
	return errors.As(e.wrapped, iErr)
}

func (e Error) Code() string {
	return e.code
}

func (e Error) Msg() string {
	return e.msg
}

func (e Error) WithMsg(format string, args ...interface{}) Error {
	return Error{
		code:    e.code,
		msg:     fmt.Sprintf(format, args...),
		wrapped: e.wrapped,
	}
}

func (e Error) Wrap(format string, args ...interface{}) Error {
	return Error{
		code:    e.code,
		msg:     e.msg,
		wrapped: fmt.Errorf(format, args...),
	}
}

func (e Error) Alert() Error {
	msg := fmt.Sprintf("ALERT: [%s] %s -- %v", e.code, e.msg, e.wrapped)

	log.Println(msg)

	if AlertCallback != nil {
		AlertCallback(msg)
	}
	return e
}

// TODO: Remove and replace with Wrap.
func (e Error) WithErr(err error) Error {
	return Error{
		code:    e.code,
		msg:     e.msg,
		wrapped: err,
	}
}

// TODO: Remove this function, replace with Wrap / Alert.
func UnexpectedError(err error, format string, args ...interface{}) error {
	args = append(args, err)
	return Unexpected.Wrap(format+" -- %w", args...)
}
