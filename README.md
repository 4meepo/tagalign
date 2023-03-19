# Golang Tag Align Linter

This linter is used to align golang struct's tags. It is built for integrating with [golangci-lint](https://golangci-lint.run/usage/quick-start/).

For example:

Before aligned:

```go
type FooBar struct {
    Foo    int    `json:"foo" validate:"required"`
    Bar    string `json:"bar" validate:"required"`
    FooFoo int8   `json:"foo_foo" validate:"required"`
    BarBar int    `json:"bar_bar" validate:"required"`
    FooBar struct {
    Foo    int    `json:"foo" yaml:"foo" validate:"required"`
    Bar222 string `json:"bar222" validate:"required" yaml:"bar"`
    } `json:"foo_bar" validate:"required"`
    BarFoo    string `json:"bar_foo" validate:"required"`
    BarFooBar string `json:"bar_foo_bar" validate:"required"`
}
```

After aligned:

```go
type FooBar struct {
    Foo    int    `json:"foo"     validate:"required"`
    Bar    string `json:"bar"     validate:"required"`
    FooFoo int8   `json:"foo_foo" validate:"required"`
    BarBar int    `json:"bar_bar" validate:"required"`
    FooBar struct {
        Foo    int    `json:"foo"    yaml:"foo"          validate:"required"`
        Bar222 string `json:"bar222" validate:"required" yaml:"bar"`
    } `json:"foo_bar" validate:"required"`
    BarFoo    string `json:"bar_foo"     validate:"required"`
    BarFooBar string `json:"bar_foo_bar" validate:"required"`
}
```

In addition to alignment, it can also sort tags with fixed order. For example, if we enable auto-sort with fixed order `json,xml`, the following code

```go
type SortExample struct {
    Foo    int `json:"foo,omitempty" yaml:"bar" xml:"baz" binding:"required" gorm:"column:foo" zip:"foo" validate:"required"`
    Bar    int `validate:"required"  yaml:"foo" xml:"bar" binding:"required" json:"bar,omitempty" gorm:"column:bar" zip:"bar" `
    FooBar int `gorm:"column:bar" validate:"required"   xml:"bar" binding:"required" json:"bar,omitempty"  zip:"bar" yaml:"foo"`
}
```

will be sorted and aligned to:

```go
type SortExample struct {
    Foo    int `json:"foo,omitempty" xml:"baz" binding:"required" gorm:"column:foo" validate:"required" yaml:"bar" zip:"foo"`
    Bar    int `json:"bar,omitempty" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" yaml:"foo" zip:"bar"`
    FooBar int `json:"bar,omitempty" xml:"bar" binding:"required" gorm:"column:bar" validate:"required" yaml:"foo" zip:"bar"`
}
```

The fixed order is `json,xml`, so the tags `json` and `xml` will be sorted and aligned first, and the rest tags will be sorted and aligned in the order of appearance.

## Install

```bash
go install github.com/4meepo/tagalign/cmd/tagalign
```

## Usage

```bash
# basic
tagalign -fix {package path}
# enable auto sort with fixed order
tagalign -fix -auto-sort -fixed-order "json,xml" {package path}
```

## Reference

[Golang AST Visualizer](http://goast.yuroyoro.net/)

[Create New Golang CI Linter](https://golangci-lint.run/contributing/new-linters/)

[Autofix Example](https://github.com/golangci/golangci-lint/pull/2450/files)

[Integraing](https://disaev.me/p/writing-useful-go-analysis-linter/#integrating)
