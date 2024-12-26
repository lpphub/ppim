package errs

import "github.com/lpphub/golib/web"

var (
	ErrServerInternal = web.Error{
		Code: -1,
		Msg:  "server internal error",
	}

	ErrInvalidParam = web.Error{
		Code: 1001,
		Msg:  "invalid param",
	}

	ErrRecordNotFound = web.Error{
		Code: 2001,
		Msg:  "record not found",
	}
)
