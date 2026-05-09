package config

import "time"

const (
	ServiceName        = "sb-sync-singbox"
	ServiceDisplayName = "Sing-Box (sb-sync managed)"
	ServiceDescription = "Sing-Box proxy service managed by sb-sync"

	RestartDelay         = 10 * time.Second
	StopDelay            = 5 * time.Second
	DefaultSyncInterval  = 60 * time.Minute
	MinimumSyncInterval  = 1 * time.Minute
	VersionCheckInterval = 1 * time.Hour

	GitHubProxyDefault = "https://gh-proxy.com/"
	SingBoxRepoOwner   = "SagerNet"
	SingBoxRepoName    = "sing-box"
)

const (
	AppName        = "sb-sync"
	AppDescription = "A cross-platform sing-box installer and config synchronizer"
)

const (
	LogPrefixInfo  = "[INFO]"
	LogPrefixWarn  = "[WARN]"
	LogPrefixError = "[ERROR]"
	LogPrefixDebug = "[DEBUG]"
)
