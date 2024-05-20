package gocontemplate

type Config struct {
	Packages []*ConfigPackage `json:"packages" yaml:"packages"`
}

func (c *Config) Package(path string) *ConfigPackage {
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
