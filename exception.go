package python

// See:
//   https://docs.python.org/3/c-api/exceptions.html

import (
	"errors"
)

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

func HasException() bool {
	return C.PyErr_Occurred() != nil
}

func GetError() error {
	if exception := FetchException(); exception != nil {
		return exception
	} else {
		return errors.New("Python error without an exception")
	}
}

//
// Exception
//

type Exception struct {
	Type      *Reference
	Value     *Reference
	Traceback *Reference
}

func FetchException() *Exception {
	var type_, value, traceback *C.PyObject
	C.PyErr_Fetch(&type_, &value, &traceback)
	if type_ != nil {
		defer C.PyErr_Restore(type_, value, traceback)

		var type__, value_, traceback_ *Reference

		if type_ != nil {
			type__ = NewReference(type_)
		}

		if value != nil {
			value_ = NewReference(value)
		}

		if traceback != nil {
			traceback_ = NewReference(traceback)
		}

		return NewExceptionRaw(type__, value_, traceback_)
	} else {
		return nil
	}
}

func NewExceptionRaw(type_ *Reference, value *Reference, traceback *Reference) *Exception {
	return &Exception{
		Type:      type_,
		Value:     value,
		Traceback: traceback,
	}
}

// error signature
func (self *Exception) Error() string {
	// TODO: include traceback?
	if self.Value != nil {
		return self.Value.String()
	} else if self.Type != nil {
		return self.Type.String()
	} else {
		return "malformed Python exception"
	}
}
