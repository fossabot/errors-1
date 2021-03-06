package errors

import errors "github.com/pkg/errors"

type causer interface {
	Cause() error
}

type kinder interface {
	Kind() Kind
}

type operation interface {
	Op() Op
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	return errors.Cause(err)
}

// Unwrap returns the first occurrence of the underlying *Error type
// using causer or create a new *Error type containing the original error's message
func Unwrap(err error) *Error {
	if err == nil {
		return nil
	}
	for err != nil {
		e, ok := err.(*Error)
		if ok {
			return e
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return &Error{
		msg:   "Internal error or inconsistency",
		kind:  Internal,
		cause: err,
	}
}

// KindOf returns the underlying kind type of the error, if possible.
// An error value has a kind if it implements the following
// interface:
//
//     type kinder interface {
//            Kind() Kind
//     }
//
// If the error does not implement Kind, the original error will
// be returned. If the error is nil, Internal will be returned without further
// investigation.
func KindOf(err error) Kind {
	for err != nil {
		kind, ok := err.(kinder)
		if ok {
			return kind.Kind()
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return Internal
}

// IsKind returns a boolean signifying if the error
// matches the Kind provided in the input
func IsKind(err error, k Kind) bool {
	return KindOf(err) == k
}
