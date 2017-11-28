package main

import (
	"fmt"
	"math"
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
	l := newLogger("gl.log")
	err := run()
	if err != nil {
		l.error(err)
		if err := l.close(); err != nil {
			fmt.Println(".. in addition the log file had problem closing", err)
		}
		os.Exit(1)
	}
	if err := l.close(); err != nil {
		fmt.Println(".. in addition the log file had problem closing", err)
	}
}

func run() error {
	width := 800
	height := 600

	engine, err := core.NewEngine(width, height, "graphics")
	if err != nil {
		return err
	}

	useShadows := true

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(mgl32.Vec3{0, 3, 6})
	cameraObject.Transform().LookAt(mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 1, 0})
	cameraObject.AddComponent(components.NewCamera(70, width, height, 0.01, 1000))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	cameraObject.AddComponent(&components.HeadHeight{})
	engine.AddObject(cameraObject)

	whiteMaterial := rendering.NewMaterial()
	whiteMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/white.png"))

	tealMaterial := rendering.NewMaterial()
	tealMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/teal.png"))

	floor := core.NewGameObject()
	floor.Transform().SetScale(mgl32.Vec3{100, 0.01, 100})
	floor.Transform().SetPos(mgl32.Vec3{0, -0.005, 0})
	handle(core.LoadModel(floor, "res/meshes/cube/model.obj", whiteMaterial))
	engine.AddObject(floor)

	dirLight := core.NewGameObject()
	dirLight.Transform().SetPos(mgl32.Vec3{8, 8, 0})
	dirLight.Transform().SetScale(mgl32.Vec3{0.5, 0.1, 0.5})
	dirLight.AddComponent(lights.NewDirectional(useShadows, 0.9, 0.9, 0.9, 1))
	//core.LoadModel(dirLight, "res/meshes/cube/model.obj", tealMaterial)
	engine.AddObject(dirLight)

	pointLight := core.NewGameObject()
	pointLight.Transform().SetPos(mgl32.Vec3{-3, 1, 0})
	pointLight.Transform().SetScale(mgl32.Vec3{0.05, 0.05, 0.05})
	pointLight.AddComponent(lights.NewPoint(1, 0, 0, 10))
	//core.LoadModel(pointLight, "res/meshes/cube/model.obj", tealMaterial)
	engine.AddObject(pointLight)

	spot := core.NewGameObject()
	spot.Transform().SetPos(mgl32.Vec3{-5, 2, 5})
	spot.Transform().SetScale(mgl32.Vec3{0.05, 0.05, 0.3})
	spot.Transform().LookAt(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	spot.AddComponent(lights.NewSpot(useShadows, 0.9, 0.4, 0.1, 20, 45))
	//core.LoadModel(spot, "res/meshes/cube/model.obj", tealMaterial)
	engine.AddObject(spot)

	bot := core.NewGameObject()
	bot.Transform().SetPos(mgl32.Vec3{0, 0, 0})
	bot.AddComponent(components.NewRotator(mgl32.Vec3{0, -1, 0}, 23))
	if err := core.LoadModel(bot, "res/meshes/sphere_bot/model.obj", whiteMaterial); err != nil {
		return err
	}
	engine.AddObject(bot)

	cube := core.NewGameObject()
	cube.Transform().SetScale(mgl32.Vec3{1, 2, 8})
	cube.Transform().SetPos(mgl32.Vec3{4, 1, 0})
	if err := core.LoadModel(cube, "res/meshes/cube/model.obj", whiteMaterial); err != nil {
		return err
	}
	engine.AddObject(cube)

	{
		cube := core.NewGameObject()
		cube.Transform().SetScale(mgl32.Vec3{1, 2, 8})
		cube.Transform().SetPos(mgl32.Vec3{-5, 1, -7})
		cube.Transform().SetRot(mgl32.QuatRotate(math.Pi/2, mgl32.Vec3{0, 1, 0}))
		if err := core.LoadModel(cube, "res/meshes/cube/model.obj", whiteMaterial); err != nil {
			return err
		}
		engine.AddObject(cube)
	}

	{
		cube := core.NewGameObject()
		cube.Transform().SetPos(mgl32.Vec3{4, 3.5, 0})
		cube.Transform().SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
		if err := core.LoadModel(cube, "res/meshes/cube/model.obj", tealMaterial); err != nil {
			return err
		}
		engine.AddObject(cube)
	}

	{
		cube := core.NewGameObject()
		cube.Transform().SetPos(mgl32.Vec3{-8, 0.5, 5})
		cube.Transform().SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
		if err := core.LoadModel(cube, "res/meshes/cube/model.obj", tealMaterial); err != nil {
			return err
		}
		engine.AddObject(cube)
	}

	{
		cube := core.NewGameObject()
		cube.Transform().SetPos(mgl32.Vec3{-5, 0.5, 8})
		cube.Transform().SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
		if err := core.LoadModel(cube, "res/meshes/cube/model.obj", tealMaterial); err != nil {
			return err
		}
		engine.AddObject(cube)
	}

	{
		cube := core.NewGameObject()
		cube.Transform().SetPos(mgl32.Vec3{-5, 0.5, -0})
		cube.Transform().SetScale(mgl32.Vec3{0.5, 0.5, 0.5})
		if err := core.LoadModel(cube, "res/meshes/cube/model.obj", tealMaterial); err != nil {
			return err
		}
		engine.AddObject(cube)
	}

	{
		cube := core.NewGameObject()
		cube.Transform().SetPos(mgl32.Vec3{2.4, 0, 0})
		cube.Transform().SetScale(mgl32.Vec3{0.1, 0.5, 0.5})
		if err := core.LoadModel(cube, "res/meshes/cube/model.obj", tealMaterial); err != nil {
			return err
		}
		engine.AddObject(cube)
	}

	engine.Start()

	return nil
}

func handle(err error) {
	if err == nil {
		return
	}
	fmt.Println(err)
	os.Exit(1)
}
