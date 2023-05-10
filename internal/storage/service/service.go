// Package service provides a high flexibility when choosing a database.
// In the future, we can easily add nosql database support.
//
// For now, there is only one error, but when we have more database,
// we will be able to add more general stuff in this package.
package service

import "errors"

// ErrNotFound is returned when a resource cannot be found.
var ErrNotFound = errors.New("not found")
