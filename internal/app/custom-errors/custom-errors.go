package errors

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultMsg = "internal server error"
const unwrappedErrorMsg = "unknown error"

var (
	// Unknown error. An example of where this error may be returned is
	// if a status value received from another address space belongs to
	// an error-space that is not known in this address space. Also
	// errors raised by APIs that do not return enough error information
	// may be converted to this error.
	Unknown errorBuilder

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument errorBuilder

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	NotFound errorBuilder

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists errorBuilder

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation. It must not be used for rejections
	// caused by exhausting some resource (use ResourceExhausted
	// instead for those errors). It must not be
	// used if the caller cannot be identified (use Unauthenticated
	// instead for those errors).
	PermissionDenied errorBuilder

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Aborted errorBuilder

	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal errorBuilder

	// Unavailable indicates the service is currently unavailable.
	// This is a most likely a transient condition and may be corrected
	// by retrying with a backoff. Note that it is not always safe to retry
	// non-idempotent operations.
	//
	// See litmus test above for deciding between FailedPrecondition,
	// Aborted, and Unavailable.
	Unavailable errorBuilder

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated errorBuilder
)

func init() {
	Unknown = errorBuilder{status: codes.Unknown}
	InvalidArgument = errorBuilder{status: codes.InvalidArgument}
	NotFound = errorBuilder{status: codes.NotFound}
	AlreadyExists = errorBuilder{status: codes.AlreadyExists}
	PermissionDenied = errorBuilder{status: codes.PermissionDenied}
	Aborted = errorBuilder{status: codes.Aborted}
	Internal = errorBuilder{status: codes.Internal}
	Unavailable = errorBuilder{status: codes.Unavailable}
	Unauthenticated = errorBuilder{status: codes.Unauthenticated}
}

// ---------- ERROR BUILDER ---------- //

// stores data about the error type
type errorBuilder struct {
	status codes.Code
}

// returns an error to send to the client via gRPC
func (eb errorBuilder) New(_ context.Context, errorString string) error {
	return newServiceError(eb.status, false, errorString)
}

// wrapping a service error in a client error and logging it
func (eb errorBuilder) NewWrap(_ context.Context, clientError string, serviceErr error) error {
	log.Println(serviceErr)
	return newServiceError(eb.status, true, clientError)
}

// checks the type of error passed
func (eb errorBuilder) IsErr(err error) bool {
	srvErr, ok := err.(*serviceError)
	if ok {
		return srvErr.status == eb.status
	} else {
		return status.Code(err) == eb.status
	}
}

// ---------- ERROR ---------- //

// service error struct which allows you to convert
// the error that occurred to gRPC error
type serviceError struct {
	status    codes.Code
	isWrapped bool
	value     string
}

func newServiceError(status codes.Code, isWrapped bool, value string) *serviceError {
	return &serviceError{status: status, isWrapped: isWrapped, value: value}
}

// implements error interface
func (se serviceError) Error() string {
	return se.value
}

// convert error to gRPC error
// and it doesn't allow sending internal errors to the client
func (se serviceError) toGRPC() error {
	var msg string

	// checks whether the error is internal and whether it was wrapped
	if se.status == codes.Internal && !se.isWrapped {
		msg = defaultMsg
	} else {
		msg = se.Error()
	}

	return status.Error(se.status, msg)
}

// ---------- gRPC WRAPPER ---------- //

// transform received err to gRPC error
func ToGRPC(err error) error {
	if err == nil {
		return nil
	}

	gRPCerr, ok := err.(interface {
		toGRPC() error
	})

	if !ok {
		return status.Error(codes.Unknown, unwrappedErrorMsg)
	}

	return gRPCerr.toGRPC()
}
