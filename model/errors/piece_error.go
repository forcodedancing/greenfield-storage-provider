package errors

import (
	"errors"
)

var (
	// ErrUnsupportedDataType defines the unsupported payload data type error
	ErrUnsupportedDataType = errors.New("unsupported payload data type")
	// ErrOverflowDataSize defines the payload size overflow error
	ErrOverflowDataSize = errors.New("payload data size overflow")
	// ErrZeroDataSize defines the data size zeroed error
	ErrZeroDataSize = errors.New("payload data size zero")
)
