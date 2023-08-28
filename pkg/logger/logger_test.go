package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	logger, err := GetLogger()
	require.Empty(t, logger)
	require.Error(t, err)

	err = InitLogger()
	require.NoError(t, err)

	logger, err = GetLogger()
	require.NotEmpty(t, logger)
	require.NoError(t, err)
}
