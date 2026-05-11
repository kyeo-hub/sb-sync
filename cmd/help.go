package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Show detailed help information",
	Run: func(cmd *cobra.Command, args []string) {
		printHelp()
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)
}

func printHelp() {
	fmt.Print(`╔══════════════════════════════════════════════════════════════════╗
║                        sb-sync 使用帮助                          ║
╚══════════════════════════════════════════════════════════════════╝

📦 安装与更新
  sb-sync install           安装或更新 sing-box 内核
  sb-sync update            更新 sing-box 到最新版本
  sb-sync uninstall         卸载 sb-sync

🔄 自动迁移
  sb-sync auto              检测现有 sing-box 安装
  sb-sync auto --all       一键迁移（停止进程+导入配置+禁用服务）
  sb-sync auto --kill      停止现有 sing-box 进程
  sb-sync auto --import    导入检测到的配置
  sb-sync auto --migrate   停止并禁用 systemd 服务

🔧 配置管理
  sb-sync config set-dav --url <url> --user <user> --pass <pass> --path <path>
  sb-sync config set-proxy --url <url>   设置 GitHub 代理
  sb-sync config set-interval <minutes>  设置自动同步间隔
  sb-sync config show                   显示当前配置

🔄 配置同步
  sb-sync sync             从 WebDAV 下载并同步配置
  sb-sync sync --dry-run   预览同步操作（不实际执行）

🌐 网络测试
  sb-sync test             测试网络连接（Google/YouTube/GitHub）
  sb-sync test --url <url> 测试自定义 URL
  sb-sync test --timeout <sec>    设置超时时间（默认 10 秒）
  sb-sync test --proxy   通过 sing-box HTTP 代理测试

💻 服务管理
  sb-sync service install  安装服务
  sb-sync service start    启动服务
  sb-sync service stop     停止服务
  sb-sync service restart  重启服务
  sb-sync service status   查看服务状态

🔍 健康检查
  sb-sync doctor           检查安装和配置状态
  sb-sync doctor --json    JSON 格式输出

📋 常用命令流程
  1. 首次安装:
     sb-sync install
     sb-sync config set-dav --url <url> --user <user> --pass <pass> --path <path>
     sb-sync service install
     sb-sync service start

  2. 从现有安装迁移:
     sb-sync auto --all
     sb-sync service start

  3. 日常使用:
     sb-sync service status    查看状态
     sb-sync sync              同步配置
     sb-sync test              测试连接
     sb-sync doctor            健康检查

📚 环境变量
  GH_PROXY / GITHUB_PROXY / SB_SYNC_GITHUB_PROXY   GitHub 代理
  SB_SYNC_WEBDAV_URL                               WebDAV 地址
  SB_SYNC_WEBDAV_USER                              WebDAV 用户名
  SB_SYNC_WEBDAV_PASS                              WebDAV 密码
  SB_SYNC_WEBDAV_PATH                              WebDAV 配置文件路径

📂 默认路径
  安装目录: ~/.sb-sync/
  配置文件: ~/.sb-sync/config.yaml
  sing-box: ~/.sb-sync/bin/sing-box
  日志文件: ~/.sb-sync/sing-box.log
  PID 文件: ~/.sb-sync/sing-box.pid

🔗 相关链接
  项目地址: https://github.com/kyeo-hub/sb-sync
  sing-box: https://github.com/SagerNet/sing-box
`)
}
