package gocontemplate_test

import (
	"testing"

	"github.com/foomo/gocontemplate"
	_ "github.com/foomo/sesamy-go"              // force inclusion
	_ "github.com/foomo/sesamy-go/event/params" // force inclusion
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLoader(t *testing.T) {
	t.Parallel()
	goctpl, err := gocontemplate.Load(&gocontemplate.Config{
		Packages: []*gocontemplate.ConfigPackage{
			{
				Path:  "github.com/foomo/sesamy-go/event",
				Types: []string{"PageView"},
			},
		},
	})
	require.NoError(t, err)

	assert.Len(t, goctpl.Packages, 4)
}

func TestLoader_LookupTypesByType(t *testing.T) {
	t.Parallel()
	goctpl, err := gocontemplate.Load(&gocontemplate.Config{
		Packages: []*gocontemplate.ConfigPackage{
			{
				Path:  "github.com/foomo/sesamy-go/event",
				Types: []string{"PageView"},
			},
		},
	})
	require.NoError(t, err)

	pkg := goctpl.Package("github.com/foomo/sesamy-go")
	require.NotNil(t, pkg)
	pkgType := pkg.LookupType("Event")
	require.NotNil(t, pkgType)

	pkgTypes := goctpl.LookupTypesByType(pkgType)
	assert.NotEmpty(t, pkgTypes)
}
