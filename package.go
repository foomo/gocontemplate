package gocontemplate

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

type Package struct {
	l          *Contemplate
	pkg        *packages.Package
	exprs      map[string]ast.Expr
	types      map[string]types.Object
	scopeExprs map[string]ast.Expr
	scopeTypes map[string]types.Object
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewPackage(l *Contemplate, pkg *packages.Package) *Package {
	exprs := make(map[string]ast.Expr)
	for expr, value := range pkg.TypesInfo.Defs {
		if value != nil {
			switch value.(type) {
			case *types.Const:
				exprs[value.Name()] = expr
			case *types.Func, *types.TypeName:
				exprs[value.Name()] = expr
			}
		}
	}

	typess := make(map[string]types.Object)
	for _, name := range pkg.Types.Scope().Names() {
		typess[name] = pkg.Types.Scope().Lookup(name)
	}

	inst := &Package{
		l:          l,
		pkg:        pkg,
		types:      typess,
		exprs:      map[string]ast.Expr{},
		scopeExprs: map[string]ast.Expr{},
		scopeTypes: map[string]types.Object{},
	}

	inst.addExprs(pkg.TypesInfo.Defs)

	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Getter
// ------------------------------------------------------------------------------------------------

func (s *Package) Name() string {
	return s.pkg.Name
}

func (s *Package) Path() string {
	return s.pkg.PkgPath
}

func (s *Package) Exprs() map[string]ast.Expr {
	return s.exprs
}

func (s *Package) Types() map[string]types.Object {
	return s.types
}

func (s *Package) ScopeTypes() map[string]types.Object {
	return s.scopeTypes
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (s *Package) AddScopeTypes(names ...string) {
	for _, name := range names {
		if _, ok := s.scopeTypes[name]; !ok {
			scopeType := s.LookupType(name)
			scopeExpr := s.LookupExpr(name)
			if scopeType != nil && scopeExpr != nil {
				s.scopeTypes[name] = scopeType
				s.scopeExprs[name] = scopeExpr
				s.addScopeTypeAstExpr(scopeExpr)
			}
		}
	}
}

func (s *Package) LookupExpr(name string) ast.Expr {
	return s.exprs[name]
}

func (s *Package) LookupScopeExpr(name string) ast.Expr {
	return s.scopeExprs[name]
}

func (s *Package) FilterExprsByTypeExpr(expr ast.Expr) []ast.Expr {
	var ret []ast.Expr
	if exprIdent := TC[*ast.Ident](expr); exprIdent != nil {
		for _, child := range s.exprs {
			if childIdent := TC[*ast.Ident](child); childIdent != nil && childIdent.Obj != nil {
				if childDecl := TC[*ast.ValueSpec](childIdent.Obj.Decl); childDecl != nil {
					if childDeclType := TC[*ast.Ident](childDecl.Type); childDeclType != nil {
						if childDeclType.Obj == exprIdent.Obj {
							ret = append(ret, child)
						}
					}
				}
			}
		}
	}
	return ret
}

func (s *Package) LookupType(name string) types.Object {
	return s.types[name]
}

func (s *Package) LookupScopeType(name string) types.Object {
	return s.scopeTypes[name]
}

func (s *Package) LookupAstIdentDef(typeName string) *ast.Ident {
	for defAstIdent, defTypeObject := range s.pkg.TypesInfo.Defs {
		if defTypeObject != nil && defTypeObject.Name() == typeName {
			return defAstIdent
		}
	}
	return nil
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (s *Package) addScopeTypeAstExpr(input ast.Expr) {
	switch t := input.(type) {
	case *ast.Ident:
		if t.Obj != nil {
			s.addScopeTypeAstObject(t.Obj.Decl)
		} else {
			s.l.addPackageTypeNames(s.pkg, t.Name)
		}
	case *ast.StructType:
		for _, field := range t.Fields.List {
			s.addScopeTypeAstExpr(field.Type)
		}
	case *ast.IndexExpr:
		s.addScopeTypeAstExpr(t.X)
		s.addScopeTypeAstExpr(t.Index)
	case *ast.SelectorExpr:
		s.addScopeTypeAstSelectorExpr(t)
	}
}

func (s *Package) addScopeTypeAstSelectorExpr(input *ast.SelectorExpr) {
	if x := TC[*ast.Ident](input.X); x != nil {
		if xPkgName := TC[*types.PkgName](s.pkg.TypesInfo.Uses[x]); xPkgName != nil {
			if selIdent := TC[*ast.Ident](input.Sel); selIdent != nil {
				for node, object := range s.pkg.TypesInfo.Implicits {
					if object == xPkgName {
						if nodeImportSepc := TC[*ast.ImportSpec](node); nodeImportSepc != nil {
							v := strings.Trim(nodeImportSepc.Path.Value, "\"")
							s.l.addPackageTypeNames(s.pkg.Imports[v], selIdent.Name)
						}
					}
				}
			}
		}
	}
}

func (s *Package) addScopeTypeAstObject(input any) {
	switch t := input.(type) {
	case *ast.TypeSpec:
		s.addScopeTypeAstExpr(t.Type)
	}
}

func (s *Package) addExprs(source map[*ast.Ident]types.Object) {
	for expr, object := range source {
		if object != nil {
			switch object.(type) {
			case *types.Func:
				s.exprs[object.Name()] = expr
			case *types.Const:
				s.exprs[object.Name()] = expr
			case *types.TypeName:
				s.exprs[object.Name()] = expr
			}
		}
	}
}
