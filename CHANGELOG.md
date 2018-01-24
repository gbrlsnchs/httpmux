# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [0.6.1] - 2018-01-24
### Fixed
- `Router` and `Subrouter` tests for handling requests with params.

## [0.6.0] - 2018-01-23
### Added
- Mock-up test file.
- Helpers test file.

### Changed
- Router tests.
- Structure stored in the radix tree (for performance reasons).
- Key used to retrieve params in the request's context.

## [0.5.1] - 2017-12-21
### Fixed
- `Cancel` not canceling middleware stacks.
- `404 Not Found` not being sent when no handler is found.

## [0.5.0] - 2017-12-21
### Added
- A method for printing the multiplexer structure.
- New test cases.

### Fixed
- `Subrouter.Use` now correctly spreads handlers to the parent subrouter's handlers slice.

## [0.4.0] - 2017-12-17
### Added
- Makefile support.
- `Router` and `Subrouter` structures.
- Tests for the new structures.
- Cancel method for making short-circuiting middlewares easier.
- Documentation for the new structs, methods and functions.

### Changed
- Travis CI script.
- The PATRICIA tree that holds handlers is now a radix tree.
- `doc.go` now says about middlewares.
- Example now shows middlewares usage.
- README file.

### Removed
- Travis CI `goimports` script.
- `mux.go`, `mux_test.go`, `submux.go` and `path.go` files.

## [0.3.0] - 2017-11-11
### Added
- Submux.
- PATRICIA tree algorithm to resolve the HTTP requests.

### Changed
- Mux's design.
- Example.
- Update this file to use "changelog" in lieu of "change log".

### Removed
- Handler file.
- Comparison to other packages in benchmark.

## [0.2.0] - 2017-09-26
### Added
- Request header filter.

### Changed
- Hash maps are now allocated only when used.

## 0.1.0 - 2017-09-25
### Added
- This changelog file.
- README file.
- MIT License.
- Travis CI configuration file and scripts.
- Git ignore file.
- Editorconfig file.
- This package's source code, including examples and tests.
- Go dep files.

[0.6.1]: https://github.com/gbrlsnchs/httpmux/compare/v0.6.0...v0.6.1
[0.6.0]: https://github.com/gbrlsnchs/httpmux/compare/v0.5.1...v0.6.0
[0.5.1]: https://github.com/gbrlsnchs/httpmux/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/gbrlsnchs/httpmux/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/gbrlsnchs/httpmux/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/gbrlsnchs/httpmux/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/gbrlsnchs/httpmux/compare/v0.1.0...v0.2.0
