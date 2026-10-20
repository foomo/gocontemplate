package service

import (
	"github.com/foomo/gocontemplate/test/event"
)

type Event[T any] interface {
	PageView() event.PageView
}
