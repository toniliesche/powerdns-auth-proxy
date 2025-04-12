# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-04-12
### Added

- Added PowerDNS Admin to docker test setup

### Changed

- Filter result from list zones PowerDNS result rather than blocking access to endpoint
- Adjusted error result JSON format to match PowerDNS responses

### Fixed

- Improved naming in docker test setup
- Fixed user mapper in list users API endpoint

## [1.0.0] - 2025-03-11
### Added

- Initial release
- General project setup and build processes
- Added admin API endpoints
- Added authenticated PowerDNS API proxy endpoints
- Added basic CLI management commands