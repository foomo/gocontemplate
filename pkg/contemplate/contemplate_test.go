package contemplate_test

import (
	"testing"

	"github.com/foomo/gocontemplate/pkg/contemplate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	t.Parallel()
	ctpl, err := contemplate.Load(t.Context(), &contemplate.Config{
		Packages: []*contemplate.PackageConfig{
			{
				Path:  "github.com/foomo/gocontemplate/test/event",
				Types: []string{"PageView"},
			},
		},
	})
	require.NoError(t, err)

	assert.Len(t, ctpl.Packages, 3)
}

func TestLoader_LookupTypesByType(t *testing.T) {
	t.Parallel()
	ctpl, err := contemplate.Load(t.Context(), &contemplate.Config{
		Packages: []*contemplate.PackageConfig{
			{
				Path:  "github.com/foomo/gocontemplate/test/event",
				Types: []string{"PageView"},
			},
		},
	})
	require.NoError(t, err)

	pkg := ctpl.Package("github.com/foomo/gocontemplate/test")
	require.NotNil(t, pkg)
	pkgType := pkg.LookupType("Event")
	require.NotNil(t, pkgType)

	pkgTypes := ctpl.LookupTypesByType(pkgType)
	assert.NotEmpty(t, pkgTypes)
}
