package python

// See:
//   https://docs.python.org/3/c-api/concrete.html

import (
	"fmt"
	"unsafe"
)

/*
#define PY_SSIZE_T_CLEAN
#include <Python.h>
*/
import "C"

func NewPrimitiveReference(value interface{}) (*Reference, error) {
	if value == nil {
		return None, nil
	}

	switch value_ := value.(type) {
	case *Reference:
		return value_, nil
	case bool:
		if value_ {
			return True, nil
		} else {
			return False, nil
		}
	case int64:
		return NewLong(value_)
	case int32:
		return NewLong(int64(value_))
	case int8:
		return NewLong(int64(value_))
	case int:
		return NewLong(int64(value_))
	case float64:
		return NewFloat(value_)
	case float32:
		return NewFloat(float64(value_))
	case string:
		return NewUnicode(value_)
	case []interface{}:
		return NewList(value_...)
	case []byte:
		return NewBytes(value_)
	}

	return nil, fmt.Errorf("unsupported primitive: %s", value)
}

//
// None
//

var None = NewReference(C.Py_None)

//
// Bool
//

var BoolType = NewType(&C.PyBool_Type)

var True = NewReference(C.Py_True)
var False = NewReference(C.Py_False)

func (self *Reference) IsBool() bool {
	return self.Type().IsSubtype(BoolType)
}

func (self *Reference) ToBool() bool {
	switch self.Object {
	case C.Py_True:
		return true
	case C.Py_False:
		return false
	}
	return false
}

//
// Long
//

var LongType = NewType(&C.PyLong_Type)

