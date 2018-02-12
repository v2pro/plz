# Reflect2 

unsafe reflect api

unsafe.Pointer is used instead of reflect.Value

## convert variable to unsafe.Pointer

* if v is a value (bool, int, float, struct, array), use `unsafe.Pointer(&v)`
* if v is interface{}, use `unsafe.Pointer(&v)`
* if v is non-empty interface, use `unsafe.Pointer(&v)`
* for other cases (map, slice, string, chan, func), use `relfect2.PtrOf(v)`

## convert unsafe.Pointer back to interface{}

use `type2.PackEFace(ptr)`