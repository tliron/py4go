package python

// See:
//  https://docs.python.org/3/c-api/
//  https://github.com/golang/go/wiki/cgo
//  https://www.datadoghq.com/blog/engineering/cgo-and-python/
//  https://github.com/sbinet/go-python

// #cgo pkg-config: python3-embed
// #cgo LDFLAGS: -lpython3
import "C"
