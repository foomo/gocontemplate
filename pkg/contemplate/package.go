package contemplate

import (
	"go/ast"
	"go/types"

	"github.com/foomo/gocontemplate/pkg/assume"
	"golang.org/x/tools/go/packages"
)

// Package wraps a loaded [packages.Package] and tracks which of its types
// are in scope. Create it with [NewPackage].
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

// NewPackage returns a new Package for pkg that belongs to l. pkg must be
// loaded with type information and syntax. No types are in scope initially.
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

// Name returns the package name.
func (s *Package) Name() string {
	return s.pkg.Name
}

// Path returns the package import path.
func (s *Package) Path() string {
	return s.pkg.PkgPath
}

// Exprs returns the defining identifiers of all package-level constants,
// functions and types, keyed by name. The map is owned by the Package and
// must not be modified.
func (s *Package) Exprs() map[string]ast.Expr {
	return s.exprs
}

// Types returns all objects in the package scope, keyed by name. The map is
// owned by the Package and must not be modified.
func (s *Package) Types() map[string]types.Object {
	return s.types
}

// ScopeTypes returns the objects that are in scope, keyed by name: the
// requested types and the types they reference. The map is owned by the
// Package and must not be modified.
func (s *Package) ScopeTypes() map[string]types.Object {
	return s.scopeTypes
}

// Raw returns the underlying [packages.Package].
func (s *Package) Raw() *packages.Package {
	return s.pkg
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

// AddScopeTypes brings the named types into scope, together with all types
// they reference, in this and in other loaded packages. Names that are
// unknown or already in scope are ignored.
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

// LookupExpr returns the defining identifier of the package-level constant,
// function or type called name, or nil if there is none.
func (s *Package) LookupExpr(name string) ast.Expr {
	return s.exprs[name]
}

// LookupScopeExpr returns the defining identifier of the in-scope type
// called name, or nil if name is not in scope.
func (s *Package) LookupScopeExpr(name string) ast.Expr {
	return s.scopeExprs[name]
}

// FilterExprsByTypeExpr returns the defining identifiers of the package-level
// constants whose declared type is the type identified by expr, such as the
// values of an enum type. It returns nil if expr is not an [*ast.Ident].
// The order of the result is unspecified.
func (s *Package) FilterExprsByTypeExpr(expr ast.Expr) []ast.Expr {
	var ret []ast.Expr

	if exprIdent := assume.T[*ast.Ident](expr); exprIdent != nil {
		for _, child := range s.exprs {
			if childIdent := assume.T[*ast.Ident](child); childIdent != nil && childIdent.Obj != nil {
				if childDecl := assume.T[*ast.ValueSpec](childIdent.Obj.Decl); childDecl != nil {
					if childDeclType := assume.T[*ast.Ident](childDecl.Type); childDeclType != nil {
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

// LookupType returns the object called name in the package scope, or nil if
// there is none.
func (s *Package) LookupType(name string) types.Object {
	return s.types[name]
}

// LookupScopeType returns the in-scope object called name, or nil if name is
// not in scope.
func (s *Package) LookupScopeType(name string) types.Object {
	return s.scopeTypes[name]
}

// LookupAstIdentDef returns an identifier in the package that defines an
// object called typeName, or nil if there is none. Local definitions are
// included; if several match, which one is returned is unspecified.
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
		}

		if t.IsExported() {
			s.l.addPackageTypeNames(s.pkg, t.Name)
		}
	case *ast.StructType:
		for _, field := range t.Fields.List {
			s.addScopeTypeAstExpr(field.Type)
		}
	case *ast.IndexExpr:
		s.addScopeTypeAstExpr(t.X)
		s.addScopeTypeAstExpr(t.Index)
	case *ast.ArrayType:
		s.addScopeTypeAstExpr(t.Elt)
	case *ast.SelectorExpr:
		s.addScopeTypeAstSelectorExpr(t)
	default:
		// fmt.Println(input, t)
	}
}

func (s *Package) addScopeTypeAstSelectorExpr(input *ast.SelectorExpr) {
	if inputTypeNamed := assume.T[*types.Named](s.pkg.TypesInfo.TypeOf(input)); inputTypeNamed != nil {
		s.l.addPackageTypeNames(s.pkg.Imports[inputTypeNamed.Obj().Pkg().Path()], inputTypeNamed.Obj().Name())
	}
}

func (s *Package) addScopeTypeAstObject(input any) {
	switch t := input.(type) {
	case *ast.TypeSpec:
		s.addScopeTypeAstFieldList(t.TypeParams)
		s.addScopeTypeAstExpr(t.Type)
	}
}

func (s *Package) addScopeTypeAstFieldList(input *ast.FieldList) {
	if input != nil {
		for _, field := range input.List {
			s.addScopeTypeAstField(field)
		}
	}
}

func (s *Package) addScopeTypeAstField(input *ast.Field) {
	// switch t := input.(type) {
	// case *ast.TypeSpec:
	// 	s.addScopeTypeAstFieldList(t.TypeParams)
	// 	s.addScopeTypeAstExpr(t.Type)
	// }
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
