package lights

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
)

type BaseLight struct {
	components.GameComponent

	color  mgl32.Vec3
	shader components.Shader
}

func (b *BaseLight) SetShader(shader components.Shader) {
	b.shader = shader
}

func (b *BaseLight) Shader() components.Shader {
	return b.shader
}

func (b *BaseLight) Color() mgl32.Vec3 {
	return b.color
}

func (b *BaseLight) SetColor(color mgl32.Vec3) {
	b.color = color
}

func (b *BaseLight) Position() mgl32.Vec3 {
	return b.Parent().Transform().Pos()
}
