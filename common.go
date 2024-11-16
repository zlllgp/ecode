package ecode

import "net/http"

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = http.StatusInternalServerError
	// ClientClosed is non-standard http status code,
	// which defined by nginx.
	// https://httpstatus.in/499/
	ClientClosed = 499
)

var (
	Success               = NewErrNo(http.StatusOK, "success")
	RequestErr            = NewErrNo(http.StatusBadRequest, "request param error")
	UnauthorizedErr       = NewErrNo(http.StatusUnauthorized, "sign error")
	ForbiddenErr          = NewErrNo(http.StatusForbidden, "no auth")
	NotFoundErr           = NewErrNo(http.StatusNotFound, "resource not found")
	TooManyRequestErr     = NewErrNo(http.StatusTooManyRequests, "rate limit exceeded")
	ServerErr             = NewErrNo(http.StatusInternalServerError, "server error")
	BadGatewayErr         = NewErrNo(http.StatusBadGateway, "service offline, unavailable")
	ServiceUnavailableErr = NewErrNo(http.StatusServiceUnavailable, "service protected, unavailable")
)
