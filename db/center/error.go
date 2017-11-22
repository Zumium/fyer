package center

import "errors"

var (
	//ErrUnsetField is returned when getting a field from a new doc
	ErrUnsetField = errors.New("field is not set")
)
