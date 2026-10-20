package event

import (
	"github.com/foomo/gocontemplate/test"
	"github.com/foomo/gocontemplate/test/params"
)

// PageView is the page view event.
type PageView test.Event[params.PageView]
