package main

import (
	"reflect"
	"unsafe"
)

func getUnexportedValue(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}