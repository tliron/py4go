package main

import (
	"fmt"
	"sync"

	python "github.com/tliron/py4go"
	"github.com/tliron/py4go/examples/hello-world/api"
)

func version() {
	fmt.Printf("Go >> Python version:\n%s\n", python.Version())
}

func typeChecking() {
	fmt.Println("Go >> Type checking:")
	float, _ := python.NewPrimitiveReference(1.0)
	fmt.Printf("Go >> IsFloat: %t\n", float.IsFloat())
}

func callPythonFunction(module *python.Reference) {
	fmt.Println("Go >> Calling a Python function:")

	hello, _ := module.GetAttr("hello")
	defer hello.Release()

	r, _ := hello.Call("Tal")
	defer r.Release()

	r_, _ := r.ToString()
	fmt.Printf("Go >> Python function returned: %s\n", r_)
}

func callPythonMethod(module *python.Reference) {
	fmt.Println("Go >> Calling a Python method:")

	person, _ := module.GetAttr("person")
	defer person.Release()

	greet, _ := person.GetAttr("greet")
	defer greet.Release()

	greet.Call()
}

func getPythonException(module *python.Reference) {
	fmt.Println("Go >> Python exception as Go error:")

	bad, _ := module.GetAttr("bad")
	defer bad.Release()

	if _, err := bad.Call(); err != nil {
		fmt.Printf("Go >> Error message: %s\n", err)
	}
}

func callGoFromPython(module *python.Reference) {
	goodbye, _ := module.GetAttr("goodbye")
	defer goodbye.Release()
	goodbye.Call()

	sayName, _ := module.GetAttr("say_name")
	defer sayName.Release()
	sayName.Call()

	sayNameFast, _ := module.GetAttr("say_name_fast")
	defer sayNameFast.Release()
	sayNameFast.Call()
}

func concurrency(module *python.Reference) {
	fmt.Println("Go >> Concurrency:")

	grow, _ := module.GetAttr("grow")
	defer grow.Release()

	func() {
		// Release Python's lock on our main thread, allowing other threads to execute
		// (Without this our calls to "grow", from other threads, will block forever)
		threadState := python.SaveThreadState()
		defer threadState.Restore()

		// Parallel work:

		var waitGroup sync.WaitGroup
		defer waitGroup.Wait()

		for i := 0; i < 5; i++ {
			waitGroup.Add(1)

			go func() {
				defer waitGroup.Done()

				for i := 0; i < 100; i++ {
					// We must manually acquire Python's Global Interpreter Lock (GIL),
					// because Python doesn't know about our Go "threads" (goroutines)
					gs := python.EnsureGilState()
					defer gs.Release()

					// (Note: We could also have acquired the GIL outside of this for-loop;
					// it's up to how we want to balance concurrency with the cost of context
					// switching)

					grow.Call(1)
				}
			}()
		}
	}()

	size, _ := module.GetAttr("size")
	defer size.Release()

	size_, _ := size.ToInt64()
	fmt.Printf("Go >> Size is %d\n", size_)
}

func main() {
	python.PrependPythonPath(".")

	python.Initialize()
	defer python.Finalize()

	version()
	fmt.Println()

	typeChecking()
	fmt.Println()

	api, _ := api.CreateModule()
	defer api.Release()
	api.EnableModule()

	foo, _ := python.Import("foo")
	defer foo.Release()

	callPythonFunction(foo)
	fmt.Println()

	callPythonMethod(foo)
	fmt.Println()

	getPythonException(foo)
	fmt.Println()

	callGoFromPython(foo)
	fmt.Println()

	concurrency(foo)
}
