#!/bin/bash

set -e

echo "Running tests and generating coverage..."
packages=$(go list ./... | grep -v -E 'examples' | grep -v -E 'test' | tr '\n' ',' | sed 's/,$//')
test_packages=$(go list ./... | grep -v -E 'examples' | grep -v -E 'test' | grep -v '^github.com/mikefarah/yq/v4$')
go test -coverprofile=coverage.out -coverpkg="$packages" -v $test_packages

echo "Generating HTML coverage report..."
go tool cover -html=coverage.out -o coverage.html

echo ""
echo "Generating sorted coverage table..."

# Create a simple approach using grep and sed to extract file coverage
# First, get the total coverage
total_coverage=$(go tool cover -func=coverage.out | grep "^total:" | sed 's/.*([^)]*)[[:space:]]*\([0-9.]*\)%.*/\1/')

# Extract file-level coverage by finding the last occurrence of each file
go tool cover -func=coverage.out | grep -E "\.go:[0-9]+:" | \
sed 's/^\([^:]*\.go\):.*[[:space:]]\([0-9.]*\)%.*/\2 \1/' | \
sort -k2 | \
awk '{file_coverage[$2] = $1} END {for (file in file_coverage) printf "%.2f %s\n", file_coverage[file], file}' | \
sort -nr > coverage_sorted.txt

# Add total coverage to the file
if [[ -n "$total_coverage" && "$total_coverage" != "0" ]]; then
    echo "TOTAL: $total_coverage" >> coverage_sorted.txt
fi

echo ""
echo "Coverage Summary (sorted by percentage - lowest coverage first):"
echo "================================================================="
printf "%-60s %10s %12s\n" "FILE" "COVERAGE" "STATUS"
echo "================================================================="

# Display results with status indicators
tail -n +1 coverage_sorted.txt | while read percent file; do
    if [[ "$file" == "TOTAL:" ]]; then
        echo ""
        printf "%-60s %8s%% %12s\n" "OVERALL PROJECT COVERAGE" "$percent" "📊 TOTAL"
        echo "================================================================="
        continue
    fi
    
    filename=$(basename "$file")
    status=""
    if (( $(echo "$percent < 50" | bc -l 2>/dev/null || echo "0") )); then
        status="🔴 CRITICAL"
    elif (( $(echo "$percent < 70" | bc -l 2>/dev/null || echo "0") )); then
        status="🟡 LOW"
    elif (( $(echo "$percent < 90" | bc -l 2>/dev/null || echo "0") )); then
        status="🟢 GOOD"
    else
        status="✅ EXCELLENT"
    fi
    
    printf "%-60s %8s%% %12s\n" "$filename" "$percent" "$status"
done

echo ""
echo "Top 10 files by uncovered statements:"
echo "================================================="
# Calculate uncovered statements for each file and sort by that
go tool cover -func=coverage.out | grep -E "\.go:[0-9]+:" | \
awk '{
    # Extract filename and percentage
    split($1, parts, ":")
    file = parts[1]
    pct = $NF
    gsub(/%/, "", pct)
    
    # Track stats per file
    total[file]++
    covered[file] += pct
}
END {
    for (file in total) {
        avg_pct = covered[file] / total[file]
        uncovered = total[file] * (100 - avg_pct) / 100
        covered_count = total[file] - uncovered
        printf "%.0f %d %.0f %.1f %s\n", uncovered, total[file], covered_count, avg_pct, file
    }
}' | sort -rn | head -10 | while read uncovered total covered pct file; do
    filename=$(basename "$file")
    printf "%-60s %4d uncovered (%4d/%4d, %5.1f%%)\n" "$filename" "$uncovered" "$covered" "$total" "$pct"
done

echo ""
echo "Coverage reports generated:"
echo "- HTML report: coverage.html (detailed line-by-line coverage)"
echo "- Sorted table: coverage_sorted.txt"
echo "- Use 'go tool cover -func=coverage.out' for function-level details"
