package test

type Event[T any] struct {
	Name   string `json:"name"`
	Params T      `json:"params,omitempty"`
}
