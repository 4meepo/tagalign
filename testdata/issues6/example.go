package issues6

type FooBar struct {
	Foo    int    `json: "foo"    validate:"required"` // want `bad syntax for struct tag value`
	Bar    string `json:bar`                           // want `bad syntax for struct tag value`
	FooFoo int8   `json:"foo_foo" validate:"required"`
	BarBar int    `json:"bar_bar" validate:"required"`
}
