package framebuffer

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/debug"
	"github.com/stojg/graphics/lib/images"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func NewCubeMap(width, height int32, mipMap bool) *CubeMap {

	if width == 0 || height == 0 {
		panic("texture cannot have zero height or width")
	}

	texture := &CubeMap{
		attachment: gl.COLOR_ATTACHMENT0,
		width:      width,
		height:     height,
	}
	texture.init(mipMap)
	texture.reserve()

	if mipMap {
		gl.GenerateMipmap(gl.TEXTURE_CUBE_MAP)
	}

	debug.CheckForError("framebuffer.CubeMap end")
	debug.FramebufferComplete("framebuffer.LDRCubeMap")

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return texture
}

type CubeMap struct {
	id         uint32
	attachment uint32
	width      int32
	height     int32

	fbo uint32
	rbo uint32
}

func (cubeMap *CubeMap) Height() int32 {
	return cubeMap.height
}

func (cubeMap *CubeMap) Width() int32 {
	return cubeMap.width
}

func (cubeMap *CubeMap) ID() uint32 {
	return cubeMap.id
}

func (cubeMap *CubeMap) Activate(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubeMap.id)
}

func (cubeMap *CubeMap) BindFrameBuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, cubeMap.fbo)
	gl.Viewport(0, 0, cubeMap.width, cubeMap.height)
}

func (cubeMap *CubeMap) UnbindFrameBuffer() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (cubeMap *CubeMap) init(mipMap bool) {
	gl.GenFramebuffers(1, &cubeMap.fbo)
	gl.GenRenderbuffers(1, &cubeMap.rbo)

	gl.BindFramebuffer(gl.FRAMEBUFFER, cubeMap.fbo)
	gl.BindRenderbuffer(gl.RENDERBUFFER, cubeMap.rbo)

	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT24, cubeMap.width, cubeMap.height)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, cubeMap.rbo)

	gl.GenTextures(1, &cubeMap.id)

	gl.BindTexture(gl.TEXTURE_CUBE_MAP, cubeMap.id)
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, cubeMap.attachment, gl.TEXTURE_CUBE_MAP_POSITIVE_X, cubeMap.id, 0)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	if mipMap {
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_BASE_LEVEL, 0)
		gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAX_LEVEL, 0)
	}

	if debug.CheckForError("initCubeMap") {
		panic("init cubemap")
	}
}

func (cubeMap *CubeMap) reserve() {
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X, cubeMap.id, 0)

	for i := 0; i < 6; i++ {
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.RGB16F, cubeMap.width, cubeMap.height, 0, gl.RGB, gl.FLOAT, nil)
	}

}

func (cubeMap *CubeMap) loadEquiRectangular(filename string) {

	if cubeMap.width == 0 || cubeMap.height == 0 {
		panic("texture cannot have zero height or width")
	}
	for i := 0; i < 6; i++ {
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.RGB16F, cubeMap.width, cubeMap.height, 0, gl.RGB, gl.FLOAT, nil)
	}

	width, height, data, err := images.RGBEImagedata(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	images.FlipRaw(width, height, data)
	hdrTexture, err := createTextureResource(width, height, gl.RGB32F, gl.FLOAT, data), nil
	if err != nil {
		fmt.Println(err)
		return
	}

	captureProjection := CubeProjection()
	captureViews := CubeViews()

	shad := shader.NewShader("equirectangular_to_cubemap")
	shad.Bind()

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, hdrTexture.id)

	shad.UpdateUniform("equirectangularMap", int32(0))
	shad.UpdateUniform("projection", captureProjection)

	gl.Viewport(0, 0, cubeMap.width, cubeMap.height)
	gl.Disable(gl.CULL_FACE)
	for i := 0; i < 6; i++ {
		shad.UpdateUniform("view", captureViews[i])
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, cubeMap.attachment, gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), cubeMap.id, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		primitives.DrawCube()
	}
	gl.Enable(gl.CULL_FACE)

	hdrTexture.Delete()
}

func CubeProjection() mgl32.Mat4 {
	return mgl32.Perspective(float32((90*math.Pi)/180.0), 1, 0.1, 10)
}

func CubeViews() []mgl32.Mat4 {
	return []mgl32.Mat4{
		mgl32.LookAt(0, 0, 0, 1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, -1, 0, 0, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 1, 0, 0, 0, 1),
		mgl32.LookAt(0, 0, 0, 0, -1, 0, 0, 0, -1),
		mgl32.LookAt(0, 0, 0, 0, 0, 1, 0, -1, 0),
		mgl32.LookAt(0, 0, 0, 0, 0, -1, 0, -1, 0),
	}
}

func (cubeMap *CubeMap) loadFromFiles(files [6]string) {
	img, err := images.RGBAImagedata(files[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	cubeMap.width = int32(img.Rect.Size().X)
	cubeMap.height = int32(img.Rect.Size().Y)

	for i := 0; i < 6; i++ {
		img, err := images.RGBAImagedata(files[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		if cubeMap.width == 0 || cubeMap.height == 0 {
			panic("texture cannot have zero height or width")
		}
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.SRGB, cubeMap.width, cubeMap.height, 0, gl.RGBA, gl.UNSIGNED_INT_8_8_8_8_REV, gl.Ptr(img.Pix))
	}

}
