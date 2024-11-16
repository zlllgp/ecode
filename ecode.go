package ecode

import (
	"errors"
	"fmt"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// ErrorNo struct
type ErrorNo struct {
	Status
	cause error
}

func (e *ErrorNo) Error() string {
	return fmt.Sprintf("error: code = %d msg = %s metadata = %v cause = %v", e.Status.Code, e.Status.Msg, e.Metadata, e.cause)
}

// Code returns the code of the error.
func (e *ErrorNo) Code() int64 { return e.Status.Code }

// Msg returns the msg of the error.
func (e *ErrorNo) Msg() string { return e.Status.Msg }

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *ErrorNo) Unwrap() error { return e.cause }

// Is matches each error in the chain with the target value.
func (e *ErrorNo) Is(err error) bool {
	if se := new(ErrorNo); errors.As(err, &se) {
		return se.Status.Code == e.Status.Code
	}
	return false
}

// Equal matches error from code.
func (e *ErrorNo) Equal(code int) bool {
	se := &ErrorNo{Status: Status{
		Code: int64(code),
	}}
	return se.Status.Code == e.Status.Code
}

// GRPCStatus returns the Status represented by error.
func (e *ErrorNo) GRPCStatus() *status.Status {
	gs, _ := status.New(DefaultConverter.ToGRPCCode(int(e.Status.Code)), e.Status.Msg).
		WithDetails(&errdetails.ErrorInfo{Metadata: e.Metadata})
	return gs
}

// WithCause with the underlying cause of the error.
func (e *ErrorNo) WithCause(cause error) *ErrorNo {
	err := DeepClone(e)
	err.cause = cause
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *ErrorNo) WithMetadata(md map[string]string) *ErrorNo {
	err := DeepClone(e)
	err.Metadata = md
	return err
}

// ============================================================================================================

// NewErrNo returns an error object for the code, msg.
func NewErrNo(code int64, msg string) *ErrorNo {
	return &ErrorNo{Status: Status{
		Code: code,
		Msg:  msg,
	}}
}

// DeepClone deep clone error to a new error.
func DeepClone(err *ErrorNo) *ErrorNo {
	if err == nil {
		return nil
	}
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &ErrorNo{
		cause: err.cause,
		Status: Status{
			Code:     err.Status.Code,
			Msg:      err.Status.Msg,
			Metadata: metadata,
		},
	}
}

// FromError try to convert an error to *ErrorNo.
// It supports wrapped errors.
func FromError(err error) *ErrorNo {
	if err == nil {
		return Success
	}
	if se := new(ErrorNo); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if !ok {
		return NewErrNo(UnknownCode, err.Error())
	}
	ret := NewErrNo(DefaultConverter.FromGRPCCode(gs.Code()), gs.Message())
	for _, detail := range gs.Details() {
		switch d := detail.(type) {
		case *errdetails.ErrorInfo:
			return ret.WithMetadata(d.Metadata)
		}
	}
	return ret
}

// AnalyseError analyse error info
func AnalyseError(err error) (e *ErrorNo) {
	if err == nil {
		return Success
	}
	if errors.As(err, &e) {
		return e
	}
	return errStringToError(err.Error())
}

func errStringToError(e string) *ErrorNo {
	if e == "" {
		return Success
	}
	i, err := strconv.Atoi(e)
	if err != nil {
		return NewErrNo(-1, e)
	}
	return NewErrNo(int64(i), e)
}
