package api

// Here we define our functions in plain Go

import (
	"fmt"
)

func sayGoodbye() {
	fmt.Println("Go >> Goodbye from Go!")
}

func concat(a string, b string) string {
	fmt.Printf("Go >> Concatenating %q and %q\n", a, b)
	return a + " " + b
}
