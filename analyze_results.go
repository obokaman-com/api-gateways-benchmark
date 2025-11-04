package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Metrics represents the parsed metrics from a hey output file
type Metrics struct {
	TotalRequests        int                       `json:"total_requests"`
	TotalTime            float64                   `json:"total_time"`
	RequestsPerSec       float64                   `json:"requests_per_sec"`
	AvgResponseTime      float64                   `json:"avg_response_time"`
	Slowest              float64                   `json:"slowest"`
	Fastest              float64                   `json:"fastest"`
	LatencyDistribution  map[string]float64        `json:"latency_distribution"`
	StatusCodes          map[string]int            `json:"status_codes"`
	Errors               map[string]int            `json:"errors"`
}

// BenchmarkData holds all benchmark results organized by gateway, endpoint, and concurrency
type BenchmarkData map[string]map[string]map[int]*Metrics

// parseHeyOutput parses a hey output file and extracts metrics
func parseHeyOutput(filepath string) (*Metrics, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	text := string(content)
	metrics := &Metrics{
		LatencyDistribution: make(map[string]float64),
		StatusCodes:         make(map[string]int),
		Errors:              make(map[string]int),
	}

	// Extract total requests and time
	re := regexp.MustCompile(`(\d+) requests in ([\d.]+)s`)
	if match := re.FindStringSubmatch(text); match != nil {
		metrics.TotalRequests, _ = strconv.Atoi(match[1])
		metrics.TotalTime, _ = strconv.ParseFloat(match[2], 64)
	}

	// Extract requests/sec
	re = regexp.MustCompile(`Requests/sec:\s+([\d.]+)`)
	if match := re.FindStringSubmatch(text); match != nil {
		metrics.RequestsPerSec, _ = strconv.ParseFloat(match[1], 64)
	}

	// Extract average response time
	re = regexp.MustCompile(`Average:\s+([\d.]+) secs`)
	if match := re.FindStringSubmatch(text); match != nil {
		metrics.AvgResponseTime, _ = strconv.ParseFloat(match[1], 64)
	}

	// Extract slowest response
	re = regexp.MustCompile(`Slowest:\s+([\d.]+) secs`)
	if match := re.FindStringSubmatch(text); match != nil {
		metrics.Slowest, _ = strconv.ParseFloat(match[1], 64)
	}

	// Extract fastest response
	re = regexp.MustCompile(`Fastest:\s+([\d.]+) secs`)
	if match := re.FindStringSubmatch(text); match != nil {
		metrics.Fastest, _ = strconv.ParseFloat(match[1], 64)
	}

	// Extract latency distribution
	re = regexp.MustCompile(`\s+([\d.]+)%\s+in\s+([\d.]+)\s+secs`)
	for _, match := range re.FindAllStringSubmatch(text, -1) {
		percentile := match[1]
		time, _ := strconv.ParseFloat(match[2], 64)
		metrics.LatencyDistribution["p"+percentile] = time
	}

	// Extract status code distribution
	re = regexp.MustCompile(`\[(\d+)\]\s+(\d+) responses`)
	statusSection := extractSection(text, "Status code distribution:")
	for _, match := range re.FindAllStringSubmatch(statusSection, -1) {
		code := match[1]
		count, _ := strconv.Atoi(match[2])
		metrics.StatusCodes[code] = count
	}

	// Extract error distribution
	re = regexp.MustCompile(`\[(\d+)\]\s+(.+)`)
	errorSection := extractSection(text, "Error distribution:")
	for _, match := range re.FindAllStringSubmatch(errorSection, -1) {
		count, _ := strconv.Atoi(match[1])
		errorMsg := strings.TrimSpace(match[2])
		metrics.Errors[errorMsg] = count
	}

	return metrics, nil
}

