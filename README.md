py4go
=====

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/tliron/py4go)](https://goreportcard.com/report/github.com/tliron/py4go)

Execute Python 3 code from within your Go program.

With py4go you can also expose Go functions to be called from that Python code.

To get started try running [`scripts/example`](scripts/example/). Note that you need the Python
development libraries. E.g. in Fedora:

    sudo dnf install python3-devel


Example Usage
-------------

```go
package main

import (
    "fmt"
    "github.com/tliron/py4go"
)

func main() {
    // Initialize Python
    python.Initialize()
    defer python.Finalize()

    // Import Python code (foo.py)
    foo, _ := python.Import("foo")
    defer foo.Release()

    // Get access to a Python function
    hello, _ := foo.GetAttr("hello")
    defer hello.Release()

    // Call the function with arguments
    r, _ := hello.Call("myargument")
    defer r.Release()
    fmt.Printf("Returned: %s\n", r.String())

    // Expose a Go function to Python via a C wrapper
    // (Just use "import api" from Python)
    api, _ := python.CreateModule("api")
    defer api.Release()
    api.AddModuleCFunctionNoArgs("my_function", C.api_my_function)
    api.EnableModule()
}
```

Calling Python code from Go is easy because Python is a dynamic language and CPython is an
interpreted runtime. Exposing Go code to Python is more involved as it requires writing wrapper
functions in C, which we omitted in the example above. See the [examples](examples/) directory for
more detail.


Caveats
-------

This is *not* an implementation of Python in Go. Rather, py4go works by embedding CPython into your
Go program using [cgo](https://github.com/golang/go/wiki/cgo) functionality. The advantage of this
approach is that you are using the standard Python runtime and can thus make use of the entire
ecosystem of Python libraries, including wrappers for C libraries. But there are several issues to
be aware of:

* It's more difficult to distribute your Go program because you *must* have the CPython library
  available on the target operating system with a specific name. Because different operating systems
  have their own conventions for naming this library, to create a truly portable distribution it may
  be best to distribute your program as a packaged container, e.g. using Flatpak or Docker.
* It is similarly more difficult to *build* your Go program. We are using `pkg-config: python3-embed` to
  locate the CPython SDK, which works on Fedora-based operating systems. But, because where you
  *build* will determine the requirements for where you will *run*, it may be best to build on
  Fedora, either directly or in a virtual machine or container. Unfortunately cgo does not let us
  parameterize that `pkg-config` directive, thus you will have to modify our source files in order to
  build on/for other operating systems.
* Calling functions and passing data between these two high-level language's runtime environments
  obviously incurs some overhead. Notably strings are sometimes copied multiple times internally,
  and may be encoded and decoded (Go normally uses UTF-8, Python defaults to UCS4). If you are
  frequently calling back and forth be aware of possible performance degradation. As always, if you
  experience a problem measure first and identify the bottleneck before prematurely optimizing!
* Similarly, be aware that you are simultaneously running two memory management runtimes, each with
  its own heap allocation and garbage collection threads, and that Go is unaware of Python's. Your
  Go code will thus need to explicitly call `Release` on all Python references to ensure that they are
  garbage collected. Luckily, the `defer` keyword makes this easy enough in many circumstances.
* Concurrency is a bit tricky in Python due to its infamous Global Interpreter Lock (GIL). If
  you are calling Python code from a Goroutine make sure to call `python.SaveThreadState` and
  `python.EnsureGilState` as appropriate. See the examples for more detail.


References
----------

* [go-python](https://github.com/sbinet/go-python) is a similar and more mature project for Python
  2.
* [goPy](https://github.com/qur/gopy) is a much older project for Python 2.
* [gopy](https://github.com/go-python/gopy) generates Python wrappers for Go functions.
* [setuptools-golang](https://github.com/asottile/setuptools-golang) allows you to include Go
  libraries in Python packages.
