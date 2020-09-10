package api

// Here we add our "py_" functions to the module

// Note that this file could potentially be combined with "py_.go", but we preferred separation
// For that reason we must forward-declare the "py_" functions in the cgo preamble

import (
	"github.com/tliron/py4go"
)

/*
#cgo pkg-config: python3-embed

#define PY_SSIZE_T_CLEAN
#include <Python.h>

PyObject *py_api_sayGoodbye(PyObject *self, PyObject *unused);
PyObject *py_api_concat(PyObject *self, PyObject *args);
PyObject *py_api_concat_fast(PyObject *self, PyObject **args, Py_ssize_t nargs);
*/
import "C"

func CreateModule() (*python.Reference, error) {
	if module, err := python.CreateModule("api"); err == nil {
		if err := module.AddModuleCFunctionNoArgs("say_goodbye", C.py_api_sayGoodbye); err != nil {
			module.Release()
			return nil, err
		}

		if err := module.AddModuleCFunctionArgs("concat", C.py_api_concat); err != nil {
			module.Release()
			return nil, err
		}

		if err := module.AddModuleCFunctionFastArgs("concat_fast", C.py_api_concat_fast); err != nil {
			module.Release()
			return nil, err
		}

		return module, nil
	} else {
		return nil, err
	}
}
