package fieldgetter

import (
	"regexp"
)

func GetExprs(data interface{}) []string {
	var exprs []string
	switch v := data.(type) {
	case []interface{}:
		for _, item := range v {
			exprs = append(exprs, GetExprs(item)...)
		}
	case map[string]interface{}:
		if expr, ok := v["expr"]; ok {

			exprs = append(exprs, replaceVariables(expr.(string)))
		}
		for _, value := range v {
			exprs = append(exprs, GetExprs(value)...)
		}
	}
	return exprs
}

func GetDefinition(data interface{}) []string {
	var  defs []string
	switch v := data.(type) {
	case []interface{}:
		for _, item := range v {
			defs = append(defs, GetDefinition(item)...)
		}
	case map[string]interface{}:
		if definition, ok := v["definition"]; ok {
			defs = append(defs, replaceVariables(definition.(string)))
		}
		for _, value := range v {
			defs = append(defs, GetDefinition(value)...)
		}
	}
	return defs
}

func replaceVariables(s string) string {
	re := regexp.MustCompile(`\$\{?\w+\}?`)
	return re.ReplaceAllString(s, "10m")
}

