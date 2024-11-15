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
	Success               = New(http.StatusOK, "success")
	RequestErr            = New(http.StatusBadRequest, "request param error")
	UnauthorizedErr       = New(http.StatusUnauthorized, "sign error")
	ForbiddenErr          = New(http.StatusForbidden, "no auth")
	NotFoundErr           = New(http.StatusNotFound, "resource not found")
	TooManyRequestErr     = New(http.StatusTooManyRequests, "rate limit exceeded")
	ServerErr             = New(http.StatusInternalServerError, "server error")
	BadGatewayErr         = New(http.StatusBadGateway, "service offline, unavailable")
	ServiceUnavailableErr = New(http.StatusServiceUnavailable, "service protected, unavailable")
)
