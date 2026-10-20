package contemplate

// Config describes which packages and types [Load] loads.
type Config struct {
	// Directory is the directory containing the go.mod file. Environment
	// variables are expanded. An empty value uses the current directory.
	Directory string `json:"directory" yaml:"directory"`
	// Packages lists the packages to load and the types to scope in each.
	Packages []*PackageConfig `json:"packages" yaml:"packages"`
}

// Package returns the [PackageConfig] with the given import path, or nil if
// none matches.
func (c *Config) Package(path string) *PackageConfig {
	for _, value := range c.Packages {
		if value.Path == path {
			return value
		}
	}

	return nil
}

// PackagePaths returns the import paths of all configured packages, in
// configuration order.
func (c *Config) PackagePaths() []string {
	ret := make([]string, len(c.Packages))
	for i, value := range c.Packages {
		ret[i] = value.Path
	}

	return ret
}
