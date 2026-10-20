package contemplate

// PackageConfig selects a package and the types to scope within it.
type PackageConfig struct {
	// Path is the import path of the package.
	Path string `json:"path" yaml:"path"`
	// Types lists the names of the types to scope. Types they reference are
	// scoped as well.
	Types []string `json:"types" yaml:"types"`
}
