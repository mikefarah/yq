package yqlib

import (
	"regexp"
	"strings"
	"time"
)

// DateTimePreprocessor handles automatic timestamp tagging for Goccy parser compatibility
type DateTimePreprocessor struct {
	enabled bool
}

// NewDateTimePreprocessor creates a new datetime preprocessor
func NewDateTimePreprocessor(enabled bool) *DateTimePreprocessor {
	return &DateTimePreprocessor{enabled: enabled}
}

// ISO8601 date/time patterns that should be automatically tagged as timestamps
var dateTimePatterns = []*regexp.Regexp{
	// RFC3339 / ISO8601 with timezone: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00
	regexp.MustCompile(`^\s*([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(?:\.[0-9]+)?(?:Z|[+-][0-9]{2}:[0-9]{2}))\s*$`),
	// Date only: 2006-01-02
	regexp.MustCompile(`^\s*([0-9]{4}-[0-9]{2}-[0-9]{2})\s*$`),
	// RFC3339 without timezone: 2006-01-02T15:04:05
	regexp.MustCompile(`^\s*([0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}(?:\.[0-9]+)?)\s*$`),
}

// isValidDateTime checks if a string represents a valid ISO8601/RFC3339 datetime
func isValidDateTime(value string) bool {
	// Try to parse with RFC3339 first
	if _, err := time.Parse(time.RFC3339, value); err == nil {
		return true
	}

	// Try to parse date-only format
	if _, err := time.Parse("2006-01-02", value); err == nil {
		return true
	}

	// Try RFC3339 without timezone
	if _, err := time.Parse("2006-01-02T15:04:05", value); err == nil {
		return true
	}

	return false
}

// PreprocessDocument automatically adds !!timestamp tags to ISO8601 datetime strings
// in YAML documents when using Goccy parser for consistent datetime arithmetic behaviour
func (dtp *DateTimePreprocessor) PreprocessDocument(yamlContent string) string {
	if !dtp.enabled {
		return yamlContent
	}

	lines := strings.Split(yamlContent, "\n")
	var result strings.Builder

	for i, line := range lines {
		processed := dtp.processLine(line)
		result.WriteString(processed)

		// Add newline except for the last line (preserve original ending)
		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// processLine processes a single line of YAML, adding timestamp tags where appropriate
func (dtp *DateTimePreprocessor) processLine(line string) string {
	// Skip lines that are comments, already have tags, or are part of multi-line constructs
	if isSkippableLine(line) {
		return line
	}

	// Look for key-value pairs: "key: value" or "- value"
	trimLeft := strings.TrimLeft(line, " \t")
	if strings.HasPrefix(trimLeft, "- ") {
		// Handle array items first (before key-value pairs)
		return dtp.processArrayItemLine(line)
	} else if colonIndex := strings.Index(line, ":"); colonIndex != -1 {
		// Handle map key-value pairs
		return dtp.processKeyValueLine(line, colonIndex)
	}

	return line
}

// isSkippableLine checks if a line should be skipped from datetime preprocessing
func isSkippableLine(line string) bool {
	trimmed := strings.TrimSpace(line)

	// Skip empty lines, comments, directives, document separators
	if trimmed == "" || strings.HasPrefix(trimmed, "#") ||
		strings.HasPrefix(trimmed, "%") || strings.HasPrefix(trimmed, "---") ||
		strings.HasPrefix(trimmed, "...") {
		return true
	}

	// Skip lines that already have explicit tags
	if strings.Contains(line, "!!") {
		return true
	}

	// Skip multi-line scalar indicators
	if strings.HasSuffix(trimmed, "|") || strings.HasSuffix(trimmed, ">") ||
		strings.HasSuffix(trimmed, "|-") || strings.HasSuffix(trimmed, ">-") ||
		strings.HasSuffix(trimmed, "|+") || strings.HasSuffix(trimmed, ">+") {
		return true
	}

	return false
}

// processKeyValueLine processes a line containing a key-value pair
func (dtp *DateTimePreprocessor) processKeyValueLine(line string, colonIndex int) string {
	beforeColon := line[:colonIndex]
	afterColon := line[colonIndex+1:]

	// Check if the value part looks like a datetime
	trimmedValue := strings.TrimSpace(afterColon)

	// Skip if value is empty, quoted, or already complex
	if trimmedValue == "" ||
		strings.HasPrefix(trimmedValue, "\"") ||
		strings.HasPrefix(trimmedValue, "'") ||
		strings.HasPrefix(trimmedValue, "{") ||
		strings.HasPrefix(trimmedValue, "[") ||
		strings.HasPrefix(trimmedValue, "&") ||
		strings.HasPrefix(trimmedValue, "*") {
		return line
	}

	// Check if it matches datetime patterns and is valid
	if dtp.matchesDateTimePattern(trimmedValue) && isValidDateTime(trimmedValue) {
		// Insert !!timestamp tag before the value
		return beforeColon + ": !!timestamp " + trimmedValue
	}

	return line
}

// processArrayItemLine processes a line containing an array item
func (dtp *DateTimePreprocessor) processArrayItemLine(line string) string {
	// Find the position of "- " after any leading whitespace
	trimLeft := strings.TrimLeft(line, " \t")
	if !strings.HasPrefix(trimLeft, "- ") {
		return line
	}

	// Find the actual position of "- " in the original line
	leadingWhitespace := line[:len(line)-len(trimLeft)]
	dashIndex := len(leadingWhitespace)

	beforeDash := line[:dashIndex+2] // Include leading whitespace and "- "
	afterDash := line[dashIndex+2:]

	trimmedValue := strings.TrimSpace(afterDash)

	// Skip if value is empty, quoted, or already complex
	if trimmedValue == "" ||
		strings.HasPrefix(trimmedValue, "\"") ||
		strings.HasPrefix(trimmedValue, "'") ||
		strings.HasPrefix(trimmedValue, "{") ||
		strings.HasPrefix(trimmedValue, "[") ||
		strings.HasPrefix(trimmedValue, "&") ||
		strings.HasPrefix(trimmedValue, "*") {
		return line
	}

	// Check if it matches datetime patterns and is valid
	if dtp.matchesDateTimePattern(trimmedValue) && isValidDateTime(trimmedValue) {
		// Insert !!timestamp tag before the value
		return beforeDash + "!!timestamp " + trimmedValue
	}

	return line
}

// matchesDateTimePattern checks if a value matches any of the datetime regex patterns
func (dtp *DateTimePreprocessor) matchesDateTimePattern(value string) bool {
	for _, pattern := range dateTimePatterns {
		if pattern.MatchString(value) {
			return true
		}
	}
	return false
}
