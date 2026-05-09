package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
)

var (
	testURL      string
	testTimeout  int
	testUseProxy bool
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test network connection through proxy",
	Run: func(cmd *cobra.Command, args []string) {
		runConnectionTest()
	},
}

func init() {
	testCmd.Flags().StringVar(&testURL, "url", "", "Custom URL to test")
	testCmd.Flags().IntVar(&testTimeout, "timeout", 10, "Connection timeout in seconds")
	testCmd.Flags().BoolVar(&testUseProxy, "proxy", true, "Use proxy settings from config")
	rootCmd.AddCommand(testCmd)
}

type TestResult struct {
	URL       string
	Success   bool
	Latency   time.Duration
	Error     string
	IPAddress string
}

func runConnectionTest() {
	fmt.Println("=== Connection Test ===")
	fmt.Println()

	tests := []struct {
		name string
		url  string
	}{
		{"Google", "https://www.google.com"},
		{"YouTube", "https://www.youtube.com"},
		{"GitHub", "https://github.com"},
	}

	if testURL != "" {
		tests = append(tests, struct {
			name string
			url  string
		}{"Custom", testURL})
	}

	proxyURL := getProxyURL()
	results := make([]TestResult, 0, len(tests))

	for _, test := range tests {
		result := testConnection(test.url, proxyURL, time.Duration(testTimeout)*time.Second)
		result.URL = test.name
		results = append(results, result)

		status := "✅ SUCCESS"
		if !result.Success {
			status = "❌ FAILED"
		}

		latency := ""
		if result.Success && result.Latency > 0 {
			latency = fmt.Sprintf(" (%.0fms)", float64(result.Latency)/float64(time.Millisecond))
		}

		fmt.Printf("[%s] %s%s", status, test.name, latency)
		if result.Error != "" {
			fmt.Printf(" - %s", result.Error)
		}
		fmt.Println()
	}

	successCount := 0
	for _, r := range results {
		if r.Success {
			successCount++
		}
	}

	fmt.Println()
	if successCount == len(results) {
		fmt.Println("[INFO] All tests passed!")
	} else if successCount > 0 {
		fmt.Printf("[WARN] %d/%d tests passed\n", successCount, len(results))
	} else {
		fmt.Println("[ERROR] All tests failed. Check your proxy settings.")
	}

	if proxyURL != "" {
		fmt.Printf("[INFO] Using proxy: %s\n", proxyURL)
	} else {
		fmt.Println("[INFO] Testing without proxy")
	}
}

func testConnection(targetURL string, proxyURL string, timeout time.Duration) TestResult {
	result := TestResult{URL: targetURL}

	start := time.Now()

	client := &http.Client{
		Timeout: timeout,
	}

	if proxyURL != "" && testUseProxy {
		proxy, err := url.Parse(proxyURL)
		if err == nil {
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
			client.Transport = transport
		}
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	req.Header.Set("User-Agent", "sb-sync/1.0")

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			result.Error = "Connection timeout"
		} else if strings.Contains(err.Error(), "connection refused") {
			result.Error = "Connection refused"
		} else {
			result.Error = err.Error()
		}
		return result
	}
	defer resp.Body.Close()

	result.Latency = time.Since(start)
	result.Success = resp.StatusCode < 500

	if resp.StatusCode >= 500 {
		result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return result
}

func getProxyURL() string {
	if config.AppConfig.GithubProxy != "" {
		return config.AppConfig.GithubProxy
	}
	return ""
}