// extractSection extracts a section of text after a header
func extractSection(text, header string) string {
	lines := strings.Split(text, "\n")
	inSection := false
	var section []string

	for _, line := range lines {
		if strings.Contains(line, header) {
			inSection = true
			continue
		}
		if inSection {
			if strings.TrimSpace(line) == "" || (len(line) > 0 && line[0] != ' ' && line[0] != '\t' && line[0] != '[') {
				break
			}
			section = append(section, line)
		}
	}

	return strings.Join(section, "\n")
}

// loadResults loads all result files from the results directory
func loadResults(resultsDir string) (BenchmarkData, error) {
	data := make(BenchmarkData)
	pattern := regexp.MustCompile(`test_results_(.+?)_(test|auth)_(\d+)\.txt`)

	files, err := filepath.Glob(filepath.Join(resultsDir, "test_results_*.txt"))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no result files found in '%s'", resultsDir)
	}

	for _, file := range files {
		basename := filepath.Base(file)
		match := pattern.FindStringSubmatch(basename)
		if match == nil {
			continue
		}

		gateway := match[1]
		endpoint := match[2]
		concurrency, _ := strconv.Atoi(match[3])

		metrics, err := parseHeyOutput(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", basename, err)
			continue
		}

		if data[gateway] == nil {
			data[gateway] = make(map[string]map[int]*Metrics)
		}
		if data[gateway][endpoint] == nil {
			data[gateway][endpoint] = make(map[int]*Metrics)
		}
		data[gateway][endpoint][concurrency] = metrics
	}

	return data, nil
}

// generateMarkdownReport generates a comprehensive markdown report
func generateMarkdownReport(data BenchmarkData) string {
	if len(data) == 0 {
		return "No benchmark data available.\n"
	}

	var report strings.Builder

	report.WriteString("# API Gateway Benchmark Results\n\n")
	report.WriteString("Total Requests: 100,000 per test\n\n")
	report.WriteString("Endpoints Tested: /test (unauthenticated), /auth (JWT-protected)\n\n")
	report.WriteString("Concurrency Levels: 50, 100, 250, 500\n\n")

	gateways := getSortedGateways(data)
	endpoints := []string{"test", "auth"}
	concurrencyLevels := []int{50, 100, 250, 500}

	// Summary table for each endpoint
	for _, endpoint := range endpoints {
		report.WriteString(fmt.Sprintf("## Results for /%s\n\n", endpoint))

		// Requests/sec comparison
		report.WriteString("### Throughput (Requests/sec)\n\n")
		report.WriteString("| Gateway | C=50 | C=100 | C=250 | C=500 |\n")
		report.WriteString("|---------|------|-------|-------|-------|\n")

		for _, gateway := range gateways {
			row := []string{gateway}
			for _, c := range concurrencyLevels {
				if metrics, ok := getMetrics(data, gateway, endpoint, c); ok {
					row = append(row, fmt.Sprintf("%.2f", metrics.RequestsPerSec))
				} else {
					row = append(row, "N/A")
				}
			}
			report.WriteString(fmt.Sprintf("| %s |\n", strings.Join(row, " | ")))
		}
		report.WriteString("\n")

		// Average response time comparison
		report.WriteString("### Average Response Time (seconds)\n\n")
		report.WriteString("| Gateway | C=50 | C=100 | C=250 | C=500 |\n")
		report.WriteString("|---------|------|-------|-------|-------|\n")

		for _, gateway := range gateways {
			row := []string{gateway}
			for _, c := range concurrencyLevels {
				if metrics, ok := getMetrics(data, gateway, endpoint, c); ok {
					row = append(row, fmt.Sprintf("%.4f", metrics.AvgResponseTime))
				} else {
					row = append(row, "N/A")
				}
			}
			report.WriteString(fmt.Sprintf("| %s |\n", strings.Join(row, " | ")))
		}
		report.WriteString("\n")

		// Latency percentiles for highest concurrency
		report.WriteString("### Latency Distribution at C=500 (seconds)\n\n")
		report.WriteString("| Gateway | 10% | 50% | 75% | 90% | 95% | 99% |\n")
		report.WriteString("|---------|-----|-----|-----|-----|-----|-----|\n")

		for _, gateway := range gateways {
			row := []string{gateway}
			if metrics, ok := getMetrics(data, gateway, endpoint, 500); ok {
				for _, p := range []string{"10", "50", "75", "90", "95", "99"} {
					key := "p" + p
					if val, exists := metrics.LatencyDistribution[key]; exists {
						row = append(row, fmt.Sprintf("%.4f", val))
					} else {
						row = append(row, "N/A")
					}
				}
			} else {
				for i := 0; i < 6; i++ {
					row = append(row, "N/A")
				}
			}
			report.WriteString(fmt.Sprintf("| %s |\n", strings.Join(row, " | ")))
		}
		report.WriteString("\n")
	}

	// Winner analysis
	report.WriteString("## Performance Summary\n\n")

	for _, endpoint := range endpoints {
		report.WriteString(fmt.Sprintf("### /%s Endpoint\n\n", endpoint))

		// Find best throughput at each concurrency
		for _, c := range concurrencyLevels {
			bestGateway := ""
			bestRPS := 0.0

			for _, gateway := range gateways {
				if metrics, ok := getMetrics(data, gateway, endpoint, c); ok {
					if metrics.RequestsPerSec > bestRPS {
						bestRPS = metrics.RequestsPerSec
						bestGateway = gateway
					}
				}
			}

			if bestGateway != "" {
				report.WriteString(fmt.Sprintf("**Concurrency %d**: %s - %.2f req/s\n\n", c, bestGateway, bestRPS))
			}
		}
	}

	// Error analysis
	report.WriteString("## Error Analysis\n\n")
	hasErrors := false

	for _, gateway := range gateways {
		for _, endpoint := range endpoints {
			for _, c := range concurrencyLevels {
				if metrics, ok := getMetrics(data, gateway, endpoint, c); ok {
					// Check for non-200 status codes
					errorCodes := make(map[string]int)
					for code, count := range metrics.StatusCodes {
						if code != "200" {
							errorCodes[code] = count
						}
					}

					if len(metrics.Errors) > 0 || len(errorCodes) > 0 {
						hasErrors = true
						report.WriteString(fmt.Sprintf("**%s** (/%s, C=%d):\n", gateway, endpoint, c))
						if len(errorCodes) > 0 {
							report.WriteString(fmt.Sprintf("  - Status codes: %v\n", errorCodes))
						}
						if len(metrics.Errors) > 0 {
							report.WriteString(fmt.Sprintf("  - Errors: %v\n", metrics.Errors))
						}
						report.WriteString("\n")
					}
				}
			}
		}
	}

	if !hasErrors {
		report.WriteString("No errors detected in any test runs.\n\n")
	}

	return report.String()
}

