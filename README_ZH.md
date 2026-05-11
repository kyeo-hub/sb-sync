# sb-sync

`sb-sync` 是一个轻量级的跨平台命令行工具，旨在自动化 [sing-box](https://github.com/SagerNet/sing-box) 的安装、配置同步以及后台服务管理。

[English](./README.md) | **中文说明**

## ⚠️ 重要说明

**如果你已经手动安装了 sing-box**，使用 `sb-sync auto --all` 可以一键完成：
- 自动检测现有安装位置
- 停止正在运行的 sing-box 进程
- 导入现有配置文件
- 禁用原有的 systemd 服务

这样可以避免多个 sing-box 实例同时运行造成的冲突！

## ✨ 功能特性

- **🚀 一键安装/更新**：自动从 GitHub 下载并安装最新版本的 sing-box 内核。
- **🔄 配置同步**：通过 WebDAV 协议无缝同步你的 sing-box 配置文件。
- **🛡️ 原子更新**：采用临时文件交换机制，确保同步过程中配置文件不会损坏。
- **⚙️ 服务管理**：支持将 sing-box 安装为系统服务，使用 PID 文件管理进程（无需 systemd）。
- **🌐 代理支持**：内置 GitHub 代理支持，加速国内等受限地区的下载速度。
- **⏲️ 定期同步**：后台服务会根据设定的间隔自动检查并更新 WebDAV 上的配置。
- **🔍 健康检查**：内置诊断工具，验证安装和配置状态。
- **🧪 模拟模式**：在执行前预览操作内容。
- **🤖 自动检测**：自动检测现有 sing-box 安装，自动停止旧进程，避免冲突。

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

### 场景一：首次安装（机器上没有 sing-box）

```bash
# 1. 安装 sing-box 内核
sb-sync install

# 2. 配置 WebDAV 服务器
sb-sync config set-dav --url https://your-webdav.com --user your_user --pass your_password --path /config.json

# 3. 配置 GitHub 代理（国内用户可选）
sb-sync config set-proxy --url https://gh-proxy.com/

# 4. 启动服务
sb-sync service install
sb-sync service start

# 5. 检查状态
sb-sync doctor
```

### 场景二：从现有安装迁移（推荐）

如果你的机器上已经安装了 sing-box，直接使用一键迁移：

```bash
# 一键完成所有迁移步骤
sb-sync auto --all
```

这会自动：
1. 检测现有 sing-box 安装位置
2. 停止正在运行的旧进程
3. 导入配置文件
4. 禁用原有的 systemd 服务

然后启动新服务：
```bash
sb-sync service start
```

### 场景三：分步迁移

如果你想逐步控制，可以使用以下命令：

```bash
# 1. 先看看检测到什么
sb-sync auto

# 2. 只停止旧进程
sb-sync auto --kill

# 3. 只导入配置
sb-sync auto --import

# 4. 只禁用 systemd 服务
sb-sync auto --migrate
```

## 💡 日常使用

```bash
# 查看服务状态
sb-sync service status

# 同步最新配置
sb-sync sync

# 测试网络连接
sb-sync test

# 健康检查
sb-sync doctor

# 重启服务
sb-sync service restart

# 更新 sing-box
sb-sync update
```

## ⚙️ 配置说明

配置文件默认存储在 `~/.sb-sync/config.yaml`。

| 命令 | 说明 |
|---------|-------------|
| `config set-dav` | 设置 WebDAV 服务器详情 |
| `config set-proxy` | 设置 GitHub 代理地址 |
| `config set-interval` | 设置自动同步频率（单位：分钟） |
| `config show` | 查看当前所有配置 |

### 环境变量

| 变量 | 说明 |
|------|------|
| `SB_SYNC_WEBDAV_URL` | WebDAV 地址 |
| `SB_SYNC_WEBDAV_USER` | WebDAV 用户名 |
| `SB_SYNC_WEBDAV_PASS` | WebDAV 密码 |
| `SB_SYNC_WEBDAV_PATH` | WebDAV 配置文件路径 |
| `SB_SYNC_GITHUB_PROXY` | GitHub 代理地址 |
| `GH_PROXY` / `GITHUB_PROXY` | 代理地址（兼容） |

### 默认路径

| 路径 | 说明 |
|------|------|
| `~/.sb-sync/` | 安装目录 |
| `~/.sb-sync/config.yaml` | sb-sync 配置文件 |
| `~/.sb-sync/bin/sing-box` | sing-box 二进制文件 |
| `~/.sb-sync/config.json` | sing-box 配置 |
| `~/.sb-sync/sing-box.log` | 日志文件 |
| `~/.sb-sync/sing-box.pid` | PID 文件 |

## 🗑️ 卸载

```bash
# 卸载 sb-sync（保留配置）
sb-sync uninstall --keep-config

# 完全卸载（包括配置）
sb-sync uninstall
```

## 🤝 参与贡献

欢迎提交 Pull Request 或报告 Issue！

## 📄 开源协议

本项目采用 MIT 协议开源 - 详情请参阅 [LICENSE](LICENSE) 文件。
