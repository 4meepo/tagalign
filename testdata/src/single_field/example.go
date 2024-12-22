package singlefield

type FooBar struct {
	Foo int    `json:"foo" validate:"required"`
	Bar string `json:"bar" validate:"required"`

	FooFoo int8 `json:"foo_foo"     validate:"required"` // want `json:"foo_foo" validate:"required"`

	BarBar int `json:"bar_bar" validate:"required"`
}
