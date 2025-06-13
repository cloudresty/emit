#!/bin/bash

# Emit Logging Library Performance Benchmark Script
# Runs comprehensive benchmarks comparing emit, zap, and logrus
# Results exported to GitHub-friendly Markdown format

set -e

BENCHMARK_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=================================================="
echo "EMIT LOGGING LIBRARY BENCHMARK SUITE"
echo "=================================================="
echo "Modular Architecture: emit vs zap vs logrus"
echo "Platform: macOS"
echo "Date: $(date)"
echo "Results will be saved to: $RESULTS_DIR"
echo "=================================================="
echo

# Navigate to benchmark directory
cd "$BENCHMARK_DIR"

# Check dependencies
echo "Checking Go module dependencies..."
go mod tidy
echo "✅ Dependencies ready"
echo

# Build the benchmark application
echo "Building benchmark application..."
go build .
echo "✅ Build complete"
echo

# Set environment for consistent benchmarking
export GOMAXPROCS=$(sysctl -n hw.ncpu)
export GODEBUG=""

echo "Running benchmarks with $GOMAXPROCS CPU cores..."
echo "This may take several minutes..."
echo

# Run the benchmark application
echo "🚀 Starting comprehensive benchmark suite..."
./benchmarks

echo
echo "=================================================="
echo "BENCHMARK COMPLETE!"
echo "=================================================="
echo "Results saved to:"
echo "  📊 benchmark-results.md  - GitHub-friendly report"
echo "=================================================="
echo

# Display quick summary
if [ -f "benchmark-results.md" ]; then
    echo "QUICK PREVIEW:"
    echo "=============="
    # Show the first performance table
    sed -n '/## Performance Summary/,/^$/p' benchmark-results.md | head -15
fi

echo
echo "💡 View benchmark-results.md on GitHub for full interactive report"
echo "💡 All results include security vs performance analysis"
