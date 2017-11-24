package main

import (
	"math/rand"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/core"
	"github.com/stojg/graphics/lib/lights"
	"github.com/stojg/graphics/lib/rendering"
)

func main() {
	rand.Seed(19)
	l := NewLogger("gl.log")
	err := run()
	if err != nil {
		l.Error(err)
		l.Close()
		os.Exit(1)
	}
	l.Close()
}

func run() error {
	width := 800
	height := 600

	engine, err := core.NewEngine(width, height, "graphics")
	if err != nil {
		return err
	}

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(mgl32.Vec3{0, 3, 6})
	cameraObject.Transform().LookAt(mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 1, 0})
	cameraObject.AddComponent(components.NewCamera(70, width, height, 0.01, 100))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	engine.AddObject(cameraObject)

	whiteMaterial := rendering.NewMaterial()
	whiteMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/white.png"))

	tealMaterial := rendering.NewMaterial()
	tealMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/teal.png"))

	floor := core.NewGameObject()
	floor.Transform().SetScale(mgl32.Vec3{100, 0.01, 100})
	floor.Transform().SetPos(mgl32.Vec3{0, -0.005, 0})
	core.LoadModel(floor, "res/meshes/cube/model.obj", whiteMaterial)
	engine.AddObject(floor)

	dirLight := core.NewGameObject()
	dirLight.Transform().SetPos(mgl32.Vec3{2, 10, 1})
	dirLight.Transform().SetScale(mgl32.Vec3{0.5, 0.1, 0.5})
	dirLight.AddComponent(lights.NewDirectional(0.98, 0.98, 0.98, 1))
	//core.LoadModel(dirLight, "res/meshes/cube/model.obj", tealMaterial)
	engine.AddObject(dirLight)

	pointLight := core.NewGameObject()
	pointLight.Transform().SetPos(mgl32.Vec3{2, 0.5, 0})
	pointLight.Transform().SetScale(mgl32.Vec3{0.05, 0.05, 0.05})
	pointLight.AddComponent(lights.NewPoint(1, 0, 0, 5))
	core.LoadModel(pointLight, "res/meshes/cube/model.obj", tealMaterial)
	engine.AddObject(pointLight)

	spot := core.NewGameObject()
	spot.Transform().SetPos(mgl32.Vec3{-6, 4, 4})
	spot.Transform().SetScale(mgl32.Vec3{0.05, 0.05, 0.3})
	spot.Transform().LookAt(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	spot.AddComponent(lights.NewSpot(0.8, 0.9, 1, 5, 22))
	core.LoadModel(spot, "res/meshes/cube/model.obj", tealMaterial)
	engine.AddObject(spot)

	bot := core.NewGameObject()
	bot.Transform().SetPos(mgl32.Vec3{0, 0, 0})
	bot.AddComponent(components.NewRotator(mgl32.Vec3{0, -1, 0}, 23))
	core.LoadModel(bot, "res/meshes/sphere_bot/model.obj", whiteMaterial)
	engine.AddObject(bot)

	cube := core.NewGameObject()
	cube.Transform().SetPos(mgl32.Vec3{4, 1, 0})
	core.LoadModel(cube, "res/meshes/cube/model.obj", whiteMaterial)
	engine.AddObject(cube)

	{
		cube := core.NewGameObject()
		cube.Transform().SetPos(mgl32.Vec3{4, 3, 0})
		cube.Transform().SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
		core.LoadModel(cube, "res/meshes/cube/model.obj", whiteMaterial)
		engine.AddObject(cube)
	}

	engine.Start()

	return nil
}
