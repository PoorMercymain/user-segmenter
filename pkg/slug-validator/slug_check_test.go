package slugvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSlug(t *testing.T) {
	isSlug := IsSlug("~not~a~slug~")
	require.False(t, isSlug)

	isSlug = IsSlug("a-slug")
	require.True(t, isSlug)

	isSlug = IsSlug("A_SLUG")
	require.True(t, isSlug)
}
