package buffers

import (
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
)

/* #nosec */
const (
	SizeOfUint32  = unsafe.Sizeof(uint32(0))
	SizeOfFloat32 = int(unsafe.Sizeof(float32(1)))
	SizeOfMat4    = int(unsafe.Sizeof(mgl32.Ident4()))
)
