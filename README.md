# sb-sync

`sb-sync` is a lightweight, cross-platform CLI tool designed to automate the installation, configuration synchronization, and background management of [sing-box](https://github.com/SagerNet/sing-box).

**English** | [中文说明](./README_ZH.md)

## ⚠️ Important Note

**If you already have sing-box installed manually**, use `sb-sync auto --all` to automatically:

- Detect existing installation paths
- Stop running sing-box processes
- Import existing configuration files
- Disable original systemd services

This prevents conflicts from multiple sing-box instances running simultaneously!

## ✨ Features

- **🚀 One-Click Install/Update**: Automatically download and install the latest version of sing-box from GitHub.
- **🔄 Configuration Sync**: Seamlessly sync your sing-box configuration from any WebDAV server.
- **🛡️ Atomic Updates**: Prevents configuration corruption during synchronization.
- **⚙️ Service Management**: Manage sing-box as a background service using PID file-based process management (no systemd required).
- **🌐 Proxy Support**: Built-in support for GitHub proxies to speed up downloads in restricted regions.
- **⏲️ Periodic Sync**: Background service automatically checks for configuration updates.
- **🔍 Health Check**: Built-in diagnostic tool to verify installation and configuration status.
- **🧪 Dry-Run Mode**: Preview actions before executing them.
- **💻 Shell Completion**: Support for Bash, Zsh, Fish, and PowerShell auto-completion.
- **🤖 Auto-Detection**: Automatically detect existing sing-box installations, stop old processes, avoid conflicts.

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

### Scenario 1: Fresh Install (No sing-box on system)

```bash
# 1. Install sing-box
sb-sync install

# 2. Configure WebDAV server
sb-sync config set-dav --url https://your-webdav.com --user your_user --pass your_password --path /config.json

# 3. Configure GitHub proxy (optional, recommended for China)
sb-sync config set-proxy --url https://gh-proxy.com/

# 4. Start service
sb-sync service install
sb-sync service start

# 5. Check status
sb-sync doctor
```

### Scenario 2: Migrate from Existing Installation (Recommended)

If you already have sing-box installed, use one-command migration:

```bash
# Run all migration steps at once
sb-sync auto --all
```

This will automatically:
1. Detect existing sing-box installation
2. Stop running processes
3. Import configuration files
4. Disable original systemd services

Then start the new service:
```bash
sb-sync service start
```

### Scenario 3: Step-by-Step Migration

If you prefer granular control:

```bash
# 1. See what is detected
sb-sync auto

# 2. Stop old processes only
sb-sync auto --kill

# 3. Import config only
sb-sync auto --import

# 4. Disable systemd service only
sb-sync auto --migrate
```

## 💡 Daily Usage

```bash
# Check service status
sb-sync service status

# Sync latest configuration
sb-sync sync

# Test network connection
sb-sync test

# Health check
sb-sync doctor

# Restart service
sb-sync service restart

# Update sing-box
sb-sync update
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

| Variable | Description |
|----------|-------------|
| `SB_SYNC_WEBDAV_URL` | WebDAV server URL |
| `SB_SYNC_WEBDAV_USER` | WebDAV username |
| `SB_SYNC_WEBDAV_PASS` | WebDAV password |
| `SB_SYNC_WEBDAV_PATH` | WebDAV config file path |
| `SB_SYNC_GITHUB_PROXY` | GitHub proxy URL |
| `GH_PROXY` / `GITHUB_PROXY` | Proxy URL (compatible) |

### Default Paths

| Path | Description |
|------|-------------|
| `~/.sb-sync/` | Installation directory |
| `~/.sb-sync/config.yaml` | sb-sync config file |
| `~/.sb-sync/bin/sing-box` | sing-box binary |
| `~/.sb-sync/config.json` | sing-box config |
| `~/.sb-sync/sing-box.log` | Log file |
| `~/.sb-sync/sing-box.pid` | PID file |

## 🗑️ Uninstall

```bash
# Uninstall sb-sync (keep config)
sb-sync uninstall --keep-config

# Complete uninstall (including config)
sb-sync uninstall
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
