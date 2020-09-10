package python

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

//
// Reference
//

type Reference struct {
	Object *C.PyObject
}

func NewReference(pyObject *C.PyObject) *Reference {
	return &Reference{pyObject}
}

func (self *Reference) Type() *Type {
	return NewType(self.Object.ob_type)
}
