# GVM SSH Setup Examples

## Standard GVM Account Setup

```bash
# Interactive wizard with GVM defaults
./gvm-ssh wizard

# Non-interactive setup for GVM account
./gvm-ssh setup \
  --account gvm \
  --git-alias gvm-git \
  --git-host gitlab.globalvision.com.au \
  --git-port 2122 \
  --gitlab-group global-vision-media \
  --admin-alias gvm-host \
  --admin-host 203.32.94.10 \
  --admin-port 2122 \
  --upload-key \
  --name "Your Name (GVM)" \
  --email "your.name+gvm@globalvision.com.au"
```

## Multi-Account Setup

```bash
# Setup for different client accounts
./gvm-ssh setup \
  --account client1 \
  --git-alias client1-git \
  --git-host gitlab.globalvision.com.au \
  --git-port 2122 \
  --gitlab-group global-vision-media \
  --upload-key \
  --name "Your Name (Client1)" \
  --email "your.name+client1@globalvision.com.au"

./gvm-ssh setup \
  --account client2 \
  --git-alias client2-git \
  --git-host gitlab.globalvision.com.au \
  --git-port 2122 \
  --gitlab-group global-vision-media \
  --upload-key \
  --name "Your Name (Client2)" \
  --email "your.name+client2@globalvision.com.au"
```

## Deploy Keys for CI/Servers

```bash
# Generate deploy key for a project (admin only)
./gvm-ssh deploy-key \
  --project myproject \
  --gitlab-group global-vision-media \
  --title "CI Deploy Key"

# Read-only deploy key
./gvm-ssh deploy-key \
  --project myproject \
  --read-only \
  --gitlab-group global-vision-media
```

## Testing and Validation

```bash
# Check SSH configuration
./gvm-ssh check --alias gvm-git --git-host gitlab.globalvision.com.au --git-port 2122

# Test SSH authentication
./gvm-ssh test --alias gvm-git

# Test Git repository access
./gvm-ssh test --alias gvm-git --repo Global-Vision-Media/some-repo
```

## Docker Usage

```bash
# Set GitLab token first
export GITLAB_TOKEN="your-gitlab-token"

# Interactive setup via Docker
docker run --rm -it \
  -u $(id -u):$(id -g) \
  -e REAL_HOME=/hosthome \
  -e GITLAB_TOKEN="$GITLAB_TOKEN" \
  -v $HOME:/hosthome \
  ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest wizard

# Non-interactive setup via Docker
docker run --rm -it \
  -u $(id -u):$(id -g) \
  -e REAL_HOME=/hosthome \
  -e GITLAB_TOKEN="$GITLAB_TOKEN" \
  -v $HOME:/hosthome \
  ghcr.io/hhawkinsgvm/gvm-ssh-setup:latest \
  setup \
    --account gvm \
    --git-alias gvm-git \
    --git-host gitlab.globalvision.com.au \
    --git-port 2122 \
    --gitlab-group global-vision-media \
    --upload-key
```

## Environment Variables

- `GITLAB_TOKEN` - GitLab API token for authentication and key upload
- `REAL_HOME` - Used in Docker to specify host home directory