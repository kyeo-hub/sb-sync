# sb-sync

`sb-sync` is a lightweight, cross-platform CLI tool designed to automate the installation, configuration synchronization, and background management of [sing-box](https://github.com/SagerNet/sing-box).

**English** | [中文说明](./README_ZH.md)

## ✨ Features

- **🚀 One-Click Install/Update**: Automatically download and install the latest version of sing-box from GitHub.
- **🔄 Configuration Sync**: Seamlessly sync your sing-box configuration from any WebDAV server.
- **🛡️ Atomic Updates**: Prevents configuration corruption during synchronization.
- **⚙️ Service Management**: Install sing-box as a system service (Windows, Linux, macOS) with auto-restart on crash.
- **🌐 Proxy Support**: Built-in support for GitHub proxies to speed up downloads in restricted regions.
- **⏲️ Periodic Sync**: Background service automatically checks for configuration updates.
- **🔍 Health Check**: Built-in diagnostic tool to verify installation and configuration status.
- **🧪 Dry-Run Mode**: Preview actions before executing them.
- **💻 Shell Completion**: Support for Bash, Zsh, Fish, and PowerShell auto-completion.

## 📥 Installation

### Linux & macOS (Quick Install)

Run the following command in your terminal to automatically download and install the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/kyeo-hub/sb-sync/main/install.sh | bash
```

**If you have trouble accessing GitHub, use the proxy command:**

```bash
curl -fsSL https://gh-proxy.com/https://raw.githubusercontent.com/kyeo-hub/sb-sync/main/install.sh | GH_PROXY=https://gh-proxy.com/ bash
```

### Windows

Download the latest `.zip` file from the [Releases](https://github.com/kyeo-hub/sb-sync/releases) page, extract it, and place `sb-sync.exe` in your system's PATH.

### Build from Source

```bash
go install github.com/kyeo-hub/sb-sync@latest
```

## 🛠️ Usage

### 1. Configure WebDAV
First, set up your WebDAV credentials to sync your `config.json`:

```bash
sb-sync config set-dav --url https://your-webdav-server.com --user your_user --pass your_password --path /path/to/config.json
```

### 2. Configure GitHub Proxy (Optional)
If you are in a region with slow access to GitHub, set a proxy:

```bash
sb-sync config set-proxy --url https://ghproxy.com/
```

### 3. Install sing-box
Install the latest core:

```bash
sb-sync install
```

### 4. Manage Service
Install and start the background service:

```bash
sb-sync service install
sb-sync service start
```

### 5. Check Status & Update
```bash
# Check service status
sb-sync service status

# Restart the service
sb-sync service restart

# Update sing-box to latest version
sb-sync update

# Show current configuration
sb-sync config show
```

### 6. Health Check
Verify your installation and configuration:

```bash
# Run health check
sb-sync doctor

# Output in JSON format
sb-sync doctor --json
```

### 7. Dry-Run Mode
Preview actions before executing them:

```bash
# Preview installation
sb-sync install --dry-run

# Preview configuration sync
sb-sync sync --dry-run

# Preview update
sb-sync update --dry-run
```

### 8. Shell Completion
Enable auto-completion for your shell:

```bash
# Bash
sb-sync completion bash > /etc/bash_completion.d/sb-sync

# Zsh
sb-sync completion zsh > "${fpath[1]}/_sb-sync"

# Fish
sb-sync completion fish > ~/.config/fish/completions/sb-sync.fish

# PowerShell
sb-sync completion powershell > sb-sync.ps1
```

## ⚙️ Configuration Reference

The configuration file is stored at `~/.sb-sync/config.yaml`.

| Command | Description |
|---------|-------------|
| `config set-dav` | Set WebDAV server details |
| `config set-proxy` | Set GitHub proxy URL |
| `config set-interval` | Set sync frequency (minutes) |
| `config show` | View current settings |

### Environment Variables

You can also configure using environment variables (prefix: `SB_SYNC_`):

| Variable | Description |
|----------|-------------|
| `SB_SYNC_GITHUB_PROXY` | GitHub proxy URL |
| `SB_SYNC_WEBDAV_URL` | WebDAV server URL |
| `SB_SYNC_WEBDAV_USER` | WebDAV username |
| `SB_SYNC_WEBDAV_PASS` | WebDAV password |
| `SB_SYNC_WEBDAV_PATH` | Remote config file path |
| `SB_SYNC_INSTALL_DIR` | Installation directory |
| `SB_SYNC_SYNC_INTERVAL` | Sync interval (minutes) |

## 🧰 Build from Source

```bash
# Clone the repository
git clone https://github.com/kyeo-hub/sb-sync.git
cd sb-sync

# Build
make build

# Or use Go directly
go build -o sb-sync .

# Run tests
make test

# Run linting
make lint
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
