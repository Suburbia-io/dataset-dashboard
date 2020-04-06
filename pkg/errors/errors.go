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

	JobQueueFull = NewErr("JobQueueFull")

	HttpBadRequestArgs = NewErr("HttpBadRequestArgs")
	HttpBadMethod      = NewErr("HttpBadMethod")
	HttpBadContentType = NewErr("HttpBadContentType")
	HttpNotAuthorized  = NewErr("HttpNotAuthorized")
	HttpForbidden      = NewErr("HttpForbidden")
	NotAuthorized      = NewErr("NotAuthorized")

	InvalidSessionToken = NewErr("InvalidSessionToken")
	SessionNotFound     = NewErr("SessionNotFound")
	AuthTokenNotFound   = NewErr("AuthTokenNotFound")

	UserNotAllowedPassword = NewErr("UserNotAllowedPassword")

	EventInvalidName = NewErr("EventInvalidName")

	IXRuleNotFound        = NewErr("IXRuleNotFound")
	IXRuleInvalidGroup    = NewErr("IXRuleInvalidGroup")
	IXRuleInvalidLabel    = NewErr("IXRuleInvalidLabel")
	IXRuleDuplicate       = NewErr("IXRuleDuplicate")
	IXRuleInvalidIncludes = NewErr("IXRuleInvalidIncludes")
	IXRuleInvalidExcludes = NewErr("IXRuleInvalidExcludes")

	BrandNotFound           = NewErr("BrandNotFound")
	BrandInvalidLabel       = NewErr("BrandInvalidLabel")
	BrandInvalidDescription = NewErr("BrandInvalidDescription")
	BrandDuplicateSlug      = NewErr("BrandDuplicateSlug")

	CorporationNotFound        = NewErr("CorporationNotFound")
	CorporationInvalidSymbol   = NewErr("CorporationInvalidSymbol")
	CorporationInvalidName     = NewErr("CorporationInvalidName")
	CorporationDuplicateSymbol = NewErr("CorporationDuplicateSymbol")
	CorporationDuplicateISIN   = NewErr("CorporationDuplicateISIN")
	CorporationInvalidISIN     = NewErr("CorporationInvalidISIN")
	CorporationDuplicateCUSIP  = NewErr("CorporationDuplicateCUSIP")
	CorporationInvalidCUSIP    = NewErr("CorporationInvalidCUSIP")
	CorporationDuplicate       = NewErr("CorporationDuplicate")

	MappingNotFound        = NewErr("MappingNotFound")
	MappingIllegalChild    = NewErr("MappingIllegalChild")
	MappingIllegalParent   = NewErr("MappingIllegalParent")
	MappingInvalidLocation = NewErr("MappingInvalidLocation")
	MappingChildNotFound   = NewErr("MappingChildNotFound")
	MappingParentNotFound  = NewErr("MappingParentNotFound")
	MappingInvalidStart    = NewErr("MappingInvalidStart")

	DatasetInvalidName = NewErr("DatasetInvalidName")
	DatasetNotFound    = NewErr("DatasetNotFound")
	DatasetDuplicate   = NewErr("DatasetDuplicate")

	CustomerInvalidName = NewErr("CustomerInvalidName")
	CustomerNotFound    = NewErr("CustomerNotFound")

	LocationNotFound      = NewErr("LocationNotFound")
	LocationAlreadyExists = NewErr("LocationAlreadyExists")

	IllegalOffset      = NewErr("IllegalOffset")
	IllegalLimit       = NewErr("IllegalLimit")
	IllegalSort        = NewErr("IllegalSort")
	InvalidQueryParams = NewErr("InvalidQueryParams")

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
