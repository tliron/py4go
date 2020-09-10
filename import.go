package python

// See:
//   https://docs.python.org/3/c-api/import.html

import (
	"unsafe"
)

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

func Import(name string) (*Reference, error) {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	if import_ := C.PyImport_ImportModule(name_); import_ != nil {
		return NewReference(import_), nil
	} else {
		return nil, GetError()
	}
}
