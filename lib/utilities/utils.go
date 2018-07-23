package utilities

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

func CreateUint32BO(vao uint32, target, usage uint32, indices []uint32) uint32 {
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindVertexArray(vao)
	gl.BindBuffer(target, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*int(SizeOfUint32), gl.Ptr(indices), usage)
	gl.BindVertexArray(0)
	return ebo
}

func CreateEmptyVBO(floatCount int, usage uint32) uint32 {
	var bufferObject uint32
	gl.GenBuffers(1, &bufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, bufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, floatCount*SizeOfFloat32, nil, usage)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return bufferObject
}

func AddInstancedAttribute(vao, vbo uint32, attribute uint32, dataSizeInFloats int32, instanceDataLength int, offset int) {
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(attribute, dataSizeInFloats, gl.FLOAT, false, int32(instanceDataLength*SizeOfFloat32), gl.PtrOffset(offset*SizeOfFloat32))
	gl.VertexAttribDivisor(attribute, 1)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func AddAttribute(vao, vbo uint32, index uint32, sizeInFloats int32, instanceDataLength int, offset int) {
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(index, sizeInFloats, gl.FLOAT, false, int32(instanceDataLength*SizeOfFloat32), gl.PtrOffset(offset*SizeOfFloat32))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func UpdateVBO(vbo uint32, floatCount int, data interface{}, usage uint32) {
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// Buffer orphaning, a common way to improve streaming perf.
	gl.BufferData(gl.ARRAY_BUFFER, floatCount*SizeOfFloat32, nil, usage)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, floatCount*SizeOfFloat32, gl.Ptr(data))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

}
