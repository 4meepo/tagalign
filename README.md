# Golang Tag Align Linter

This linter is used to align golang struct's tags. It is built for integrating with [golangci-lint](https://golangci-lint.run/usage/quick-start/).

For example:

* Before aligned:

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

* After aligned:

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

## Reference

[Golang AST Visualizer](http://goast.yuroyoro.net/)

[Create New Golang CI Linter](https://golangci-lint.run/contributing/new-linters/)

[Autofix Example](https://github.com/golangci/golangci-lint/pull/2450/files)
