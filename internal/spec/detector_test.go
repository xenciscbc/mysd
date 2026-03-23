package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectSpecDir_MySD(t *testing.T) {
	dir, flavor, err := DetectSpecDir("../../testdata/fixtures/mysd-project")
	require.NoError(t, err)
	assert.Equal(t, ".specs", dir)
	assert.Equal(t, FlavorMySD, flavor)
}

func TestDetectSpecDir_OpenSpec(t *testing.T) {
	dir, flavor, err := DetectSpecDir("../../testdata/fixtures/openspec-project")
	require.NoError(t, err)
	assert.Equal(t, "openspec", dir)
	assert.Equal(t, FlavorOpenSpec, flavor)
}

func TestDetectSpecDir_NotFound(t *testing.T) {
	_, flavor, err := DetectSpecDir("/nonexistent/path/that/does/not/exist")
	assert.ErrorIs(t, err, ErrNoSpecDir)
	assert.Equal(t, FlavorNone, flavor)
}

func TestListChanges(t *testing.T) {
	changes, err := ListChanges("../../testdata/fixtures/mysd-project/.specs")
	require.NoError(t, err)
	assert.Contains(t, changes, "add-dark-mode")
}
