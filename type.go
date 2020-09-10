package python

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

//
// Type
//

type Type struct {
	Object *C.PyTypeObject
}

func NewType(pyTypeObject *C.PyTypeObject) *Type {
	return &Type{pyTypeObject}
}

func (self *Type) IsSubtype(type_ *Type) bool {
	return C.PyType_IsSubtype(self.Object, type_.Object) != 0
}

func (self *Type) HasFlag(flag C.ulong) bool {
	return C.PyType_GetFlags(self.Object)&flag != 0
}
