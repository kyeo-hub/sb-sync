# sb-sync (Sing-Box Sync)

一个跨平台的 sing-box 一键安装与配置同步工具，支持 Windows、macOS 和 Linux。

## 主要流程

1.  **安装**: 下载并解压 `sing-box` 核心。
2.  **配置**: 设置 WebDAV 凭据（支持 GitHub 代理）。
3.  **同步**: 从 WebDAV 下载 `config.json`。
4.  **服务**: 将 `sing-box` 注册为系统服务并后台运行。

## 使用说明

### 1. 安装与初始化
首先编译或下载 `sb-sync` 二进制文件。

```powershell
# 安装 sing-box 核心
.\sb-sync.exe install
```

### 2. 配置 WebDAV
以坚果云为例：

```powershell
.\sb-sync.exe config set-dav --url "https://dav.jianguoyun.com/dav/" --user "你的邮箱" --pass "你的应用密码" --path "/path/to/your/config.json"
```

*注意：如果下载 sing-box 缓慢，可以设置 GitHub 代理（默认为 https://ghproxy.com/）：*
```powershell
.\sb-sync.exe config set-proxy --url "https://mirror.ghproxy.com/"
```

### 3. 同步配置
```powershell
.\sb-sync.exe sync
```

### 4. 服务管理
```powershell
# 安装为系统服务
.\sb-sync.exe service install

# 启动服务
.\sb-sync.exe service start

# 查看状态
.\sb-sync.exe service status
```

## 项目结构
- `main.go`: 程序入口
- `cmd/`: CLI 命令定义 (Cobra)
- `pkg/config/`: 配置读写逻辑 (Viper)
- `pkg/downloader/`: sing-box 下载与解压逻辑
- `pkg/sync/`: WebDAV 同步逻辑
- `pkg/service/`: 跨平台服务管理逻辑 (kardianos/service)
