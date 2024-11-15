package ecode

import (
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// Error struct
type Error struct {
	Status
	cause error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d msg = %s metadata = %v cause = %v", e.Status.Code, e.Status.Msg, e.Metadata, e.cause)
}

// Code returns the code of the error.
func (e *Error) Code() int32 { return e.Status.Code }

// Msg returns the msg of the error.
func (e *Error) Msg() string { return e.Status.Msg }

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Status.Code == e.Status.Code
	}
	return false
}

// Equal matches error from code.
func (e *Error) Equal(code int) bool {
	se := &Error{Status: Status{
		Code: int32(code),
	}}
	return se.Status.Code == e.Status.Code
}

// GRPCStatus returns the Status represented by error.
func (e *Error) GRPCStatus() *status.Status {
	gs, _ := status.New(DefaultConverter.ToGRPCCode(int(e.Status.Code)), e.Status.Msg).
		WithDetails(&errdetails.ErrorInfo{Metadata: e.Metadata})
	return gs
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := DeepClone(e)
	err.cause = cause
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := DeepClone(e)
	err.Metadata = md
	return err
}

// ============================================================================================================

// New returns an error object for the code, msg.
func New(code int32, msg string) *Error {
	return &Error{Status: Status{
		Code: code,
		Msg:  msg,
	}}
}

// DeepClone deep clone error to a new error.
func DeepClone(err *Error) *Error {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &Error{
		cause: err.cause,
		Status: Status{
			Code:     err.Status.Code,
			Msg:      err.Status.Msg,
			Metadata: metadata,
		},
	}
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return Success
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return New(UnknownCode, err.Error())
	}
	ret := New(DefaultConverter.FromGRPCCode(gs.Code()), gs.Message())
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}

// AnalyseError analyse error info
func AnalyseError(err error) (e2 *Error) {
	if err == nil {
		return Success
	}
	if errors.As(err, &e2) {
		return e2
	}
	return errStringToErrorV2(err.Error())
}

func errStringToErrorV2(e string) *Error {
	if e == "" {
		return Success
	}
	i, err := strconv.Atoi(e)
	if err != nil {
		return New(-1, e)
	}
	return New(int32(i), e)
}
