package dpet

import "errors"

var (
	WrongFileTypeError   = errors.New("not dpet file or file damage")
	UnmarshalError       = errors.New("cannot unmarshal file header content")
	UnknownMarshalMethod = errors.New("unknown marshal method")
	UnknownFileType      = errors.New("unknown file type")
	UnknownDrive         = errors.New("unknown unknown drive")
)
