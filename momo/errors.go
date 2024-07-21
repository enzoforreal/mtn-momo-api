package momo

import "errors"

var (
	ErrAuthFailed         = errors.New("authentication failed")
	ErrBalanceFetchFailed = errors.New("failed to fetch account balance")
	ErrRequestToPayFailed = errors.New("request to pay failed")
)
