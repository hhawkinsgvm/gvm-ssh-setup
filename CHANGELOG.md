# Changelog

All notable changes to the GVM SSH Setup Tool project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-01-XX

### Added
- Initial release of GVM SSH Setup Tool
- SSH key generation and management with Ed25519 keys
- SSH alias configuration with organized config.d structure
- Per-directory Git identity using includeIf
- GitLab CE integration with API client
- Group membership verification for access control
- SSH key upload to GitLab accounts
- Interactive wizard with GVM-specific defaults
- Non-interactive setup mode for automation
- Configuration verification and validation
- SSH and Git connectivity testing
- Deploy key support for CI/server environments
- Docker containerization with multi-stage builds
- Private registry support for secure distribution
- Comprehensive documentation and usage examples

### Features
- **Interactive Wizard**: Guided setup with sensible GVM defaults
- **GitLab CE Integration**: Native support for GitLab CE environments
- **Security**: Group membership verification and deploy key isolation
- **Flexibility**: Both interactive and non-interactive modes
- **Testing**: Built-in connectivity and configuration verification
- **Docker Support**: Containerized execution with host filesystem access

### Default Configuration
- Git Host: `gitlab.globalvision.com.au`
- Git Port: `2122`
- Admin Host: `203.32.94.10`
- Admin Port: `2122`
- Email Domain: `@globalvision.com.au`
- Default Group: `global-vision-media`

### Technical Details
- Built with Go 1.22
- Uses Cobra CLI framework
- GitLab API integration via go-gitlab library
- Multi-stage Docker builds for minimal image size
- Organized SSH configuration structure
- Git includeIf for automatic identity switching