package contemplate

import (
	"go/ast"
	"go/types"
	"slices"

	"github.com/foomo/gocontemplate/pkg/assume"
	"golang.org/x/exp/maps"
	"golang.org/x/tools/go/packages"
)

type Contemplate struct {
	cfg      *Config
	Packages map[string]*Package
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (s *Contemplate) Package(path string) *Package {
	return s.Packages[path]
}

func (s *Contemplate) LookupExpr(name string) ast.Expr {
	for _, pkg := range s.Packages {
		if value := pkg.LookupExpr(name); value != nil {
			return value
		}
	}
	return nil
}

func (s *Contemplate) LookupTypesByType(obj types.Object) []types.Object {
	var ret []types.Object

	expr := assume.T[*ast.Ident](s.LookupExpr(obj.Name()))
	if expr == nil {
		return nil
	}

	for _, pkg := range s.Packages {
		for _, object := range pkg.Types() {
			switch objectType := object.(type) {
			case *types.Const:
				if objectTypeNamed := assume.T[*types.Named](objectType.Type()); objectTypeNamed != nil {
					if objectTypeNamed.Obj() == obj {
						ret = append(ret, objectType)
					}
				}
			case *types.TypeName:
				if objectExpr := pkg.LookupExpr(object.Name()); objectExpr != nil {
					if objectExprIdent := assume.T[*ast.Ident](objectExpr); objectExprIdent != nil {
						if objectExprDecl := assume.T[*ast.TypeSpec](objectExprIdent.Obj.Decl); objectExprDecl != nil {
							if objectExprType, ok := pkg.pkg.TypesInfo.Types[objectExprDecl.Type]; ok {
								if objectExprTypeNamed := assume.T[*types.Named](objectExprType.Type); objectExprTypeNamed != nil {
									if objectExprTypeNamed.Obj() == obj {
										ret = append(ret, objectType)
									}
								}
							}
						}
					}
				}
			default:
				// fmt.Println("?")
			}
		}
	}
	return ret
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (s *Contemplate) addPackages(pkgs ...*packages.Package) {
	for _, pkg := range pkgs {
		if _, ok := s.Packages[pkg.PkgPath]; !ok {
			s.Packages[pkg.PkgPath] = NewPackage(s, pkg)

			s.addPackages(maps.Values(pkg.Imports)...)
		}
	}
}

func (s *Contemplate) addPackagesConfigs(confs ...*PackageConfig) {
	for _, conf := range confs {
		s.Package(conf.Path).AddScopeTypes(conf.Types...)
	}
}

func (s *Contemplate) LookupAstIdentDefsByDeclType(input types.TypeAndValue) []types.Object {
	var pkgs []*packages.Package
	var addImports func(pkg *packages.Package)
	addImports = func(pkg *packages.Package) {
		for _, p := range pkg.Imports {
			if !slices.Contains(pkgs, p) {
				pkgs = append(pkgs, p)
				addImports(p)
			}
		}
	}

	var ret []types.Object
	for _, p := range pkgs {
		for _, name := range p.Types.Scope().Names() {
			child := p.Types.Scope().Lookup(name)
			if child.Type() == input.Type {
				ret = append(ret, child)
			}
		}
	}
	return ret
}

func (s *Contemplate) addPackageTypeNames(pkg *packages.Package, typeNames ...string) {
	if _, ok := s.Packages[pkg.PkgPath]; !ok {
		s.Packages[pkg.PkgPath] = NewPackage(s, pkg)
	}
	// add request scopes
	s.Packages[pkg.PkgPath].AddScopeTypes(typeNames...)
}
