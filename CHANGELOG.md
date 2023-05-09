# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/go-playground/pkg/compare/v5.17.2...HEAD
[5.17.2]: https://github.com/go-playground/pkg/compare/v5.17.1..v5.17.2
[5.17.1]: https://github.com/go-playground/pkg/compare/v5.17.0...v5.17.1
[5.17.0]: https://github.com/go-playground/pkg/compare/v5.16.0...v5.17.0
[5.16.0]: https://github.com/go-playground/pkg/compare/v5.15.2...v5.16.0
[5.15.2]: https://github.com/go-playground/pkg/compare/v5.15.1...v5.15.2
[5.15.1]: https://github.com/go-playground/pkg/compare/v5.15.0...v5.15.1
[5.15.0]: https://github.com/go-playground/pkg/compare/v5.14.0...v5.15.0
[5.14.0]: https://github.com/go-playground/pkg/commit/v5.14.0