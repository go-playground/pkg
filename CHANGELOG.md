# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [5.30.0] - 2024-06-01
### Changed
- Changed NanoTome to not use linkname due to Go1.23 upcoming breaking changes. 

## [5.29.1] - 2024-04-04
### Fixed
- Added HTTP 404 to non retryable status codes.

## [5.29.0] - 2024-03-24
### Added
- `asciiext` package for ASCII related functions.
- `errorsext.Retrier` configurable retry helper for any fallible operation.
- `httpext.Retrier` configurable retry helper for HTTP requests and parsing of responses.
- `httpext.DecodeResponseAny` non-generic helper for decoding HTTP responses.
- `httpext.HasRetryAfter` helper for checking if a response has a `Retry-After` header and returning duration to wait.

## [5.28.1] - 2024-02-14
### Fixed
- Additional supported types, cast to `sql.Valuer` supported types, they need to be returned to the driver for evaluation.

## [5.28.0] - 2024-02-13
### Added
- Additionally supported types, cast to `sql.Valuer` supported types.

### Changed
- Option scan to take advantage of new `sql.Null` and `reflect.TypeFor` for go1.22+.
- `BytesToString` & `StringToBytes` to use `unsafe.String` & `unsafe.Slice` for go1.21+.

### Deprecated
- `mathext.Min` & `mathext.Max` in favour of std lib min & max.

### Fixed
- Some documentation typos.

## [5.27.0] - 2024-01-29
### Changed
- `sliceext.Retain` & `sliceext.Filter` to not shuffle data in the underlying slice array but create new slice referencing the data instead. In practice, it can cause unexpected behaviour and users expectations not met when the same data is also referenced elsewhere. If anyone still requires a `shuffle` implementation for efficiency I'd be happy to add a separate function for that as well.

## [5.26.0] - 2024-01-28
### Added
- `stringsext.Join` a more ergonomic way to join strings with a separator when you don't have a slice of strings.

## [5.25.0] - 2024-01-22
### Added
- Add additional `Option.Scan` type support for `sql.Scanner` interface of Uint, Uint16, Uint32, Uint64, Int, Int, Int8, Float32, []byte, json.RawValue.

## [5.24.0] - 2024-01-21
### Added
- `appext` package for application level helpers. Specifically added setting up os signal trapping and cancellation of context.Context.

## [5.23.0] - 2024-01-14
### Added
- `And` and `AndThen` functions to `Option` & `Result` types.

## [5.22.0] - 2023-10-18
### Added
 - `UnwrapOr`, `UnwrapOrElse` and `UnwrapOrDefault` functions to `Option` & `Result` types.

## [5.21.3] - 2023-10-11
### Fixed
- Fix SQL Scanner interface not returning None for Option when source data is nil.

## [5.21.2] - 2023-07-13
### Fixed
- Updated default form/url.Value encoder/decoder with fix for bubbling up invalid array index values.

## [5.21.1] - 2023-06-30
### Fixed
- Instant type to not be wrapped in a struct but a type itself.

## [5.21.0] - 2023-06-30
### Added
- Instant type to make working with monotonically increasing times more convenient. 

## [5.20.0] - 2023-06-17
### Added
- Expanded Option type SQL Value support to handle value custom types and honour the `driver.Valuer` interface.

### Changed
- Option sql.Scanner to support custom types.

## [5.19.0] - 2023-06-14
### Added
- strconvext.ParseBool(...) which is a drop-in replacement for the std lin strconv.ParseBool(..) with a few more supported values.
- Expanded Option type SQL Scan support to handle Scanning to an Interface, Struct, Slice, Map and anything that implements the sql.Scanner interface.

## [5.18.0] - 2023-05-21
### Added
- typesext.Nothing & valuesext.Nothing for better clarity in generic params and values that represent struct{}. This will provide better code readability and intent.

## [5.17.2] - 2023-05-09
### Fixed
- Prematurely closing http.Response Body before error with it can be intercepted for ErrUnexpectedResponse. 

