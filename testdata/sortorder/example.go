package sortorder

type SortWithOrderExample struct {
	// not aligned but sorted, should not be reported
	Foo    int `xml:"baz" json:"foo,omitempty"     yaml:"bar"     binding:"required"      gorm:"column:foo" validate:"required"     zip:"foo" `
	Bar    int `xml:"bar"        json:"bar,omitempty" yaml:"foo"            gorm:"column:bar"  validate:"required"     zip:"bar" `
	FooBar int `xml:"bar"           json:"bar,omitempty"             yaml:"foo"   gorm:"column:bar"   `
	// aligned but not sorted, should be reported
	BarFoo int `xml:"bar" yaml:"foo" json:"bar,omitempty" gorm:"column:bar" validate:"required" zip:"bar"` // want `xml:"bar" json:"bar,omitempty" yaml:"foo" gorm:"column:bar" validate:"required" zip:"bar"`
	// not aligned but sorted, should trim spaces between tags
	FooBarFoo int `xml:"bar"    json:"bar,omitempty"       yaml:"foo"       gorm:"column:bar" validate:"required" zip:"bar"` // want `xml:"bar" json:"bar,omitempty" yaml:"foo" gorm:"column:bar" validate:"required" zip:"bar"`
}
