package api

// Here we define Python wrappers in C for our "go_" functions
// They handle the conversion between Python and C types

// We are demonstrating the two supported Python ABIs:
// 1) the more common PyCFunction ABI, which passes arguments as Python tuples and dicts, and
// 2) the newer _PyCFunctionFast ABI, which passes arguments as a more efficient stack

// Unfortunately we must use C and not Go because:
// 1) the "PyArg_Parse_" functions all use variadic arguments, which are not supported by cgo, and
// 2) the "PyArg_Parse_" functions unpack arguments to pointers, which we cannot implement in Go

// Note: cgo exports cannot be in the same file as cgo preamble functions,
// which is why this file cannot be combined with "go_.go"
// and is also why must forward-declare the "go_" functions in the cgo preamble

// See:
//   https://docs.python.org/3/c-api/arg.html

/*
#cgo pkg-config: python3-embed

#define PY_SSIZE_T_CLEAN
#include <Python.h>

void go_api_sayGoodbye();
char *go_api_concat(char*, char*);

// PyCFunction signature
PyObject *py_api_sayGoodbye(PyObject *self, PyObject *unused) {
	go_api_sayGoodbye();
	return Py_None;
}

// PyCFunction signature
PyObject *py_api_concat(PyObject *self, PyObject *args) {
	char *arg1 = NULL, *arg2 = NULL;
	PyArg_ParseTuple(args, "ss", &arg1, &arg2);
	char *r = go_api_concat(arg1, arg2);
	return PyUnicode_FromString(r);
}

// _PyCFunctionFast signature
PyObject *py_api_concat_fast(PyObject *self, PyObject **args, Py_ssize_t nargs) {
	char *arg1 = NULL, *arg2 = NULL;
	_PyArg_ParseStack(args, nargs, "ss", &arg1, &arg2);
	char *r = go_api_concat(arg1, arg2);
	return PyUnicode_FromString(r);
}
*/
import "C"
