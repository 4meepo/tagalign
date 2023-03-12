# Golang Tag Align Linter

This linter is used to align golang struct's tags. It is built for integrating with [golangci-lint](https://golangci-lint.run/usage/quick-start/).

For example:

* Before aligned:

    ```go
    type Foo struct {
        Id      string `json:"id" yaml:"Id"`
        Name    string `json:"name" yaml:"name"`
        Address string `json:"Address" yaml:"Address"`
    }
    ```

* After aligned:

    ```go
    type Foo struct {
        Id      string `json:"id"      yaml:"Id"`
        Name    string `json:"name"    yaml:"name"`
        Address string `json:"Address" yaml:"Address"`
    }
    ```

## Reference

[Golang AST Visualizer](http://goast.yuroyoro.net/)

[Create New Golang CI Linter](https://golangci-lint.run/contributing/new-linters/)

[Autofix Example](https://github.com/golangci/golangci-lint/pull/2450/files)
