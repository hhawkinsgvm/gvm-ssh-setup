#!/bin/bash

# GitLab SSH Setup Script for Global Vision
# Configures SSH access to gitlab.globalvision.com.au

set -e

GITLAB_HOST="gitlab.globalvision.com.au"
SSH_KEY_PATH="$HOME/.ssh/id_rsa_gitlab_gv"
SSH_CONFIG="$HOME/.ssh/config"

echo "=== GitLab SSH Setup for Global Vision ==="
echo "Configuring SSH access to $GITLAB_HOST"
echo

# Check if SSH directory exists
if [ ! -d "$HOME/.ssh" ]; then
    echo "Creating ~/.ssh directory..."
    mkdir -p "$HOME/.ssh"
    chmod 700 "$HOME/.ssh"
fi

# Generate SSH key if it doesn't exist
if [ ! -f "$SSH_KEY_PATH" ]; then
    echo "Generating new SSH key for GitLab..."
    read -p "Enter your email address: " email
    ssh-keygen -t rsa -b 4096 -C "$email" -f "$SSH_KEY_PATH" -N ""
    echo "SSH key generated: $SSH_KEY_PATH"
else
    echo "SSH key already exists: $SSH_KEY_PATH"
fi

# Add SSH config entry for GitLab
if ! grep -q "Host $GITLAB_HOST" "$SSH_CONFIG" 2>/dev/null; then
    echo "Adding SSH configuration for $GITLAB_HOST..."
    cat >> "$SSH_CONFIG" << EOF

# GitLab Global Vision configuration
Host $GITLAB_HOST
    HostName $GITLAB_HOST
    User git
    IdentityFile $SSH_KEY_PATH
    IdentitiesOnly yes
EOF
    chmod 600 "$SSH_CONFIG"
    echo "SSH configuration added to $SSH_CONFIG"
else
    echo "SSH configuration for $GITLAB_HOST already exists"
fi

# Start SSH agent and add key
if [ -z "$SSH_AUTH_SOCK" ]; then
    echo "Starting SSH agent..."
    eval "$(ssh-agent -s)"
fi

ssh-add "$SSH_KEY_PATH" 2>/dev/null || true

echo
echo "=== Setup Complete ==="
echo "Your public key is located at: ${SSH_KEY_PATH}.pub"
echo
echo "Next steps:"
echo "1. Copy your public key to clipboard:"
echo "   cat ${SSH_KEY_PATH}.pub | pbcopy   # macOS"
echo "   cat ${SSH_KEY_PATH}.pub | xclip -sel clip   # Linux"
echo
echo "2. Add the public key to your GitLab profile:"
echo "   - Go to https://$GITLAB_HOST/-/profile/keys"
echo "   - Paste your public key and give it a descriptive title"
echo
echo "3. Test your connection:"
echo "   ssh -T git@$GITLAB_HOST"
echo
echo "You can now clone repositories using:"
echo "   git clone git@$GITLAB_HOST:group/repository.git"