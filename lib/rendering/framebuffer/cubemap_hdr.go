package framebuffer

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func NewHDRCubeMap(width, height int32, filename string) *HDRCubeMap {
	texture := &HDRCubeMap{
		attachment: gl.COLOR_ATTACHMENT0,
		width:      width,
		height:     height,
	}
	texture.init()
	texture.LoadHDR(filename)

	return texture
}

type HDRCubeMap struct {
	id         uint32
	attachment uint32
	width      int32
	height     int32

	fbo uint32
	rbo uint32
}

func (t *HDRCubeMap) Height() int32 {
	return t.height
}

func (t *HDRCubeMap) LoadHDR(filename string) {

	if t.width == 0 || t.height == 0 {
		panic("texture cannot have zero height or width")
	}
	for i := 0; i < 6; i++ {
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.RGB16F, t.width, t.height, 0, gl.RGB, gl.FLOAT, nil)
	}

	hdrTexture, err := loadHDRTextureResource(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	captureProjection := mgl32.Perspective(float32((90*math.Pi)/180.0), 1, 0.1, 10)
	captureViews := []mgl32.Mat4{
		mgl32.LookAt(0, 0, 0, 1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, -1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 1, 0, 0, 0, 1),
		mgl32.LookAt(0, 0, 0, 0, -1, 0, 0, 0, -1),
		mgl32.LookAt(0, 0, 0, 0, 0, 1, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 0, -1, 0, -1, 0),
	}

	shad := shader.NewShader("equirectangular_to_cubemap")
	shad.Bind()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, hdrTexture.id)

	shad.UpdateUniform("equirectangularMap", int32(0))
	shad.UpdateUniform("projection", captureProjection)

	gl.Viewport(0, 0, t.width, t.height)
	gl.Disable(gl.CULL_FACE)
	for i := 0; i < 6; i++ {
		shad.UpdateUniform("view", captureViews[i])
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, t.attachment, gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), t.id, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		primitives.DrawCube()
	}
	gl.Enable(gl.CULL_FACE)

	debug.CheckForError("framebuffer.HDRCubeMap end")
	debug.FramebufferComplete("framebuffer.HDRCubeMap")

	hdrTexture.Delete()
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (t *HDRCubeMap) Width() int32 {
	return t.width
}

func (t *HDRCubeMap) ID() uint32 {
	return t.id
}

func (t *HDRCubeMap) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, t.id)
}

func (t *HDRCubeMap) BindAsRenderTarget() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
}

func (t *HDRCubeMap) SetViewPort() {
	gl.Viewport(0, 0, t.width, t.height)
}

func (t *HDRCubeMap) init() {
	gl.GenFramebuffers(1, &t.fbo)
	gl.GenRenderbuffers(1, &t.rbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, t.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, t.rbo)

	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, t.width, t.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, t.rbo)

	gl.GenTextures(1, &t.id)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, t.id)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, t.attachment, gl.TEXTURE_CUBE_MAP_POSITIVE_X, t.id, 0)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BASE_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAX_LEVEL, 0)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
}
