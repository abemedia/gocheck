package untested_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/abemedia/goneat/untested"
)

func TestUntested(t *testing.T) {
	analyzer := untested.NewAnalyzer()
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "a/...")
}

func TestUntestedWithGenerated(t *testing.T) {
	analyzer := untested.NewAnalyzer()
	analyzer.Flags.Set("generated", "true")

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "b/...")
}

func TestUntestedWithInternal(t *testing.T) {
	analyzer := untested.NewAnalyzer()
	analyzer.Flags.Set("internal", "true")

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, analyzer, "c/...")
}
