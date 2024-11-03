package errs

import "github.com/lpphub/golib/render"

var (
	ErrServerInternal = render.Error{
		Code: -1,
		Msg:  "server internal error",
	}

	ErrInvalidParam = render.Error{
		Code: 1001,
		Msg:  "invalid param",
	}

	ErrRecordNotFound = render.Error{
		Code: 2001,
		Msg:  "record not found",
	}
)
