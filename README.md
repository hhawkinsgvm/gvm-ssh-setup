# GVM SSH Setup Tool

A GitLab CE-compatible SSH/Git per-account setup tool for Global Vision Media.

## Features

- Generate SSH keys for different accounts/projects
- Configure SSH aliases and Git identities per directory
- GitLab CE integration with group membership verification
- Automatic SSH key upload to GitLab
- Support for both interactive (wizard) and non-interactive setup
- Docker containerization for consistent environments

## Quick Start

### Interactive Setup (Wizard)

```bash
# Run the wizard
./gvm-ssh wizard
```

### Non-Interactive Setup

```bash
./gvm-ssh setup \
  --account gvm \
  --git-alias gitlab-git \
  --git-host gitlab.globalvision.com.au \
  --git-port 2122 \
  --gitlab-group global-vision-media \
  --upload-key \
  --name "Your Name (GVM)" \
  --email "your.name+gvm@globalvision.com.au"
```

### Docker Usage

```bash
# Interactive wizard
docker run --rm -it \
  -u $(id -u):$(id -g) \
  -e REAL_HOME=/hosthome \
  -e GITLAB_TOKEN="$GITLAB_TOKEN" \
  -v $HOME:/hosthome \
  ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest wizard

# Non-interactive setup
docker run --rm -it \
  -u $(id -u):$(id -g) \
  -e REAL_HOME=/hosthome \
  -e GITLAB_TOKEN="$GITLAB_TOKEN" \
  -v $HOME:/hosthome \
  ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest \
  setup \
    --account gvm \
    --git-alias gitlab-git \
    --git-host gitlab.globalvision.com.au \
    --git-port 2122 \
    --gitlab-group global-vision-media \
    --upload-key
```

## Commands

- `wizard` - Interactive guided setup
- `setup` - Non-interactive setup with flags
- `check` - Show effective SSH config and host fingerprints
- `test` - Test SSH auth and Git connectivity

## GitLab CE Integration

### Authentication

Set your GitLab token:
```bash
export GITLAB_TOKEN="your-gitlab-token"
# OR
glab auth login
```

### Group Membership

The tool verifies you're a member of the specified GitLab group before setup:
```bash
--gitlab-group global-vision-media
```

### SSH Key Upload

Automatically upload your SSH key to GitLab:
```bash
--upload-key
```

## Building

```bash
go mod tidy
go build -o gvm-ssh ./main.go
```

## Docker Build

```bash
docker build -t ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest .
```