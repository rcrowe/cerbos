package format_test

import (
	"testing"

	"github.com/cerbos/cerbos/internal/format"
	"github.com/cerbos/cerbos/internal/test"
	"github.com/stretchr/testify/require"
)

func TestYAML(t *testing.T) {
	testCases := test.LoadTestCases(t, "format")

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			f, err := format.YAML(testCase.Input)
			require.NoError(t, err)
			require.Equal(t, string(testCase.Want["out"]), string(f))
		})
	}
}
