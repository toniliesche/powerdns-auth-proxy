# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.1] - 2025-05-04
### Fixed

- Fixed issue with PowerDNS cache controller action not using correct user authorization

## [1.3.0] - 2025-04-13
### Added

- Added zerolog logger
- Added trace logs

### Changed

- Major Code Style Overhaul

## [1.2.0] - 2025-04-13
### Changed

- Changed way to authorize users (assume implicit reading of TLD domain)

### Fixed

- Fixed tests

## [1.1.1] - 2025-04-12
### Added

- Added /api endpoint to retrieve api version of PowerDNS server

### Fixed

- Filter X-Api-Key header from requests sent to PowerDNS API Proxy

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