// getSortedGateways returns a sorted list of gateway names
func getSortedGateways(data BenchmarkData) []string {
	gateways := make([]string, 0, len(data))
	for gateway := range data {
		gateways = append(gateways, gateway)
	}
	sort.Strings(gateways)
	return gateways
}

// getMetrics retrieves metrics for a specific gateway, endpoint, and concurrency level
func getMetrics(data BenchmarkData, gateway, endpoint string, concurrency int) (*Metrics, bool) {
	if data[gateway] == nil {
		return nil, false
	}
	if data[gateway][endpoint] == nil {
		return nil, false
	}
	metrics, ok := data[gateway][endpoint][concurrency]
	return metrics, ok
}

func main() {
	resultsDir := "results"

	// Load results
	data, err := loadResults(resultsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading results: %v\n", err)
		os.Exit(1)
	}

	// Generate markdown report
	markdown := generateMarkdownReport(data)
	outputFile := filepath.Join(resultsDir, "BENCHMARK_REPORT.md")
	if err := os.WriteFile(outputFile, []byte(markdown), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing markdown report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Markdown report generated: %s\n", outputFile)
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Print(markdown)

	// Generate JSON report
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
		os.Exit(1)
	}
	jsonFile := filepath.Join(resultsDir, "benchmark_data.json")
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing JSON report: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("JSON report generated: %s\n", jsonFile)
}
