package contemplate

import (
	"go/ast"
	"go/types"

	"github.com/foomo/gocontemplate/pkg/assume"
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

func (s *Package) Raw() *packages.Package {
	return s.pkg
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
