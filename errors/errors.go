package errors

import pkgErrors "errors"

var ErrEntityNotFound = pkgErrors.New("error entity not found")
var ErrMissingParameters = pkgErrors.New("error missing parameters")
