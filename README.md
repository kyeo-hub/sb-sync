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

## 📥 Installation

### Linux & macOS (Quick Install)

Run the following command in your terminal to automatically download and install the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/kyeo-hub/sb-sync/main/install.sh | bash
```

**If you have trouble accessing GitHub, use the proxy command:**

```bash
curl -fsSL https://ghproxy.com/https://raw.githubusercontent.com/kyeo-hub/sb-sync/main/install.sh | GH_PROXY=https://ghproxy.com/ bash
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

# Update sing-box to latest version
sb-sync update

# Show current configuration
sb-sync config show
```

## ⚙️ Configuration Reference

The configuration file is stored at `~/.sb-sync/config.yaml`.

| Command | Description |
|---------|-------------|
| `config set-dav` | Set WebDAV server details |
| `config set-proxy` | Set GitHub proxy URL |
| `config set-interval` | Set sync frequency (minutes) |
| `config show` | View current settings |

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