## [5.17.1] - 2023-05-09
### Fixed
- ErrRetryableStatusCode passing the *http.Response to have access to not only the status code but headers etc. related to retrying.
- Added ErrUnexpectedResponse to pass back when encountering an unexpected response code to allow the caller to decide what to do.

## [5.17.0] - 2023-05-08
### Added
- bytesext.Bytes alias to int64 for better code clarity.
- errorext.DoRetryable(...) building block for automating retryable errors.
- sqlext.DoTransaction(...) building block for abstracting away transactions.
- httpext.DoRetryableResponse(...) & httpext.DoRetryable(...) building blocks for automating retryable http requests.
- httpext.DecodeResponse(...) building block for decoding http responses.
- httpext.ErrRetryableStatusCode error for retryable http status code detection and handling.
- errorsext.ErrMaxAttemptsReached error for retryable retryable logic & reuse.

## [5.16.0] - 2023-04-16
### Added
- sliceext.Reverse(...)

## [5.15.2] - 2023-03-06
### Remove
- Unnecessary second type param for Mutex2.

## [5.15.1] - 2023-03-06
### Fixed
- New Mutex2 functions and guards; checked in the wrong code accidentally last commit.

## [5.15.0] - 2023-03-05
### Added
- New Mutex2 and RWMutex2 which corrects the original Mutex's design issues.
- Deprecation warning for original Mutex usage.

## [5.14.0] - 2023-02-25
### Added
- Added `timext.NanoTime` for fast low level monotonic time with nanosecond precision.

[Unreleased]: https://github.com/go-playground/pkg/compare/v5.30.0...HEAD
[5.30.0]: https://github.com/go-playground/pkg/compare/v5.29.1..v5.30.0
[5.29.1]: https://github.com/go-playground/pkg/compare/v5.29.0..v5.29.1
[5.29.0]: https://github.com/go-playground/pkg/compare/v5.28.1..v5.29.0
[5.28.1]: https://github.com/go-playground/pkg/compare/v5.28.0..v5.28.1
[5.28.0]: https://github.com/go-playground/pkg/compare/v5.27.0..v5.28.0
[5.27.0]: https://github.com/go-playground/pkg/compare/v5.26.0..v5.27.0
[5.26.0]: https://github.com/go-playground/pkg/compare/v5.25.0..v5.26.0
[5.25.0]: https://github.com/go-playground/pkg/compare/v5.24.0..v5.25.0
[5.24.0]: https://github.com/go-playground/pkg/compare/v5.23.0..v5.24.0
[5.23.0]: https://github.com/go-playground/pkg/compare/v5.22.0..v5.23.0
[5.22.0]: https://github.com/go-playground/pkg/compare/v5.21.3..v5.22.0
[5.21.3]: https://github.com/go-playground/pkg/compare/v5.21.2..v5.21.3
[5.21.2]: https://github.com/go-playground/pkg/compare/v5.21.1..v5.21.2
[5.21.1]: https://github.com/go-playground/pkg/compare/v5.21.0..v5.21.1
[5.21.0]: https://github.com/go-playground/pkg/compare/v5.20.0..v5.21.0
[5.20.0]: https://github.com/go-playground/pkg/compare/v5.19.0..v5.20.0
[5.19.0]: https://github.com/go-playground/pkg/compare/v5.18.0..v5.19.0
[5.18.0]: https://github.com/go-playground/pkg/compare/v5.17.2..v5.18.0
[5.17.2]: https://github.com/go-playground/pkg/compare/v5.17.1..v5.17.2
[5.17.1]: https://github.com/go-playground/pkg/compare/v5.17.0...v5.17.1
[5.17.0]: https://github.com/go-playground/pkg/compare/v5.16.0...v5.17.0
[5.16.0]: https://github.com/go-playground/pkg/compare/v5.15.2...v5.16.0
[5.15.2]: https://github.com/go-playground/pkg/compare/v5.15.1...v5.15.2
[5.15.1]: https://github.com/go-playground/pkg/compare/v5.15.0...v5.15.1
[5.15.0]: https://github.com/go-playground/pkg/compare/v5.14.0...v5.15.0
[5.14.0]: https://github.com/go-playground/pkg/commit/v5.14.0