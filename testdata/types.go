package testdata

type FooBar struct {
	Foo    int    `  json:"foo"  validate:"required"`
	Bar    string `    json:"___bar___,omitempty"     validate:"required"`
	FooFoo int8   `json:"foo_foo" validate:"required" yaml:"fooFoo"`
	BarBar int    `json:"bar_bar" validate:"required"`
	FooBar struct {
		Foo    int    `json:"foo"    yaml:"foo"          validate:"required"`
		Bar222 string `json:"bar222" validate:"required" yaml:"bar"`
	} `json:"foo_bar" validate:"required"`
	BarFoo    string `json:"bar_foo"     validate:"required"`
	BarFooBar string `json:"bar_foo_bar" validate:"required"`
}
