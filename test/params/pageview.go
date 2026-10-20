package params

import (
	"github.com/foomo/gocontemplate/test"
)

// PageView holds the parameters of a page view event.
type PageView struct {
	Currency test.Currency `json:"currency,omitempty"`
}
