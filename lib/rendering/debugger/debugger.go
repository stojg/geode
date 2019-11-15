package debugger

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/geode/lib/components"
	"github.com/stojg/geode/lib/rendering/framebuffer"
	"github.com/stojg/geode/lib/rendering/primitives"
	"github.com/stojg/geode/lib/rendering/shader"
)

const (
	perRow       = 3
	gutter int32 = 10
)

var texture components.Texture
var shaders map[string]components.Shader
var w, h int32
var panelWidth, panelHeight int32
var numPanels int

func New(width, height int, s components.RenderState) {
	w, h = int32(width), int32(height)
	panelWidth = w/perRow - gutter
	panelHeight = h/perRow - gutter
	texture = framebuffer.NewTexture(gl.COLOR_ATTACHMENT0, width, height, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, gl.NEAREST, false)
	shaders = make(map[string]components.Shader)
	shaders["shadow"] = shader.NewShader("debug_shadow")
	shaders["rgb"] = shader.NewShader("filter_null")
}

func Clear() {
	texture.BindFrameBuffer()
	clearColor := [4]uint32{0, 0, 0, 0}
	gl.ClearBufferuiv(gl.COLOR, 0, &clearColor[0])
	numPanels = 0
}

func AddRGBTexture(in, out components.Texture) {
	posX, posY := getNextSlot()

	gl.Disable(gl.DEPTH_TEST)

	if out == nil {
		//gl.Viewport(0, 0, int32(w), int32(h))
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	} else {
		out.BindFrameBuffer()
	}

	gl.Viewport(posX, posY, panelWidth, panelHeight)

	rgbShader := shaders["rgb"]
	in.Activate(0)
	rgbShader.Bind()
	primitives.DrawQuad()
	if out != nil {
		out.UnbindFrameBuffer()
	}

	gl.Enable(gl.DEPTH_TEST)
}

func AddShadowTexture(in, out components.Texture) {
	posX, posY := getNextSlot()

	gl.Disable(gl.DEPTH_TEST)

	if out == nil {
		gl.Viewport(0, 0, int32(w), int32(h))
		gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	} else {
		out.BindFrameBuffer()
	}

	gl.Viewport(posX, posY, panelWidth, panelHeight)

	rgbShader := shaders["shadow"]
	in.Activate(0)
	rgbShader.Bind()
	primitives.DrawQuad()
	if out != nil {
		out.UnbindFrameBuffer()
	}

	gl.Enable(gl.DEPTH_TEST)
}

func Texture() components.Texture {
	return texture
}

func getNextSlot() (int32, int32) {
	x := int32(numPanels % perRow)
	y := int32(numPanels / perRow)

	posX := x*(panelWidth+gutter) + gutter/2
	posY := y*(panelHeight+gutter) + gutter/2
	posY = h - posY - panelHeight

	numPanels++

	var clearSlots = []int{1, 3, 4, 5, 7}
	for _, s := range clearSlots {
		if numPanels == s {
			numPanels++
		}
	}

	return posX, posY
}
