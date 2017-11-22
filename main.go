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
	cameraObject.Transform().SetPos(mgl32.Vec3{8, 5, 8})
	cameraObject.Transform().LookAt(mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 1, 0})
	cameraObject.AddComponent(components.NewCamera(70, width, height, 0.01, 100))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	engine.AddObject(cameraObject)

	whiteMaterial := rendering.NewMaterial()
	whiteMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/white.png"))

	floor := core.NewGameObject()
	floor.Transform().SetScale(mgl32.Vec3{100, 0.01, 100})
	floor.Transform().SetPos(mgl32.Vec3{0, -0.005, 0})
	core.LoadModel(floor, "res/meshes/cube/model.obj", whiteMaterial)
	engine.AddObject(floor)

	{
		dirLight := core.NewGameObject()
		dirLight.Transform().SetPos(mgl32.Vec3{2, 6, -1})
		dirLight.AddComponent(lights.NewDirectional(0.99, 0.98, 0.7, 1))
		engine.AddObject(dirLight)
	}

	bot := core.NewGameObject()
	bot.Transform().SetPos(mgl32.Vec3{0, 0, 0})
	bot.AddComponent(components.NewRotator(mgl32.Vec3{0, -1, 0}, 23))
	core.LoadModel(bot, "res/meshes/sphere_bot/model.obj", whiteMaterial)
	engine.AddObject(bot)

	engine.Start()

	return nil
}
