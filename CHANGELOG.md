# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
## [v0.2.0] - 2021-10-21
Quite a number of breaking changes this release to ensure API consistency
throughout the library.

### Added
- storage/permissions: adds `SignedPermission` enum.
- storage: implements `Stringer` across all packages.

### Changed
- storage/versions: removes camel casing from versions to conform to Go style practices.
- storage/ips: renames package `signedip` to `ips` for consistency.
- storage/protocols: renames types `Protocol` to `SignedProtocol` and `Protocols` to `SignedProtocols`.
- storage/services: renames types `Service` to `SignedService` and `Services` to `SignedServices`.
- storage/resourcetypes: renames types `ResourceType` to `SignedResourceType` and `ResourceTypes` to `SignedResourceTypes`.
- storage/resources: renames types `Resource` to `SignedResource` and `Resources` to `SignedResources`.
- storage: updates `AccountSAS` properties with the new type changes.
- storage: renames `AccountSAS.ApiVersion` to `AccountSAS.APIVersion` to conform to Go style practices.

### Fixed
- storage: ensures all parsing functions trim whitespace and cast inputs to lower case where applicable.

## [v0.1.0] - 2021-09-03
Sassy GA!

[Unreleased]: https://github.com/matthewhartstonge/sassy/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/matthewhartstonge/sassy/releases/tag/v0.2.0
[v0.1.0]: https://github.com/matthewhartstonge/sassy/releases/tag/v0.1.0
