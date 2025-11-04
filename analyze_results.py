#!/usr/bin/env python3
"""
Benchmark Results Analyzer for API Gateway Performance Tests

Parses 'hey' output and generates comparison reports.
"""

import os
import re
import json
from pathlib import Path
from collections import defaultdict
import sys


class BenchmarkAnalyzer:
    def __init__(self, results_dir="results"):
        self.results_dir = Path(results_dir)
        self.data = defaultdict(lambda: defaultdict(dict))

    def parse_hey_output(self, filepath):
        """Parse hey output file and extract metrics."""
        with open(filepath, 'r') as f:
            content = f.read()

        metrics = {}

        # Extract total requests
        match = re.search(r'(\d+) requests in ([\d.]+)s', content)
        if match:
            metrics['total_requests'] = int(match.group(1))
            metrics['total_time'] = float(match.group(2))

        # Extract requests/sec
        match = re.search(r'Requests/sec:\s+([\d.]+)', content)
        if match:
            metrics['requests_per_sec'] = float(match.group(1))

        # Extract average response time
        match = re.search(r'Average:\s+([\d.]+) secs', content)
        if match:
            metrics['avg_response_time'] = float(match.group(1))

        # Extract slowest response
        match = re.search(r'Slowest:\s+([\d.]+) secs', content)
        if match:
            metrics['slowest'] = float(match.group(1))

        # Extract fastest response
        match = re.search(r'Fastest:\s+([\d.]+) secs', content)
        if match:
            metrics['fastest'] = float(match.group(1))

        # Extract latency distribution
        latency_dist = {}
        for line in content.split('\n'):
            match = re.match(r'\s+([\d.]+)%\s+in\s+([\d.]+)\s+secs', line)
            if match:
                percentile = match.group(1)
                time = float(match.group(2))
                latency_dist[f'p{percentile}'] = time
        metrics['latency_distribution'] = latency_dist

        # Extract status code distribution
        status_codes = {}
        match = re.search(r'Status code distribution:\s*\n((?:\s+\[.*?\]\s+\d+\s+responses\s*\n?)*)', content)
        if match:
            for line in match.group(1).split('\n'):
                code_match = re.search(r'\[(\d+)\]\s+(\d+) responses', line)
                if code_match:
                    status_codes[code_match.group(1)] = int(code_match.group(2))
        metrics['status_codes'] = status_codes

        # Extract error distribution
        errors = {}
        match = re.search(r'Error distribution:\s*\n((?:\s+\[.*?\]\s+\d+\s*\n?)*)', content)
        if match:
            for line in match.group(1).split('\n'):
                error_match = re.search(r'\[(\d+)\]\s+(.+)', line)
                if error_match:
                    errors[error_match.group(2).strip()] = int(error_match.group(1))
        metrics['errors'] = errors

        return metrics

    def load_results(self):
        """Load all result files from the results directory."""
        pattern = re.compile(r'test_results_(.+?)_(test|auth)_(\d+)\.txt')

        if not self.results_dir.exists():
            print(f"Results directory '{self.results_dir}' not found!")
            return False

        files = list(self.results_dir.glob('test_results_*.txt'))
        if not files:
            print(f"No result files found in '{self.results_dir}'!")
            return False

        for filepath in files:
            match = pattern.match(filepath.name)
            if match:
                gateway = match.group(1)
                endpoint = match.group(2)
                concurrency = int(match.group(3))

                try:
                    metrics = self.parse_hey_output(filepath)
                    self.data[gateway][endpoint][concurrency] = metrics
                except Exception as e:
                    print(f"Error parsing {filepath.name}: {e}")

        return True

    def generate_markdown_report(self):
        """Generate a comprehensive markdown report."""
        if not self.data:
            return "No benchmark data available."

        report = []
        report.append("# API Gateway Benchmark Results\n")
        report.append(f"Total Requests: 100,000 per test\n")
        report.append(f"Endpoints Tested: /test (unauthenticated), /auth (JWT-protected)\n")
        report.append(f"Concurrency Levels: 50, 100, 250, 500\n\n")

        gateways = sorted(self.data.keys())
        endpoints = ['test', 'auth']
        concurrency_levels = [50, 100, 250, 500]

        # Summary table for each endpoint
        for endpoint in endpoints:
            report.append(f"## Results for /{endpoint}\n")

            # Requests/sec comparison
            report.append(f"### Throughput (Requests/sec)\n")
            report.append("| Gateway | C=50 | C=100 | C=250 | C=500 |\n")
            report.append("|---------|------|-------|-------|-------|\n")

            for gateway in gateways:
                row = [gateway]
                for c in concurrency_levels:
                    metrics = self.data[gateway].get(endpoint, {}).get(c, {})
                    rps = metrics.get('requests_per_sec', 0)
                    row.append(f"{rps:,.2f}" if rps > 0 else "N/A")
                report.append("| " + " | ".join(row) + " |\n")
            report.append("\n")

            # Average response time comparison
            report.append(f"### Average Response Time (seconds)\n")
            report.append("| Gateway | C=50 | C=100 | C=250 | C=500 |\n")
            report.append("|---------|------|-------|-------|-------|\n")

            for gateway in gateways:
                row = [gateway]
                for c in concurrency_levels:
                    metrics = self.data[gateway].get(endpoint, {}).get(c, {})
                    avg = metrics.get('avg_response_time', 0)
                    row.append(f"{avg:.4f}" if avg > 0 else "N/A")
                report.append("| " + " | ".join(row) + " |\n")
            report.append("\n")

            # Latency percentiles for highest concurrency
            report.append(f"### Latency Distribution at C=500 (seconds)\n")
            report.append("| Gateway | 10% | 50% | 75% | 90% | 95% | 99% |\n")
            report.append("|---------|-----|-----|-----|-----|-----|-----|\n")

            for gateway in gateways:
                metrics = self.data[gateway].get(endpoint, {}).get(500, {})
                latency = metrics.get('latency_distribution', {})
                row = [gateway]
                for p in ['10', '50', '75', '90', '95', '99']:
                    val = latency.get(f'p{p}', 0)
                    row.append(f"{val:.4f}" if val > 0 else "N/A")
                report.append("| " + " | ".join(row) + " |\n")
            report.append("\n")

        # Winner analysis
        report.append("## Performance Summary\n\n")

        for endpoint in endpoints:
            report.append(f"### /{endpoint} Endpoint\n\n")

            # Find best throughput at each concurrency
            for c in concurrency_levels:
                best_gateway = None
                best_rps = 0

                for gateway in gateways:
                    metrics = self.data[gateway].get(endpoint, {}).get(c, {})
                    rps = metrics.get('requests_per_sec', 0)
                    if rps > best_rps:
                        best_rps = rps
                        best_gateway = gateway

                if best_gateway:
                    report.append(f"**Concurrency {c}**: {best_gateway} - {best_rps:,.2f} req/s\n\n")

        # Error analysis
        report.append("## Error Analysis\n\n")
        has_errors = False

        for gateway in gateways:
            for endpoint in endpoints:
                for c in concurrency_levels:
                    metrics = self.data[gateway].get(endpoint, {}).get(c, {})
                    errors = metrics.get('errors', {})
                    status_codes = metrics.get('status_codes', {})

                    # Check for non-200 status codes
                    error_codes = {k: v for k, v in status_codes.items() if k != '200'}

                    if errors or error_codes:
                        has_errors = True
                        report.append(f"**{gateway}** (/{endpoint}, C={c}):\n")
                        if error_codes:
                            report.append(f"  - Status codes: {error_codes}\n")
                        if errors:
                            report.append(f"  - Errors: {errors}\n")
                        report.append("\n")

        if not has_errors:
            report.append("No errors detected in any test runs.\n\n")

        return "".join(report)

    def generate_json_report(self):
        """Generate a JSON report of all data."""
        return json.dumps(dict(self.data), indent=2)


def main():
    analyzer = BenchmarkAnalyzer()

    if not analyzer.load_results():
        sys.exit(1)

    # Generate markdown report
    markdown = analyzer.generate_markdown_report()
    output_file = Path("results/BENCHMARK_REPORT.md")
    output_file.write_text(markdown)
    print(f"Markdown report generated: {output_file}")
    print("\n" + "="*80 + "\n")
    print(markdown)

    # Generate JSON report
    json_report = analyzer.generate_json_report()
    json_file = Path("results/benchmark_data.json")
    json_file.write_text(json_report)
    print(f"JSON report generated: {json_file}")


if __name__ == "__main__":
    main()
