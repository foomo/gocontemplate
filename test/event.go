package test

// Event is a named event carrying parameters of type T.
type Event[T any] struct {
	Name   string `json:"name"`
	Params T      `json:"params,omitempty"`
}
