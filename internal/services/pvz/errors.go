package pvz

import "errors"

var (
	ErrInaccessibleCity = errors.New("inaccessible city")
	ErrAlreadyExists    = errors.New("pvz with this id already exists")
	ErrOpenReception    = errors.New("there is an open reception in this pvz")
	ErrPvzNotExists     = errors.New("pvz does not exist")
)
