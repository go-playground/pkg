# pkg

![Project status](https://img.shields.io/badge/version-5.29.0-green.svg)
[![Lint & Test](https://github.com/go-playground/pkg/actions/workflows/go.yml/badge.svg)](https://github.com/go-playground/pkg/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/go-playground/pkg/badge.svg?branch=master)](https://coveralls.io/github/go-playground/pkg?branch=master)
[![GoDoc](https://godoc.org/github.com/go-playground/pkg?status.svg)](https://pkg.go.dev/mod/github.com/go-playground/pkg/v5)
![License](https://img.shields.io/dub/l/vibe-d.svg)

pkg extends the core Go packages with missing or additional functionality built in. All packages correspond to the std go package name with an additional suffix of `ext` to avoid naming conflicts.

## Motivation

This is a place to put common reusable code that is not quite a library but extends upon the core library, or it's failings.

## Install

`go get -u github.com/go-playground/pkg/v5`


## Highlights
- Generic Doubly Linked List.
- Result & Option types
- Generic Mutex and RWMutex.
- Bytes helper placeholders units eg. MB, MiB, GB, ...
- Detachable context.
- Retrier for helping with any fallible operation.
- Proper RFC3339Nano definition.
- unsafe []byte->string & string->[]byte helper functions.
- HTTP helper functions and constant placeholders.
- And much, much more.

## How to Contribute

Make a pull request... can't guarantee it will be added, going to strictly vet what goes in.

## License

<sup>
Licensed under either of <a href="LICENSE-APACHE">Apache License, Version
2.0</a> or <a href="LICENSE-MIT">MIT license</a> at your option.
</sup>

<br>

<sub>
Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in this package by you, as defined in the Apache-2.0 license, shall be
dual licensed as above, without any additional terms or conditions.
</sub>
