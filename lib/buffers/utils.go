package buffers

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/debug"
)

func AddInstancedAttribute(vao, vbo uint32, attribute uint32, dataSizeInFloats int32, instanceDataLength int, offset int) {
	gl.BindVertexArray(vao)
	debug.AddVertexBind()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(attribute, dataSizeInFloats, gl.FLOAT, false, int32(instanceDataLength*SizeOfFloat32), gl.PtrOffset(offset*SizeOfFloat32))
	gl.VertexAttribDivisor(attribute, 1)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
}

func AddAttribute(vao, vbo uint32, index uint32, sizeInFloats int32, strideInFloats int32, offset int) {
	gl.BindVertexArray(vao)
	debug.AddVertexBind()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(index, sizeInFloats, gl.FLOAT, false, strideInFloats*int32(SizeOfFloat32), gl.PtrOffset(offset*SizeOfFloat32))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
}

func CreateIntEBO(vao uint32, intCount int, data []uint32, usage uint32) uint32 {
	gl.BindVertexArray(vao)
	debug.AddVertexBind()
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, intCount*int(SizeOfUint32), gl.Ptr(data), usage)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
	return ebo
}

func CreateFloatVBO(vao uint32, floatCount int, data interface{}, usage uint32) uint32 {
	gl.BindVertexArray(vao)
	debug.AddVertexBind()
	var bufferObject uint32
	gl.GenBuffers(1, &bufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, bufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, floatCount*SizeOfFloat32, gl.Ptr(data), usage)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
	return bufferObject
}

func CreateVBO(vao uint32, size int, data interface{}, usage uint32) uint32 {
	gl.BindVertexArray(vao)
	debug.AddVertexBind()
	var bufferObject uint32
	gl.GenBuffers(1, &bufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, bufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, size, gl.Ptr(data), usage)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
	return bufferObject
}

func CreateEmptyFloatVBO(vao uint32, floatCount int, usage uint32) uint32 {
	gl.BindVertexArray(vao)
	debug.AddVertexBind()
	var bufferObject uint32
	gl.GenBuffers(1, &bufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, bufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, floatCount*SizeOfFloat32, nil, usage)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
	return bufferObject
}

func UpdateFloatVBO(vao, vbo uint32, floatCount int, data interface{}, usage uint32) {
	gl.BindVertexArray(vao)
	debug.AddVertexBind()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// Buffer orphaning, a common way to improve streaming perf.
	gl.BufferData(gl.ARRAY_BUFFER, floatCount*SizeOfFloat32, nil, usage)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, floatCount*SizeOfFloat32, gl.Ptr(data))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
}

func UpdateVBO(vao, vbo uint32, size int, data interface{}, usage uint32) {
	gl.BindVertexArray(vao)
	debug.AddVertexBind()
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// Buffer orphaning, a common way to improve streaming perf.
	gl.BufferData(gl.ARRAY_BUFFER, size, nil, usage)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, size, gl.Ptr(data))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	debug.AddVertexBind()
}

func CreateEmptyUBO(size int) uint32 {
	var bufferObject uint32
	gl.GenBuffers(1, &bufferObject)
	gl.BindBuffer(gl.UNIFORM_BUFFER, bufferObject)
	gl.BufferData(gl.UNIFORM_BUFFER, size, nil, gl.STATIC_DRAW)
	gl.BindBuffer(gl.UNIFORM_BUFFER, 0)
	return bufferObject

}
