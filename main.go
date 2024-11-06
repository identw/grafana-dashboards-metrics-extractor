package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"flag"
	"path/filepath"

	"github.com/identw/grafana-dashboards-metrics-extractor/pkg/fieldgetter"
	"github.com/identw/grafana-dashboards-metrics-extractor/pkg/parse"
)

var (
	dir = flag.String("dir", "./", "directory path to grafana dashboards")
)
func main() {

	flag.Parse()

	uniqMetrics := make(map[string]struct{})
	err := filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			metrics, err := processDashboard(path)
			if err != nil {
				log.Printf("[main] process error %s: %s", path, err)
			}

			for _, metric := range metrics {
				if metric == "" {
					continue
				}
				uniqMetrics[metric] = struct{}{}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("[main] directory walk error: %s", err)
	}

	for metric := range uniqMetrics {
		fmt.Println(metric)
	}

}

func processDashboard(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("[processDashboard] read file error: %w", err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, fmt.Errorf("[processDashboard]: parse json error: file: %v error: %w", filePath, err)
	}

	
	exprs := fieldgetter.GetExprs(jsonData)
	metrics := make([]string, 0)
	for _, expr := range exprs {
		metrics = append(metrics, parse.ExtractMetricsFromExpression(expr)...)
	}
	defs := fieldgetter.GetDefinition(jsonData)
	for _, def := range defs {
		if def == "" {
			continue
		}
		promExpr := parse.ParseVariable(def)
		metrics = append(metrics, parse.ExtractMetricsFromExpression(promExpr)...)

	}

	return metrics, nil
}




