package strict

type AlignAndSortWithOrderExample struct {
	Foo    int `binding:"required" gorm:"column:foo" json:"foo,omitempty" validate:"required" xml:"baz" yaml:"bar" zip:"foo"` // want `tag is not aligned, should be: json:"foo,omitempty" yaml:"bar" xml:"baz" binding:"required" gorm:"column:foo" validate:"required" zip:"foo"`
	Bar    int `binding:"required" gorm:"column:bar" json:"bar,omitempty" validate:"required" xml:"bar" yaml:"foo" zip:"bar"` // want `tag is not aligned, should be: json:"bar,omitempty" yaml:"foo" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
	FooBar int `binding:"required" gorm:"column:bar" json:"bar,omitempty" validate:"required" xml:"bar" yaml:"foo" zip:"bar"` // want `tag is not aligned, should be: json:"bar,omitempty" yaml:"foo" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
}

type AlignAndSortWithOrderExample2 struct {
	Foo int `binding:"required" gorm:"column:foo"                      validate:"required" xml:"baz" yaml:"bar" zip:"foo"` // want `tag is not aligned, should be:                      yaml:"bar" xml:"baz" binding:"required" gorm:"column:foo" validate:"required" zip:"foo"`
	Bar int `binding:"required" gorm:"column:bar" json:"bar,omitempty" validate:"required" xml:"bar" yaml:"foo"`           // want `tag is not aligned, should be: json:"bar,omitempty" yaml:"foo" xml:"bar" binding:"required" gorm:"column:bar" validate:"required"`
}

type AlignAndSortWithOrderExample3 struct {
	Foo    int `                   gorm:"column:foo"                                                                           zip:"foo"` // want `tag is not aligned, should be:                                                                          gorm:"column:foo"                     zip:"foo"`
	Bar    int `binding:"required" gorm:"column:bar" json:"bar,omitempty" validate:"required" xml:"barxxxxxxxxxxxx" yaml:"foo" zip:"bar"` // want `tag is not aligned, should be: json:"bar,omitempty" yaml:"foo" xml:"barxxxxxxxxxxxx" binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
	FooBar int `binding:"required" gorm:"column:bar" json:"bar,omitempty" validate:"required"                       yaml:"foo" zip:"bar"` // want `tag is not aligned, should be: json:"bar,omitempty" yaml:"foo"                       binding:"required" gorm:"column:bar" validate:"required" zip:"bar"`
}
