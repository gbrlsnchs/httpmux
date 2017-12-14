# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased
### Added
- Makefile support.

### Changed
- Travis CI script.

### Removed
- Travis CI `goimports` script.

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

[Unreleased]: https://github.com/gbrlsnchs/httpmux/compare/v0.3.0...HEAD
[0.3.0]: https://github.com/gbrlsnchs/httpmux/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/gbrlsnchs/httpmux/compare/v0.1.0...v0.2.0
