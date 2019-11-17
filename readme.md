# `github.com/robertrossmann/actions/toolkit`

[![GoDoc][godoc-badge]][godoc-home] [![GitHub Actions][verify-badge]][verify-action]

> A basic toolkit for working with GitHub Actions.

## About

This toolkit allows you to perform basic operations from within a GitHub Action context, such as

- annotating a pull request at certain lines
- Reading metadata about the current Action run
- Setting an Action's output

It was heavily inspired by the official toolkit for JavaScript: [`@actions/toolkit`][js-toolkit]. This package implements most of the `core` part of that package.

## License

See the [LICENSE](LICENSE) file for more information.

[js-toolkit]: https://github.com/actions/toolkit
[godoc-home]: https://godoc.org/github.com/robertrossmann/actions/toolkit
[godoc-badge]: https://godoc.org/github.com/robertrossmann/actions/toolkit?status.svg
[verify-badge]: https://github.com/robertrossmann/actions/workflows/verify/badge.svg
[verify-action]: https://github.com/robertrossmann/actions/actions
