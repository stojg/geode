package shader

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
)

func NewResource() *Resource {
	s := &Resource{
		Program:  gl.CreateProgram(),
		refCount: 1,
		uniforms: make(map[string]int32),
	}

	if s.Program == 0 {
		fmt.Println("Shader creation failed: Could not find valid memory location in constructor")
		os.Exit(1)
	}
	return s
}

type Resource struct {
	Program      uint32
	refCount     int
	uniforms     map[string]int32
	uniformNames []string
	uniformTypes []string
}

func (r *Resource) AddReference() {
	r.refCount++
}

func (r *Resource) Cleanup() {
	gl.DeleteBuffers(1, &r.Program)
}

func (r *Resource) UniformExists(name string) bool {
	for _, n := range r.uniformNames {
		if name == n {
			return true
		}
	}
	return false
}
func (r *Resource) AddUniformName(name string) {
	r.uniformNames = append(r.uniformNames, name)
}

func (r *Resource) AdduniformType(t string) {
	r.uniformTypes = append(r.uniformTypes, t)
}
