package parse

import (
	"regexp"
)

var (
	PrometheusLabelNamesRegex           = regexp.MustCompile(`^label_names\(\)\s*$`)
	PrometheusLabelValuesRegex          = regexp.MustCompile(`^label_values\((?:(.+),\s*)?([a-zA-Z_$][a-zA-Z0-9_]*)\)\s*$`)
	PrometheusMetricNamesRegex          = regexp.MustCompile(`^metrics\((.+)\)\s*$`)
	PrometheusQueryResultRegex          = regexp.MustCompile(`^query_result\((.+)\)\s*$`)
	PrometheusLabelNamesRegexWithMatch  = regexp.MustCompile(`^label_names\((.+)\)\s*$`)
)

func ParseVariable(expr string) (string) {
	switch {
	case PrometheusLabelNamesRegex.MatchString(expr):
		return ""
	case PrometheusLabelValuesRegex.MatchString(expr):
		matches := PrometheusLabelValuesRegex.FindStringSubmatch(expr)
		return matches[1]
	case PrometheusMetricNamesRegex.MatchString(expr):
		matches := PrometheusMetricNamesRegex.FindStringSubmatch(expr)
		return matches[1]
	case PrometheusQueryResultRegex.MatchString(expr):
		matches := PrometheusQueryResultRegex.FindStringSubmatch(expr)
		return matches[1]
	case PrometheusLabelNamesRegexWithMatch.MatchString(expr):
		matches := PrometheusLabelNamesRegexWithMatch.FindStringSubmatch(expr)
		return matches[1]
	default:
		return  ""
	}
}
