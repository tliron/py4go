package python

import (
	"unsafe"
)

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

var ModuleType = NewType(&C.PyModule_Type)

func CreateModule(name string) (*Reference, error) {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	var definition C.PyModuleDef
	definition.m_name = name_

	if module := C.PyModule_Create2(&definition, C.PYTHON_ABI_VERSION); module != nil {
		return NewReference(module), nil
	} else {
		return nil, GetError()
	}
}

func NewModuleRaw(name string) (*Reference, error) {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	if module := C.PyModule_New(name_); module != nil {
		return NewReference(module), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) GetModuleName() (string, error) {
	if name := C.PyModule_GetName(self.Object); name != nil {
		defer C.free(unsafe.Pointer(name)) // TODO: need this?

		return C.GoString(name), nil
	} else {
		return "", GetError()
	}
}

func (self *Reference) EnableModule() error {
	if name := C.PyModule_GetNameObject(self.Object); name != nil {
		name_ := NewReference(name)
		defer name_.Release()

		if moduleDict := C.PyImport_GetModuleDict(); moduleDict != nil {
			moduleDict_ := NewReference(moduleDict)
			defer moduleDict_.Release()

			return moduleDict_.SetDictItem(name_, self)
		} else {
			return GetError()
		}
	} else {
		return GetError()
	}
}

// PyCFunction signature, second argument unused
func (self *Reference) AddModuleCFunctionNoArgs(name string, function unsafe.Pointer) error {
	return self.addModuleCFunction(name, function, C.METH_NOARGS)
}

// PyCFunction signature, second argument is the Python argument
func (self *Reference) AddModuleCFunctionOneArg(name string, function unsafe.Pointer) error {
	return self.addModuleCFunction(name, function, C.METH_O)
}

// PyCFunction signature, second argument is tuple of Python arguments
func (self *Reference) AddModuleCFunctionArgs(name string, function unsafe.Pointer) error {
	return self.addModuleCFunction(name, function, C.METH_VARARGS)
}

// PyCFunctionWithKeywords signature
func (self *Reference) AddModuleCFunctionArgsAndKeywords(name string, function unsafe.Pointer) error {
	return self.addModuleCFunction(name, function, C.METH_VARARGS|C.METH_KEYWORDS)
}

// _PyCFunctionFast signature
func (self *Reference) AddModuleCFunctionFastArgs(name string, function unsafe.Pointer) error {
	return self.addModuleCFunction(name, function, C.METH_FASTCALL)
}

// _PyCFunctionFastWithKeywords signature
func (self *Reference) AddModuleCFunctionFastArgsAndKeywords(name string, function unsafe.Pointer) error {
	return self.addModuleCFunction(name, function, C.METH_FASTCALL|C.METH_KEYWORDS)
}

func (self *Reference) addModuleCFunction(name string, function unsafe.Pointer, flags C.int) error {
	name_ := C.CString(name)
	defer C.free(unsafe.Pointer(name_))

	methodDef := []C.PyMethodDef{
		{
			ml_name:  name_,
			ml_meth:  C.PyCFunction(function),
			ml_flags: flags,
		},
		{}, // NULL end of array
	}

	if C.PyModule_AddFunctions(self.Object, &methodDef[0]) == 0 {
		return nil
	} else {
		return GetError()
	}
}
