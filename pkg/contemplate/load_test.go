package contemplate_test

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/foomo/gocontemplate/pkg/contemplate"
)

func ExampleLoad() {
	ctpl, err := contemplate.Load(context.Background(), &contemplate.Config{
		Packages: []*contemplate.PackageConfig{
			{
				Path:  "github.com/foomo/gocontemplate/test/event",
				Types: []string{"PageView"},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	for _, path := range slices.Sorted(maps.Keys(ctpl.Packages)) {
		fmt.Println(path, slices.Sorted(maps.Keys(ctpl.Package(path).ScopeTypes())))
	}
	// Output:
	// github.com/foomo/gocontemplate/test [Currency Event]
	// github.com/foomo/gocontemplate/test/event [PageView]
	// github.com/foomo/gocontemplate/test/params [PageView]
}
