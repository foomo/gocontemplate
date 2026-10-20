package contemplate

import (
	"go/ast"
	"go/types"
	"maps"
	"slices"

	"github.com/foomo/gocontemplate/pkg/assume"
	"golang.org/x/tools/go/packages"
)

// Contemplate is the result of [Load]. Create it with [Load]; the zero value
// is not usable. It is not safe for concurrent modification.
type Contemplate struct {
	cfg *Config
	// Packages maps import paths to all loaded packages, including
	// transitive imports.
	Packages map[string]*Package
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

// Package returns the loaded package with the given import path, or nil if
// it was not loaded.
func (s *Contemplate) Package(path string) *Package {
	return s.Packages[path]
}

// LookupExpr returns the defining identifier of the package-level constant,
// function or type called name, searching all loaded packages. If several
// packages define name, which one is returned is unspecified. It returns nil
// if no package defines name.
func (s *Contemplate) LookupExpr(name string) ast.Expr {
	for _, pkg := range s.Packages {
		if value := pkg.LookupExpr(name); value != nil {
			return value
		}
	}

	return nil
}

// LookupTypesByType returns the package-level objects of all loaded packages
// that derive from the named type obj: constants of type obj and types
// declared in terms of obj, such as `type X obj` or `type X obj[T]`.
// It returns nil if obj's name does not resolve to an identifier via
// [Contemplate.LookupExpr].
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
					if objectExprIdent := assume.T[*ast.Ident](objectExpr); objectExprIdent != nil && objectExprIdent.Obj != nil {
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

// LookupAstIdentDefsByDeclType returns the package-level objects of the
// imported packages whose type is identical to input.Type.
func (s *Contemplate) LookupAstIdentDefsByDeclType(input types.TypeAndValue) []types.Object {
	var (
		pkgs       []*packages.Package
		addImports func(pkg *packages.Package)
	)

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

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (s *Contemplate) addPackages(pkgs ...*packages.Package) {
	for _, pkg := range pkgs {
		if _, ok := s.Packages[pkg.PkgPath]; !ok {
			s.Packages[pkg.PkgPath] = NewPackage(s, pkg)

			s.addPackages(slices.AppendSeq(make([]*packages.Package, 0, len(pkg.Imports)), maps.Values(pkg.Imports))...)
		}
	}
}

func (s *Contemplate) addPackagesConfigs(confs ...*PackageConfig) {
	for _, conf := range confs {
		s.Package(conf.Path).AddScopeTypes(conf.Types...)
	}
}

func (s *Contemplate) addPackageTypeNames(pkg *packages.Package, typeNames ...string) {
	if _, ok := s.Packages[pkg.PkgPath]; !ok {
		s.Packages[pkg.PkgPath] = NewPackage(s, pkg)
	}
	// add request scopes
	s.Packages[pkg.PkgPath].AddScopeTypes(typeNames...)
}
