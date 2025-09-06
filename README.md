# GVM SSH Setup

Automated SSH configuration tool for Global Vision's GitLab instance at `gitlab.globalvision.com.au`.

## Purpose

This tool simplifies the process of setting up SSH access to our internal GitLab server for team members. It automates SSH key generation, configuration, and provides clear instructions for completing the setup.

## Usage

Run the setup script to configure SSH access:

```bash
./setup-gitlab-ssh.sh
```

The script will:
1. Generate a new SSH key pair (if one doesn't exist)
2. Configure SSH settings for gitlab.globalvision.com.au
3. Provide instructions for adding the public key to your GitLab profile

## What it does

- Creates SSH key: `~/.ssh/id_rsa_gitlab_gv`
- Adds configuration to `~/.ssh/config` for gitlab.globalvision.com.au
- Provides step-by-step instructions for completing the setup

## Manual Setup (Alternative)

If you prefer to set up SSH manually:

1. Generate SSH key:
   ```bash
   ssh-keygen -t rsa -b 4096 -C "your.email@globalvision.com.au" -f ~/.ssh/id_rsa_gitlab_gv
   ```

2. Add to SSH config (`~/.ssh/config`):
   ```
   Host gitlab.globalvision.com.au
       HostName gitlab.globalvision.com.au
       User git
       IdentityFile ~/.ssh/id_rsa_gitlab_gv
       IdentitiesOnly yes
   ```

3. Add key to SSH agent:
   ```bash
   ssh-add ~/.ssh/id_rsa_gitlab_gv
   ```

4. Copy public key and add to GitLab profile at https://gitlab.globalvision.com.au/-/profile/keys

## Testing Your Setup

Test your SSH connection:
```bash
ssh -T git@gitlab.globalvision.com.au
```

You should see a welcome message from GitLab.

## Internal Use

This tool is designed specifically for Global Vision team members accessing our internal GitLab instance.