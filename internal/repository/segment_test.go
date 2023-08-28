package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSegment(t *testing.T) {
	seg := NewSegment(nil)
	require.Empty(t, seg)
}
