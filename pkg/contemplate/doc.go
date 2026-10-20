// Package contemplate loads Go packages with [golang.org/x/tools/go/packages]
// and narrows them down to a set of requested types and the types they depend
// on. It is intended as a building block for code generators.
//
// # Usage
//
// Describe the packages and types of interest in a [Config] and call [Load]:
//
//	ctpl, err := contemplate.Load(ctx, &contemplate.Config{
//		Packages: []*contemplate.PackageConfig{
//			{Path: "example.com/event", Types: []string{"PageView"}},
//		},
//	})
//	if err != nil {
//		return err
//	}
//	for path, pkg := range ctpl.Packages {
//		fmt.Println(path, len(pkg.ScopeTypes()))
//	}
//
// [Contemplate.Packages] holds every loaded package, including transitive
// imports. Each [Package] keeps its full type scope ([Package.Types]) and the
// subset reachable from the requested types ([Package.ScopeTypes]).
package contemplate
