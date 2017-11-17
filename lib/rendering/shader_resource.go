package rendering

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewShaderResource() *ShaderResource {
	s := &ShaderResource{
		Program:  gl.CreateProgram(),
		refCount: 1,
	}

	if s.Program == 0 {
		fmt.Println("Shader creation failed: Could not find valid memory location in constructor")
		os.Exit(1)
	}
	return s
}

type ShaderResource struct {
	Program  uint32
	refCount int
	//uniforms map[string]int
	//uniformNames []string
	//uniformTypes []string
}

func (r *ShaderResource) AddReference() {
	r.refCount++
}

func (r *ShaderResource) Cleanup() {
	//gl.DeleteBuffers(r.Program)
}
