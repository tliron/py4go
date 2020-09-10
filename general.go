package python

// See:
//   https://docs.python.org/3/c-api/init.html

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

func Initialize() {
	C.Py_Initialize()
}

func Finalize() error {
	if C.Py_FinalizeEx() == 0 {
		return nil
	} else {
		return GetError()
	}
}

func Version() string {
	return C.GoString(C.Py_GetVersion())
}
