package stopalignthresholdstrict

type Foo struct {
	Foo string `json:"foo"   description:"this is foo, it's a long str" yaml:"foo"    xml:"foo"  required:"true"` // want `json:"foo" description:"this is foo, it's a long str" yaml:"foo" required:"true" xml:"foo"`
	Bar string `json:"bar"  yaml:"bar"  required:"true"`                                                          // want `json:"bar" yaml:"bar" required:"true"`
}

//
