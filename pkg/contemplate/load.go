package contemplate

import (
	"golang.org/x/tools/go/packages"
)

func Load(cfg *Config) (*Contemplate, error) {
	inst := &Contemplate{
		cfg:      cfg,
		Packages: map[string]*Package{},
	}

	// load packages
	pkgs, err := packages.Load(&packages.Config{
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
