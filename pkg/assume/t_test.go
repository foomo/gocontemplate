package assume_test

import (
	"fmt"

	"github.com/foomo/gocontemplate/pkg/assume"
)

func ExampleT() {
	var v any = "hello"

	fmt.Printf("%q\n", assume.T[string](v))
	fmt.Println(assume.T[int](v))
	fmt.Println(assume.T[error](nil))
	// Output:
	// "hello"
	// 0
	// <nil>
}
