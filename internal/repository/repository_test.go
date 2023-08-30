package repository

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSegment(t *testing.T) {
	seg := NewSegment(nil)
	require.Empty(t, seg)

	usr := NewUser(nil)
	require.Empty(t, usr)

	rep := NewReport(nil)
	require.Empty(t, rep)

	pg := NewPostgres(nil)
	require.Empty(t, pg)
}
