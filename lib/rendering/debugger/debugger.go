package debugger

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
	"github.com/stojg/graphics/lib/rendering/shader"
)

const (
	perRow       = 4
	gutter int32 = 10
)

var texture components.Texture
var s components.Shader
var w, h int32
var panelWidth, panelHeight int32

var numPanels int

func New(width, height int) {
	w, h = int32(width), int32(height)
	panelWidth = w/perRow - gutter
	panelHeight = h/perRow - gutter
	texture = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, gl.NEAREST, false)
	s = shader.NewShader("debug_shadow")
}

func Clear() {
	texture.BindAsRenderTarget()
	clearColor := [4]uint32{0, 0, 0, 0}
	gl.ClearBufferuiv(gl.COLOR, 0, &clearColor[0])
	numPanels = 0
}

func AddTexture(in components.Texture, applyFilter func(s components.Shader, in, out components.Texture)) {
	x := int32(numPanels % perRow)
	y := int32(numPanels / perRow)
	nextSlot()

	posX := x*(panelWidth+gutter) + gutter/2
	posY := y*(panelHeight+gutter) + gutter/2
	posY = h - posY - panelHeight

	gl.Viewport(posX, posY, panelWidth, panelHeight)

	gl.Disable(gl.DEPTH_TEST)
	applyFilter(s, in, texture)
	gl.Enable(gl.DEPTH_TEST)
}

func Texture() components.Texture {
	return texture
}

func nextSlot() {
	numPanels++

	// make sure we leave a space in the middle
	var clearSlots = []int{5, 6, 9, 10}
	for _, s := range clearSlots {
		if numPanels == s {
			numPanels++
		}
	}
}
