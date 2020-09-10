package main

import (
	"fmt"

	"github.com/tliron/py4go"
	"github.com/tliron/py4go/examples/hello-world/api"
)

func main() {
	python.PrependPythonPath(".")

	python.Initialize()
	defer python.Finalize()

	fmt.Printf("Go >> Python version:\n%s\n", python.Version())
	fmt.Println()

	fmt.Println("Go >> Type checking:")
	float, _ := python.NewPrimitiveReference(1.0)
	fmt.Printf("Go >> IsFloat: %t\n", float.IsFloat())
	fmt.Println()

	api, _ := api.CreateModule()
	api.EnableModule()
	defer api.Release()

	foo, _ := python.Import("foo")
	defer foo.Release()

	fmt.Println("Go >> Calling a Python function:")
	hello, _ := foo.GetAttr("hello")
	defer hello.Release()
	r, _ := hello.Call("Tal")
	defer r.Release()
	r_, _ := r.ToString()
	fmt.Printf("Go >> Python function returned: %s\n", r_)
	fmt.Println()

	fmt.Println("Go >> Calling a Python method:")
	person, _ := foo.GetAttr("person")
	defer person.Release()
	greet, _ := person.GetAttr("greet")
	defer greet.Release()
	greet.Call()
	fmt.Println()

	fmt.Println("Go >> Python exception as Go error:")
	bad, _ := foo.GetAttr("bad")
	defer bad.Release()
	if _, err := bad.Call(); err != nil {
		fmt.Printf("Go >> Error message: %s\n", err)
	}
	fmt.Println()

	goodbye, _ := foo.GetAttr("goodbye")
	defer goodbye.Release()
	goodbye.Call()

	sayName, _ := foo.GetAttr("say_name")
	defer sayName.Release()
	sayName.Call()

	sayNameFast, _ := foo.GetAttr("say_name_fast")
	defer sayNameFast.Release()
	sayNameFast.Call()
}
