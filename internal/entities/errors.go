package entities

import "errors"

var (
	ErrNotFound                = errors.New("entity not found")
	ErrInvalidName             = errors.New("invalid argument")
	ErrRoomNotAvailable        = errors.New("room not available")
	ErrStartDateIsAfterEndDate = errors.New("start date is after end date")
	ErrNameIsRequired          = errors.New("name is required")
	ErrNameIsTooLong           = errors.New("name is too long")
)
