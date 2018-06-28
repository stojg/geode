package framebuffer

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/loaders"
	"github.com/stojg/graphics/lib/rendering/primitives"
	"github.com/stojg/graphics/lib/rendering/shader"
)

func reserveCubeMap(texture *CubeMap) {
	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X, texture.id, 0)

	if texture.width == 0 || texture.height == 0 {
		panic("texture cannot have zero height or width")
	}
	for i := 0; i < 6; i++ {
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.RGB16F, texture.width, texture.height, 0, gl.RGB, gl.FLOAT, nil)
	}

}

func loadEquirectangular(t *CubeMap, filename string) {

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

	hdrTexture.Delete()
}

func loadFromFiles(t *CubeMap, files [6]string) {
	img, err := loaders.RGBAImagedata(files[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	t.width = int32(img.Rect.Size().X)
	t.height = int32(img.Rect.Size().Y)

	for i := 0; i < 6; i++ {
		img, err := loaders.RGBAImagedata(files[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		if t.width == 0 || t.height == 0 {
			panic("texture cannot have zero height or width")
		}
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i), 0, gl.SRGB, t.width, t.height, 0, gl.RGBA, gl.UNSIGNED_INT_8_8_8_8_REV, gl.Ptr(img.Pix))
	}

}

func loadHDRTextureResource(filename string) (*TextureResource, error) {
	width, height, data, err := loaders.RGBEImagedata(filename)
	if err != nil {
		return nil, err
	}
	loaders.FlipRaw(width, height, data)
	return createTextureResource(width, height, gl.RGB32F, gl.FLOAT, data), nil
}