func NewLong(value int64) (*Reference, error) {
	if long := C.PyLong_FromLong(C.long(value)); long != nil {
		return NewReference(long), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsLong() bool {
	// More efficient to use the flag
	return self.Type().HasFlag(C.Py_TPFLAGS_LONG_SUBCLASS)
}

func (self *Reference) ToInt64() (int64, error) {
	if long := C.PyLong_AsLong(self.Object); !HasException() {
		return int64(long), nil
	} else {
		return 0, GetError()
	}
}

//
// Float
//

var FloatType = NewType(&C.PyFloat_Type)

func NewFloat(value float64) (*Reference, error) {
	if float := C.PyFloat_FromDouble(C.double(value)); float != nil {
		return NewReference(float), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsFloat() bool {
	return self.Type().IsSubtype(FloatType)
}

func (self *Reference) ToFloat64() (float64, error) {
	if double := C.PyFloat_AsDouble(self.Object); !HasException() {
		return float64(double), nil
	} else {
		return 0.0, GetError()
	}
}

//
// Unicode
//

var UnicodeType = NewType(&C.PyUnicode_Type)

func NewUnicode(value string) (*Reference, error) {
	value_ := C.CString(value)
	defer C.free(unsafe.Pointer(value_))

	if unicode := C.PyUnicode_FromString(value_); unicode != nil {
		return NewReference(unicode), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsUnicode() bool {
	// More efficient to use the flag
	return self.Type().HasFlag(C.Py_TPFLAGS_UNICODE_SUBCLASS)
}

func (self *Reference) ToString() (string, error) {
	if utf8stringBytes := C.PyUnicode_AsUTF8String(self.Object); utf8stringBytes != nil {
		defer C.Py_DecRef(utf8stringBytes)

		if utf8string := C.PyBytes_AsString(utf8stringBytes); utf8string != nil {
			return C.GoString(utf8string), nil
		} else {
			return "", GetError()
		}
	} else {
		return "", GetError()
	}
}

//
// Tuple
//

var TupleType = NewType(&C.PyTuple_Type)

func NewTuple(items ...interface{}) (*Reference, error) {
	if tuple, err := NewTupleRaw(len(items)); err == nil {
		for index, item := range items {
			if item_, err := NewPrimitiveReference(item); err == nil {
				if err := tuple.SetTupleItem(index, item_); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		return tuple, nil
	} else {
		return nil, GetError()
	}
}

func NewTupleRaw(size int) (*Reference, error) {
	if tuple := C.PyTuple_New(C.long(size)); tuple != nil {
		return NewReference(tuple), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsTuple() bool {
	// More efficient to use the flag
	return self.Type().HasFlag(C.Py_TPFLAGS_TUPLE_SUBCLASS)
}

func (self *Reference) SetTupleItem(index int, item *Reference) error {
	if C.PyTuple_SetItem(self.Object, C.long(index), item.Object) == 0 {
		return nil
	} else {
		return GetError()
	}
}

//
// List
//

var ListType = NewType(&C.PyList_Type)

func NewList(items ...interface{}) (*Reference, error) {
	if list, err := NewListRaw(len(items)); err == nil {
		for index, item := range items {
			if item_, err := NewPrimitiveReference(item); err == nil {
				if err := list.SetListItem(index, item_); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		return list, nil
	} else {
		return nil, GetError()
	}
}

func NewListRaw(size int) (*Reference, error) {
	if list := C.PyList_New(C.long(size)); list != nil {
		return NewReference(list), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsList() bool {
	// More efficient to use the flag
	return self.Type().HasFlag(C.Py_TPFLAGS_LIST_SUBCLASS)
}

func (self *Reference) SetListItem(index int, item *Reference) error {
	if C.PyList_SetItem(self.Object, C.long(index), item.Object) == 0 {
		return nil
	} else {
		return GetError()
	}
}

//
// Dict
//

var DictType = NewType(&C.PyDict_Type)

func NewDict() (*Reference, error) {
	if dict := C.PyDict_New(); dict != nil {
		return NewReference(dict), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsDict() bool {
	// More efficient to use the flag
	return self.Type().HasFlag(C.Py_TPFLAGS_DICT_SUBCLASS)
}

func (self *Reference) SetDictItem(key *Reference, value *Reference) error {
	if C.PyDict_SetItem(self.Object, key.Object, value.Object) == 0 {
		return nil
	} else {
		return GetError()
	}
}

//
// Set (mutable)
//

var SetType = NewType(&C.PySet_Type)

func NewSet(iterable *Reference) (*Reference, error) {
	if set := C.PySet_New(iterable.Object); set != nil {
		return NewReference(set), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsSet() bool {
	return self.Type().IsSubtype(SetType)
}

//
// Frozen set (immutable)
//

var FrozenSetType = NewType(&C.PyFrozenSet_Type)

func NewFrozenSet(iterable *Reference) (*Reference, error) {
	if frozenSet := C.PyFrozenSet_New(iterable.Object); frozenSet != nil {
		return NewReference(frozenSet), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsFrozenSet() bool {
	return self.Type().IsSubtype(FrozenSetType)
}

//
// Bytes (immutable)
//

var BytesType = NewType(&C.PyBytes_Type)

func NewBytes(value []byte) (*Reference, error) {
	size := len(value)
	value_ := C.CBytes(value)
	defer C.free(value_) // TODO: check this!

	if bytes := C.PyBytes_FromStringAndSize((*C.char)(value_), C.long(size)); bytes != nil {
		return NewReference(bytes), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsBytes() bool {
	// More efficient to use the flag
	return self.Type().HasFlag(C.Py_TPFLAGS_BYTES_SUBCLASS)
}

func (self *Reference) ToBytes() ([]byte, error) {
	if pointer, size, err := self.AccessBytes(); err == nil {
		return C.GoBytes(pointer, C.int(size)), nil
	} else {
		return nil, err
	}
}

func (self *Reference) AccessBytes() (unsafe.Pointer, int, error) {
	if string_ := C.PyBytes_AsString(self.Object); string_ != nil {
		size := C.PyBytes_Size(self.Object)
		return unsafe.Pointer(string_), int(size), nil
	} else {
		return nil, 0, GetError()
	}
}

//
// Byte array (mutable)
//

var ByteArrayType = NewType(&C.PyByteArray_Type)

func NewByteArray(value []byte) (*Reference, error) {
	size := len(value)
	value_ := C.CBytes(value)
	defer C.free(value_) // TODO: check this!

	if byteArray := C.PyByteArray_FromStringAndSize((*C.char)(value_), C.long(size)); byteArray != nil {
		return NewReference(byteArray), nil
	} else {
		return nil, GetError()
	}
}

func (self *Reference) IsByteArray() bool {
	return self.Type().IsSubtype(ByteArrayType)
}

func (self *Reference) ByteArrayToBytes() ([]byte, error) {
	if pointer, size, err := self.AccessByteArray(); err == nil {
		return C.GoBytes(pointer, C.int(size)), nil
	} else {
		return nil, err
	}
}

func (self *Reference) AccessByteArray() (unsafe.Pointer, int, error) {
	if string_ := C.PyByteArray_AsString(self.Object); string_ != nil {
		size := C.PyByteArray_Size(self.Object)
		return unsafe.Pointer(string_), int(size), nil
	} else {
		return nil, 0, GetError()
	}
}
