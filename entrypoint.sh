#!/usr/bin/env bash
set -e

# GVM SSH Setup Tool Entrypoint Script
# Handles Docker environment setup and user permissions

# Function to display banner
show_banner() {
    cat << 'EOF'
   _____ _    ____  __   _____ _____ _    _            
  / ____| |  / /  \/  | / ____/ ____| |  | |           
 | |  __| | / /| \  / || (___| (___ | |__| |           
 | | |_ | |/ / | |\/| | \___ \\___ \|  __  |           
 | |__| |   <  | |  | | ____) |___) | |  | |           
  \_____|_|\_\ |_|  |_||_____/_____/|_|  |_|           
                                                       
           SSH & Git Setup Tool for GitLab CE
              Global Vision Media - v1.0.0
EOF
    echo
}

# Function to check and setup environment
setup_environment() {
    # Set HOME to REAL_HOME if running in Docker with host volume mount
    if [ -n "$REAL_HOME" ]; then
        export HOME="$REAL_HOME"
        echo "🏠 Using host home directory: $HOME"
    fi

    # Ensure SSH directory exists and has correct permissions
    if [ -n "$HOME" ] && [ -d "$HOME" ]; then
        mkdir -p "$HOME/.ssh"
        chmod 700 "$HOME/.ssh" 2>/dev/null || true
        
        # Check if we can write to the home directory
        if [ ! -w "$HOME" ]; then
            echo "⚠️  Warning: Cannot write to home directory $HOME"
            echo "   Make sure you're running with correct user permissions:"
            echo "   docker run -u \$(id -u):\$(id -g) ..."
        fi
    fi

    # Check for GitLab authentication
    if [ -n "$GITLAB_TOKEN" ]; then
        echo "🔐 GitLab token found in environment"
    elif command -v glab >/dev/null 2>&1; then
        echo "🛠️  glab CLI available"
        echo "   Run 'glab auth login' to authenticate with GitLab"
    else
        echo "ℹ️  No GitLab authentication configured"
        echo "   Set GITLAB_TOKEN or install glab CLI for key upload features"
    fi
}

# Main execution
main() {
    # Show banner for interactive sessions
    if [ -t 0 ] && [ "$1" != "--help" ] && [ "$1" != "-h" ]; then
        show_banner
    fi
    
    # Setup environment
    setup_environment
    
    # Execute the main application
    exec /usr/local/bin/gvm-ssh "$@"
}

# Run main function with all arguments
main "$@"