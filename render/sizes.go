package render

import "unsafe"

var SizeFloat32 = int(unsafe.Sizeof(float32(1.0)))
