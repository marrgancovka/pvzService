package pvz

import "errors"

var (
	ErrInaccessibleCity     = errors.New("inaccessible city")
	ErrAlreadyExists        = errors.New("pvz with this id already exists")
	ErrNoOpenReception      = errors.New("no open reception found")
	ErrNoClosedReception    = errors.New("no closed reception found")
	ErrPvzNotExists         = errors.New("pvz does not exist")
	ErrIncorrectProductType = errors.New("incorrect product type")
	ErrNoProduct            = errors.New("no product found")
)
