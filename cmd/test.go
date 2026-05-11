package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	testURL     string
	testTimeout int
	testProxy   bool
	testDirect  bool
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test network connectivity",
	Long: `Test network connection through the sing-box proxy or directly.

When TUN mode is enabled, traffic is routed through sing-box automatically.
Use --proxy to test through the HTTP proxy instead.`,
	Run: func(cmd *cobra.Command, args []string) {
		runConnectionTest()
	},
}

func init() {
	testCmd.Flags().StringVar(&testURL, "url", "", "Custom URL to test")
	testCmd.Flags().IntVar(&testTimeout, "timeout", 10, "Connection timeout in seconds")
	testCmd.Flags().BoolVar(&testProxy, "proxy", false, "Test through sing-box HTTP proxy (localhost:20122)")
	testCmd.Flags().BoolVar(&testDirect, "direct", true, "Test direct connection (default: true)")
	rootCmd.AddCommand(testCmd)
}

func runConnectionTest() {
	fmt.Println("=== Connection Test ===")
	fmt.Println()

	if testProxy {
		testThroughSingBoxProxy()
		return
	}

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

	results := make([]testResult, 0, len(tests))

	for _, test := range tests {
		result := testDirectConnection(test.url, time.Duration(testTimeout)*time.Second)
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
		fmt.Println("[INFO] Network connectivity is working.")
	} else if successCount > 0 {
		fmt.Printf("[WARN] %d/%d tests passed\n", successCount, len(results))
	} else {
		fmt.Println("[ERROR] All tests failed.")
		fmt.Println()
		fmt.Println("Possible causes:")
		fmt.Println("  1. sing-box is not running - try 'sb-sync service start'")
		fmt.Println("  2. TUN mode is not configured correctly")
		fmt.Println("  3. Firewall is blocking connections")
	}

	fmt.Println()
	fmt.Println("Note: If TUN mode is enabled, traffic is automatically routed through sing-box.")
	fmt.Println("      Use 'sb-sync test --proxy' to test through the HTTP proxy instead.")
}

func testDirectConnection(targetURL string, timeout time.Duration) testResult {
	result := testResult{URL: targetURL}

	start := time.Now()

	tr := &http.Transport{
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	req.Header.Set("User-Agent", "sb-sync/1.0")

	resp, err := client.Do(req)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "timeout") {
			result.Error = "Connection timeout"
		} else if strings.Contains(errStr, "connection refused") {
			result.Error = "Connection refused"
		} else if strings.Contains(errStr, "no such host") {
			result.Error = "DNS resolution failed"
		} else {
			result.Error = errStr
		}
		return result
	}
	defer resp.Body.Close()

	result.Latency = time.Since(start)
	result.StatusCode = resp.StatusCode

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Success = true
	} else if resp.StatusCode >= 400 {
		result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return result
}

func testThroughSingBoxProxy() {
	fmt.Println("[INFO] Testing through sing-box HTTP proxy (localhost:20122)")
	fmt.Println()

	proxyURL := "http://127.0.0.1:20122"

	tests := []struct {
		name string
		url  string
	}{
		{"Google", "https://www.google.com"},
		{"YouTube", "https://www.youtube.com"},
	}

	if testURL != "" {
		tests = append(tests, struct {
			name string
			url  string
		}{"Custom", testURL})
	}

	proxy, _ := url.Parse(proxyURL)

	for _, test := range tests {
		result := testWithProxy(test.url, proxy, time.Duration(testTimeout)*time.Second)
		result.URL = test.name

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

	fmt.Println()
	fmt.Println("Note: Ensure sing-box HTTP proxy is enabled in config.")
}

func testWithProxy(targetURL string, proxy *url.URL, timeout time.Duration) testResult {
	result := testResult{URL: targetURL}

	start := time.Now()

	transport := &http.Transport{
		Proxy:             http.ProxyURL(proxy),
		DisableKeepAlives: true,
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	req.Header.Set("User-Agent", "sb-sync/1.0")

	resp, err := client.Do(req)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.Latency = time.Since(start)
	result.StatusCode = resp.StatusCode

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Success = true
	} else {
		result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}

	return result
}

type testResult struct {
	URL        string
	Success    bool
	Latency    time.Duration
	StatusCode int
	Error      string
}
