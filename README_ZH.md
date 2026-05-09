# sb-sync

`sb-sync` 是一个轻量级的跨平台命令行工具，旨在自动化 [sing-box](https://github.com/SagerNet/sing-box) 的安装、配置同步以及后台服务管理。

[English](./README.md) | **中文说明**

## ✨ 功能特性

- **🚀 一键安装/更新**：自动从 GitHub 下载并安装最新版本的 sing-box 内核。
- **🔄 配置同步**：通过 WebDAV 协议无缝同步你的 sing-box 配置文件。
- **🛡️ 原子更新**：采用临时文件交换机制，确保同步过程中配置文件不会损坏。
- **⚙️ 服务管理**：支持将 sing-box 安装为系统服务，使用 PID 文件管理进程（无需 systemd）。
- **🌐 代理支持**：内置 GitHub 代理支持，加速国内等受限地区的下载速度。
- **⏲️ 定期同步**：后台服务会根据设定的间隔自动检查并更新 WebDAV 上的配置。
- **🔍 健康检查**：内置诊断工具，验证安装和配置状态。
- **🧪 模拟模式**：在执行前预览操作内容。
- **🤖 自动检测**：自动检测现有 sing-box 安装并配置。

## 📥 快速安装

### Linux & macOS (一键脚本)

在终端运行以下命令即可自动下载并安装最新版本：

```bash
curl -fsSL https://raw.githubusercontent.com/kyeo-hub/sb-sync/main/install.sh | bash
```

**如果遇到 GitHub 访问困难，请使用国内加速命令：**

```bash
curl -fsSL https://gh-proxy.com/https://raw.githubusercontent.com/kyeo-hub/sb-sync/main/install.sh | GH_PROXY=https://gh-proxy.com/ bash
```

### Windows

请从 [Releases](https://github.com/kyeo-hub/sb-sync/releases) 页面下载最新的 `.zip` 文件，解压并将 `sb-sync.exe` 放入你的系统变量路径中。

### 源码编译

```bash
go install github.com/kyeo-hub/sb-sync@latest
```

## 🛠️ 使用教程

### 1. 自动检测现有安装（新增）

如果你已经安装了 sing-box，可以使用 auto 命令自动检测并配置：

```bash
# 检测现有安装
sb-sync auto

# 自动导入检测到的配置
sb-sync auto --import

# 显示服务迁移说明
sb-sync auto --migrate
```

### 2. 配置 WebDAV

配置你的 WebDAV 服务器信息以便同步 `config.json`：

```bash
sb-sync config set-dav --url https://your-webdav-server.com --user your_user --pass your_password --path /path/to/config.json
```

### 3. 配置 GitHub 代理 (可选)

如果你在访问 GitHub 时速度较慢，可以设置下载代理：

```bash
sb-sync config set-proxy --url https://ghproxy.com/
```

### 4. 安装 sing-box 内核

下载并安装最新的 sing-box 二进制文件：

```bash
sb-sync install

# 或者使用模拟模式预览安装
sb-sync install --dry-run
```

### 5. 服务管理

将 sing-box 安装为系统后台服务并启动：

```bash
sb-sync service install
sb-sync service start
```

### 6. 状态检查与更新

```bash
# 查看服务运行状态
sb-sync service status

# 重启服务
sb-sync service restart

# 检查并更新 sing-box 到最新版本
sb-sync update

# 同步配置
sb-sync sync

# 显示当前配置信息
sb-sync config show
```

### 7. 健康检查

验证安装和配置状态：

```bash
sb-sync doctor
```

## ⚙️ 配置说明

配置文件默认存储在 `~/.sb-sync/config.yaml`。

| 命令 | 说明 |
|---------|-------------|
| `config set-dav` | 设置 WebDAV 服务器详情 |
| `config set-proxy` | 设置 GitHub 代理地址 |
| `config set-interval` | 设置自动同步频率（单位：分钟） |
| `config show` | 查看当前所有配置 |

## 🤝 参与贡献

欢迎提交 Pull Request 或报告 Issue！

## 📄 开源协议

本项目采用 MIT 协议开源 - 详情请参阅 [LICENSE](LICENSE) 文件。