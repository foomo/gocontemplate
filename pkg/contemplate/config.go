package contemplate

type Config struct {
	// Directory containing the go.mod file
	Directory string `json:"directory" yaml:"directory"`
	// Packages to load
	Packages []*PackageConfig `json:"packages" yaml:"packages"`
}

func (c *Config) Package(path string) *PackageConfig {
	for _, value := range c.Packages {
		if value.Path == path {
			return value
		}
	}
	return nil
}

func (c *Config) PackagePaths() []string {
	ret := make([]string, len(c.Packages))
	for i, value := range c.Packages {
		ret[i] = value.Path
	}
	return ret
}
