// Package main contains the goneat command.
package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/abemedia/goneat/fieldorder"
	"github.com/abemedia/goneat/untested"
)

func main() {
	multichecker.Main(fieldorder.NewAnalyzer(), untested.NewAnalyzer())
}
