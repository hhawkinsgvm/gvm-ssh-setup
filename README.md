# GVM SSH Setup Tool

A comprehensive SSH and Git configuration tool optimized for Global Vision Media's GitLab CE environment.

## Features

- 🔑 **SSH Key Management**: Generate and manage SSH keys per account
- 🔧 **SSH Alias Configuration**: Create organized SSH aliases for different environments
- 📁 **Per-Directory Git Identity**: Automatic Git user switching based on project folder
- 🦊 **GitLab CE Integration**: Upload SSH keys and verify group membership
- 🐳 **Docker Support**: Containerized execution with host filesystem access
- 🧙 **Interactive Wizard**: Guided setup with sensible GVM defaults
- 🚀 **Deploy Key Support**: Special mode for CI/server deployments

## Quick Start

### Using Docker (Recommended)

```bash
# Interactive wizard
docker run --rm -it \
  -u $(id -u):$(id -g) \
  -e REAL_HOME=/hosthome \
  -v $HOME:/hosthome \
  -e GITLAB_TOKEN=$GITLAB_TOKEN \
  ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest wizard

# Non-interactive setup
docker run --rm -it \
  -u $(id -u):$(id -g) \
  -e REAL_HOME=/hosthome \
  -v $HOME:/hosthome \
  -e GITLAB_TOKEN=$GITLAB_TOKEN \
  ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest \
  setup --account gvm --git-alias gitlab-git --upload-key
```

### Native Installation

```bash
# Build from source
git clone https://github.com/hhawkinsgvm/gvm-ssh-setup.git
cd gvm-ssh-setup
go build -o gvm-ssh ./main.go

# Run wizard
./gvm-ssh wizard
```

## Commands

### `wizard` - Interactive Setup

The wizard provides a guided setup experience with GVM-specific defaults:

```bash
gvm-ssh wizard
```

Pre-configured defaults:
- Git Host: `gitlab.globalvision.com.au`
- Git Port: `2122`
- Admin Host: `203.32.94.10`
- Admin Port: `2122`
- Email Domain: `@globalvision.com.au`

### `setup` - Non-Interactive Setup

```bash
# Basic developer setup
gvm-ssh setup \
  --account myaccount \
  --git-alias myaccount-git \
  --upload-key

# Full setup with admin access
gvm-ssh setup \
  --account gvm \
  --git-alias gitlab-git \
  --admin-alias gitlab-host \
  --name "Your Name" \
  --email "your.email@globalvision.com.au" \
  --upload-key

# Deploy key for CI/servers
gvm-ssh setup \
  --account ci \
  --git-alias ci-git \
  --deploy-key \
  --skip-auth
```

### `check` - Configuration Verification

```bash
# Check SSH alias configuration
gvm-ssh check --alias gitlab-git

# Check with host key verification
gvm-ssh check --alias gitlab-git --git-host gitlab.globalvision.com.au
```

### `test` - Connectivity Testing

```bash
# Test SSH authentication
gvm-ssh test --alias gitlab-git

# Test Git connectivity
gvm-ssh test --alias gitlab-git --repo Global-Vision-Media/my-project
```

## Configuration

### Environment Variables

- `GITLAB_TOKEN`: GitLab personal access token for API operations
- `GITLAB_BASE_URL`: GitLab instance URL (defaults to configured Git host)
- `REAL_HOME`: Host home directory when running in Docker

### GitLab Token Setup

Create a personal access token in GitLab with the following scopes:
- `read_user`: For user verification
- `read_api`: For group membership checking
- `write_repository`: For SSH key upload (if using `--upload-key`)

```bash
# Set token in environment
export GITLAB_TOKEN="your-token-here"

# Or use glab CLI
glab auth login
```

## How It Works

### SSH Configuration

The tool creates an organized SSH configuration structure:

```
~/.ssh/
├── config              # Main config with Include directive
├── config.d/           # Organized alias configurations
│   ├── gitlab-git.conf # Git operations alias
│   └── gitlab-host.conf # Admin/shell access alias
└── id_ed25519_*        # Per-account SSH keys
```

### Git Configuration

Per-directory Git identity using `includeIf`:

```
~/
├── .gitconfig                    # Global config
├── .gitconfig-gvm               # Account-specific config
└── projects/
    └── gvm/                     # Project folder
        └── my-project/          # Git identity auto-switches here
```

## GitLab CE Specific Features

### Group Membership Verification

The tool verifies you're a member of the specified GitLab group before allowing setup:

```bash
gvm-ssh setup --account myaccount --org-group global-vision-media
```

### SSH Key Upload

Automatically upload generated SSH keys to your GitLab account:

```bash
gvm-ssh setup --account myaccount --upload-key
```

### Deploy Keys

Special mode for CI/server environments that skips personal Git configuration:

```bash
gvm-ssh setup \
  --account ci \
  --git-alias ci-git \
  --deploy-key \
  --skip-auth
```

## Docker Deployment

### Building the Image

```bash
# Build for production
docker build -t ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest .

# Push to GitLab Container Registry
docker push ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest
```

### Usage Patterns

```bash
# Development workflow
docker run --rm -it \
  -u $(id -u):$(id -g) \
  -e REAL_HOME=/hosthome \
  -v $HOME:/hosthome \
  -e GITLAB_TOKEN=$GITLAB_TOKEN \
  ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest wizard

# CI/Server deployment
docker run --rm \
  -v /etc/ssh-keys:/keys \
  -e GITLAB_TOKEN=$DEPLOY_TOKEN \
  ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest \
  setup --account ci --deploy-key --key /keys/ci_key
```

## Security Considerations

- SSH keys are generated with Ed25519 algorithm and 100 rounds
- All SSH configurations use `IdentitiesOnly yes` for security
- Deploy keys are isolated from personal accounts
- Group membership verification prevents unauthorized access
- Private registry distribution limits tool access

## Troubleshooting

### Common Issues

1. **Permission Denied**: Ensure correct user permissions when running Docker
2. **GitLab Authentication Failed**: Verify `GITLAB_TOKEN` and network access
3. **SSH Agent Issues**: The tool automatically handles ssh-agent setup
4. **Host Key Verification**: Use `gvm-ssh check` to verify host fingerprints

### Debug Mode

```bash
# Enable verbose SSH output
SSH_VERBOSE=1 gvm-ssh test --alias gitlab-git

# Check effective Git configuration
cd ~/projects/myaccount
git config --list | grep user
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

MIT License - see LICENSE file for details.