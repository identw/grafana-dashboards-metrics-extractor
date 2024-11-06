package parse

import (
	"fmt"
	"os"

	pLabels "github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/promql/parser"
)

func ExtractMetricsFromExpression(promExpr string) []string {
   

	expr, err := parser.ParseExpr(promExpr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ExtractMetricsFromExpression] parse error '%s': %v\n", promExpr, err)
		return nil
	}

    uniqMetrics := make(map[string]struct{})
    extractUniqMetricsFromExpression(expr, uniqMetrics)
    metrics := make([]string, 0, len(uniqMetrics))
    for m := range uniqMetrics {
        metrics = append(metrics, m)
    }
    return metrics
}

func extractUniqMetricsFromExpression(expr parser.Expr, metrics map[string]struct{}) {

    labels := parser.ExtractSelectors(expr)
    for _, lmatchers := range labels {
         for _, l := range lmatchers {
                if l.Name == "__name__" && l.Type == pLabels.MatchRegexp {
                    metrics[l.Value] = struct{}{}
                }
                if l.Name == "__name__" && l.Type == pLabels.MatchEqual {
                    metrics[l.Value] = struct{}{}
                }
         }
    }
	switch e := expr.(type) {
	case *parser.VectorSelector:
		metrics[e.Name] = struct{}{}
    case *parser.ParenExpr:
        extractUniqMetricsFromExpression(e.Expr, metrics)
    case *parser.AggregateExpr:
        extractUniqMetricsFromExpression(e.Expr, metrics)
    case *parser.MatrixSelector:
        extractUniqMetricsFromExpression(e.VectorSelector, metrics)
	case *parser.Call:
		for _, arg := range e.Args {
			extractUniqMetricsFromExpression(arg, metrics)
		}
	case *parser.BinaryExpr:
		extractUniqMetricsFromExpression(e.LHS, metrics)
		extractUniqMetricsFromExpression(e.RHS, metrics)
    case *parser.NumberLiteral:
        
	default:
		fmt.Fprintf(os.Stderr, "[ExtractMetricsFromExpression] not supported types of expression: %T\n", expr)
	}

}

