// Package main contains the gocheck command.
package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/abemedia/gocheck/fieldorder"
	"github.com/abemedia/gocheck/untested"
)

func main() {
	multichecker.Main(fieldorder.NewAnalyzer(), untested.NewAnalyzer())
}
