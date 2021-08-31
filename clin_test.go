package clin

import (
	"fmt"
	"strings"
)


func ExampleArgs() {

	for _, s := range Args([]string{"ordinary ", " flags", ""}) {
		fmt.Println("[" + s + "]")
	}

	// Output:
	// [ordinary ]
	// [ flags]
	// []
}

func ExampleInput_Args() {

	const stdin = `
a stream

of  
	input	
  tokens
`

	in := Default()
	in.Stream = strings.NewReader(stdin)

	for _, s := range in.Args([]string{}) {
		fmt.Println("[" + s + "]")
	}

	// Output:
	// []
	// [a stream]
	// []
	// [of  ]
	// [	input	]
	// [  tokens]
}
