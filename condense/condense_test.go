package condense_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/abemedia/goneat/condense"
)

func TestCondense(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, condense.NewAnalyzer(), "condense")
}
