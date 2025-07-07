// Package main contains the goneat command.
package main

import (
	"golang.org/x/tools/go/analysis/multichecker"

	"github.com/abemedia/goneat/fieldorder"
)

func main() {
	multichecker.Main(
		fieldorder.NewAnalyzer(),
	)
}
