package client

import "errors"

var (
	Err401StatusCode        = errors.New("unauthorized - 401")
	ErrWrongStatusCode      = errors.New("wrong status code")
	ErrCantCreateRequest    = errors.New("can't create request")
	ErrAuthInfoNotSpecified = errors.New("auth info wasn't specified")
	ErrCantUpdateToken      = errors.New("can't update token")
)
