package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/core"
	"github.com/stojg/graphics/lib/lights"
	"github.com/stojg/graphics/lib/rendering/terrain"
)

func main() {
	rand.Seed(19)
	l := newLogger("gl.log")

	err := run(l)

	if err != nil {
		l.ErrorLn(err)
		if err := l.close(); err != nil {
			fmt.Println(".. in addition the log file had problem closing", err)
		}
		os.Exit(1)
	}
	if err := l.close(); err != nil {
		fmt.Println(".. in addition the log file had problem closing", err)
	}
}

func run(l *logger) error {
	width := 20
	height := 20

	engine, err := core.NewEngine(width, height, "graphics", l)
	if err != nil {
		return err
	}

	terrain1 := terrain.New(float32(-0.5), float32(-0.5))
	terrainObj, _ := loadModelFromMesh(terrain1.Mesh(), "dry-dirt")
	terrainObj.Transform().SetPos(vec3(terrain1.X(), 0, terrain1.Z()))
	engine.AddTerrain(terrainObj)

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(vec3(10, 0, -10))
	cameraObject.Transform().SetScale(vec3(0.1, 0.1, 0.1))
	cameraObject.AddComponent(components.NewCamera(75, width, height, 0.1, 512))
	cameraObject.AddComponent(&components.FreeMove{Speed: 5})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	cameraObject.Transform().LookAt(vec3(4, 1, 1), up())
	cameraObject.AddComponent(&components.HeadHeight{Terrain: terrain1})
	engine.AddObject(cameraObject)

	directionalLight := lights.NewDirectional(10, 0.9, 0.9, 0.9, 10)
	dirLight := core.NewGameObject()
	dirLight.Transform().SetPos(vec3(1, 1, 0))
	dirLight.Transform().LookAt(vec3(0, 0, 0), up())
	dirLight.AddComponent(directionalLight)
	engine.AddObject(dirLight)

	//spotLight := lights.NewSpot(0.9, 0.4, 0.1, 500, 65)
	//spot := core.NewGameObject()
	//spot.Transform().SetPos(vec3(3, 3.5, 4.6))
	//spot.Transform().LookAt(vec3(0, 1, 0), up())
	//spot.AddComponent(spotLight)
	//engine.AddObject(spot)

	pointLight := core.NewGameObject()
	pointLight.Transform().SetPos(vec3(-2, terrain1.Height(-2, 2)+0.2, 2))
	pointLight.AddComponent(lights.NewPoint(0, 0.5, 1.0, 50))
	engine.AddObject(pointLight)

	{
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(-10, terrain1.Height(-10, 0)+0.2, 0))
		pointLight.AddComponent(lights.NewPoint(0.0, 0.5, 1.0, 50))
		//lightMaterial := resources.NewMaterial()
		//lightMaterial.SetAlbedo(mgl32.Vec3{0.1, 0.05, 0.98})
		engine.AddObject(pointLight)
	}

	tSize := float32(512)
	tHalfSize := tSize / 2
	for i := 0; i < 200; i++ {
		p := core.NewGameObject()
		p, err := loadModel("cube")
		handleError(err)
		x, z := rand.Float32()*tSize-tHalfSize, rand.Float32()*tSize-tHalfSize
		p.Transform().SetPos(vec3(x, terrain1.Height(x, z)+0.5, z))
		p.Transform().SetScale(vec3(0.5, 0.5, 0.5))
		engine.AddObject(p)
	}

	{
		bot, err := loadModel("bot")
		handleError(err)
		bot.AddComponent(components.NewRotator(vec3(0, -1, 0), 15))
		bot.Transform().SetPos(vec3(0, terrain1.Height(0, 0), 0))
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
