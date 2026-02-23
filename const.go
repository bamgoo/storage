package storage

import "errors"

const NAME = "STORAGE"

var (
	errInvalidConnection = errors.New("invalid storage connection")
	errInvalidCode       = errors.New("invalid storage code")
	errInvalidHandler    = errors.New("invalid storage handler")
	errBrowseUnsupported = errors.New("storage browse not supported")
)
