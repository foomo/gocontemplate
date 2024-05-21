package params

import (
	"github.com/foomo/gocontemplate/test"
)

type PageView struct {
	Currency test.Currency `json:"currency,omitempty"`
}
