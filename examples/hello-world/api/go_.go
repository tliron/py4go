package api

// Here we export cgo wrappers for our plain Go functions
// They handle the conversion between Go and C types

// Note: cgo exports cannot be in the same file as cgo preamble funtions,
// which is why this file cannot be combined with "py_.go"

import "C"

//export go_api_sayGoodbye
func go_api_sayGoodbye() {
	// Note that we could have just exported sayGoodbye directly because it has no arguments,
	// and thus nothing to convert, But for completion we are adding this straightforward
	// wrapper
	sayGoodbye()
}

//export go_api_concat
func go_api_concat(a *C.char, b *C.char) *C.char {
	return C.CString(concat(C.GoString(a), C.GoString(b)))
}
