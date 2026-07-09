// Линтер, который проверяет код на стандартные staticcheck правила + unusedwrite + unusedresult + unreachable + waitgroup
//
// Для получения инструкции по использованию:
//
//	$ staticlint --help
package main

import (
	"github.com/WorstOfAny/go-musthave-metrics-tpl/cmd/staticlint/osexitcheckanalyzer"
	"github.com/nishanths/exhaustive"
	"github.com/sonatard/noctx"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	mychecks := []*analysis.Analyzer{
		osexitcheckanalyzer.Analyzer,
		unusedwrite.Analyzer,
		unusedresult.Analyzer,
		unreachable.Analyzer,
		waitgroup.Analyzer,
		// Анаплизатор для работы с switch enum
		exhaustive.Analyzer,
		// Анаплизатор для поиска запросов без контекста
		noctx.Analyzer,
	}

	// Анаплизаторы для поиска багов и проблем с производительностью
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}
	// Анализатор для simplify кода
	for _, v := range simple.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	multichecker.Main(
		mychecks...,
	)
}
