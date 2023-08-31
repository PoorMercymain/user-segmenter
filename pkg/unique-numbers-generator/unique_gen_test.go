package uniquenumbersgenerator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateUniqueNonNegativeNumbers(t *testing.T) {
	aMap, err := GenerateUniqueNonNegativeNumbers(1, 0)
	require.Error(t, err)
	require.Empty(t, aMap)

	aMap, err = GenerateUniqueNonNegativeNumbers(1, 10)
	require.NoError(t, err)
	require.Len(t, aMap, 1)

	aMap, err = GenerateUniqueNonNegativeNumbers(10, 100)
	require.NoError(t, err)
	require.Len(t, aMap, 10)

	aMap, err = GenerateUniqueNonNegativeNumbers(100, 10)
	require.Error(t, err)
	require.Empty(t, aMap)
}
