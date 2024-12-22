package issues6

type FooBar struct {
	Foo    int    `json: "foo"    validate:"required"` // want `bad syntax for struct tag value`
	Bar    string `json:bar`                           // want `bad syntax for struct tag value`
	FooFoo int8   `json:"foo_foo" validate:"required"`
	BarBar int    `json:"bar_bar" validate:"required"`
}

type FooBar2 struct {
	Foo int `json:"foo" validate:"required"`

	FooFoo int8   `json:"foo_foo"`
	BarBar int    `json:"bar_bar" validate:"required"`
	XXX    int    `json:"xxx"     validate:"required"`
	Bar    string `json:bar` // want `bad syntax for struct tag value`
}
