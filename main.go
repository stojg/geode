package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/core"
	"github.com/stojg/graphics/lib/lights"
	"github.com/stojg/graphics/lib/rendering"
	"github.com/stojg/graphics/lib/rendering/terrain"
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
	width := 1024
	height := 800

	engine, err := core.NewEngine(width, height, "graphics")
	if err != nil {
		return err
	}

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(vec3(0, 1.8, 6))
	cameraObject.Transform().SetScale(vec3(0.1, 0.1, 0.1))
	cameraObject.AddComponent(components.NewCamera(75, width, height, 0.1, 2000))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	cameraObject.AddComponent(&components.HeadHeight{})
	engine.AddObject(cameraObject)

	directionalLight := lights.NewDirectional(10, 0.9, 0.9, 0.9, 1)
	dirLight := core.NewGameObject()
	dirLight.Transform().SetPos(vec3(1, 1, 0))
	dirLight.Transform().LookAt(vec3(0, 0, 0), up())
	dirLight.AddComponent(directionalLight)
	engine.AddObject(dirLight)

	spotLight := lights.NewSpot(0.9, 0.4, 0.1, 500, 65)
	spot := core.NewGameObject()
	spot.Transform().SetPos(vec3(3, 3.5, 4.6))
	spot.Transform().LookAt(vec3(0, 1, 0), up())
	spot.AddComponent(spotLight)
	engine.AddObject(spot)

	pointLight := core.NewGameObject()
	pointLight.Transform().SetPos(vec3(-2, 0.6, 2))
	pointLight.AddComponent(lights.NewPoint(0, 0.5, 1.0, 50))
	engine.AddObject(pointLight)

	{
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(-10, 0.3, 0))
		pointLight.AddComponent(lights.NewPoint(0.0, 0.5, 1.0, 50))
		lightMaterial := rendering.NewMaterial()
		lightMaterial.SetAlbedo(mgl32.Vec3{0.1, 0.05, 0.98})
		engine.AddObject(pointLight)
	}

	{
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(2, 0.4, 4))
		pointLight.AddComponent(lights.NewPoint(1.0, 1.0, 1.0, 50))
		//pointLight.AddComponent(components.NewTimeMove(mgl32.Vec3{-1, 0, 0}, func(elapsed float64) float64 {
		//	return math.Sin(glfw.GetTime())
		//}))
		engine.AddObject(pointLight)
	}

	for i := 0; i < 100; i++ {
		cube, err := loadModel("cube")
		handleError(err)
		engine.AddObject(cube)
		cube.Transform().SetPos(vec3(rand.Float32()*100-50, 1, rand.Float32()*100-50))
	}

	for i := -1; i < 1; i++ {
		for j := -1; j < 1; j++ {
			t := terrain.New(float32(i), float32(j))
			to, err := loadModelFromMesh(t.Mesh(), "dry-dirt")
			to.Transform().SetPos(vec3(t.X(), 0, t.Z()))
			handleError(err)
			engine.AddTerrain(to)
		}
	}

	{
		bot, err := loadModel("bot")
		handleError(err)
		bot.AddComponent(components.NewRotator(vec3(0, -1, 0), 15))
		engine.AddObject(bot)
	}

	engine.Start()
	return nil
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func vec3(x, y, z float32) mgl32.Vec3 {
	return mgl32.Vec3{x, y, z}
}

func up() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 0}
}
