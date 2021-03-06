package errors

import (
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

const (
	Errors_MissingName = "Helloworld_MissingName"
)

func IsMissingName(err error) bool {
	return errors.Reason(err) == Errors_MissingName
}
