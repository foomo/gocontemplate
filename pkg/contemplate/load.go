package contemplate

import (
	"context"
	"os"

	"golang.org/x/tools/go/packages"
)

// Load loads the packages listed in cfg, together with all their transitive
// imports, and scopes the configured types and the types they reference.
// It returns an error if the packages cannot be loaded.
func Load(ctx context.Context, cfg *Config) (*Contemplate, error) {
	inst := &Contemplate{
		cfg:      cfg,
		Packages: map[string]*Package{},
	}

	// load packages
	pkgs, err := packages.Load(&packages.Config{
		Context: ctx,
		Dir:     os.ExpandEnv(cfg.Directory),
		Mode: packages.NeedName | packages.NeedTypesInfo |
			packages.NeedFiles | packages.NeedImports | packages.NeedDeps |
			packages.NeedModule | packages.NeedTypes | packages.NeedSyntax,
	}, cfg.PackagePaths()...)
	if err != nil {
		return nil, err
	}

	inst.addPackages(pkgs...)
	inst.addPackagesConfigs(cfg.Packages...)

	return inst, nil
}
