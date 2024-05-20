package gocontemplate

type ConfigPackage struct {
	Path  string   `json:"path" yaml:"path"`
	Types []string `json:"types" yaml:"types"`
}
