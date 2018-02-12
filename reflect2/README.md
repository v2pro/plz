# Reflect2 

unsafe reflect api

## convert variable to unsafe.Pointer

* if v is a value (int, struct, array), use `unsafe.Pointer(&v)`
* if v is interface{}, use `unsafe.Pointer(&v)`
* if v is non-empty interface, use `unsafe.Pointer(&v)`
* for other cases (map, slice, string, chan, func), use `relfect2.PtrOf(v)`

## convert unsafe.Pointer back to interface{}

use `type2.PackEFace(ptr)`