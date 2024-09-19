package stopalignthreshold

type FooBar struct {
	Foo    string `json:"foo" yaml:"foo" required:"true" description:"this is a str" xml:"foo"` // want `json:"foo"    yaml:"foo"    required:"true" description:"this is a str" xml:"foo"`
	FooBar string `json:"fooBar" yaml:"fooBar" description:"this is a long long str str str" xml:"fooBar"`
}

type Foo struct {
	Foo string `json:"foo"   description:"this is foo" yaml:"foo" required:"true"  xml:"foo"` // want `json:"foo" description:"this is foo" yaml:"foo" required:"true" xml:"foo"`
	Bar string `json:"bar"  yaml:"bar"`                                                       // want `json:"bar" yaml:"bar"`
}
