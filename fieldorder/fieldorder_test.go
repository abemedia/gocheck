package fieldorder_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/abemedia/gocheck/fieldorder"
)

func TestFieldOrder(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, fieldorder.NewAnalyzer(), "fieldorder")
}
