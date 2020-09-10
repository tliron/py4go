package python

// See:
//   https://docs.python.org/3/c-api/object.html

import (
	"unsafe"
)

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

// fmt.Stringer interface
func (self *Reference) String() string {
	if str := C.PyObject_Str(self.Object); str != nil {
		if data := C.PyUnicode_AsUTF8String(str); data != nil {
			defer C.Py_DecRef(data)
			if string_ := C.PyBytes_AsString(data); string_ != nil {
				return C.GoString(string_)
			} else {
				return ""
			}
		} else {
			return ""
		}
	} else {
		return ""
	}
}

func (self *Reference) Acquire() {
	C.Py_IncRef(self.Object)
}

func (self *Reference) Release() {
	C.Py_DecRef(self.Object)
}

func (self *Reference) Str() (*Reference, error) {
	if str := C.PyObject_Str(self.Object); str != nil {
		return NewReference(str), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) GetAttr(name string) (*Reference, error) {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	if attr := C.PyObject_GetAttrString(self.Object, name_); attr != nil {
		return NewReference(attr), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) SetAttr(name string, reference *Reference) error {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	if C.PyObject_SetAttrString(self.Object, name_, reference.Object) == 0 {
		return nil
	} else {
		return GetError()
	}
}

func (self *Reference) Call(args ...interface{}) (*Reference, error) {
	if args_, err := NewTuple(args...); err == nil {
		if kw, err := NewDict(); err == nil {
			return self.CallRaw(args_, kw)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *Reference) CallRaw(args *Reference, kw *Reference) (*Reference, error) {
	if r := C.PyObject_Call(self.Object, args.Object, kw.Object); r != nil {
		return NewReference(r), nil
	} else {
		return nil, GetError()
	}
}
