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
	w := 1024
	h := int(float32(w) / (4.0 / 3.0))
	engine, err := core.NewEngine(w, h, "graphics", l)
	if err != nil {
		return err
	}
	width, height := engine.Width(), engine.Height()

	terrainA := terrain.New(float32(-0.5), float32(-0.5))
	terrainObj, _ := loadModelFromMesh(terrainA.Mesh(), "dry-dirt")
	terrainObj.Transform().SetPos(vec3(terrainA.X(), 0, terrainA.Z()))
	engine.AddTerrain(terrainObj)

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(vec3(10, 0, -10))
	cameraObject.Transform().SetScale(vec3(0.1, 0.1, 0.1))
	cameraObject.AddComponent(components.NewCamera(75, width, height, 0.1, 512))
	cameraObject.AddComponent(components.NewFreeMove(5))
	cameraObject.AddComponent(components.NewFreelook(width, height))
	cameraObject.Transform().LookAt(vec3(4, 1, 1), up())
	cameraObject.AddComponent(components.NewHeadHeight(terrainA))
	engine.AddObject(cameraObject)

	sun := core.NewGameObject()
	sun.Transform().SetPos(vec3(1, 1, 0))
	sun.Transform().LookAt(vec3(0, 0, 0), up())
	sun.AddComponent(lights.NewDirectional(10, 0.9, 0.9, 0.9, 10))
	engine.AddObject(sun)

	spot := core.NewGameObject()
	spot.Transform().SetPos(vec3(3, 3.5, 4.6))
	spot.Transform().LookAt(vec3(0, 1, 0), up())
	spot.AddComponent(lights.NewSpot(0.9, 0.4, 0.1, 500, 65))
	engine.AddObject(spot)

	pointLightA := core.NewGameObject()
	pointLightA.Transform().SetPos(vec3(-2, terrainA.Height(-2, 2)+0.2, 2))
	pointLightA.AddComponent(lights.NewPoint(0, 0.5, 1.0, 50))
	engine.AddObject(pointLightA)

	pointLightB := core.NewGameObject()
	pointLightB.Transform().SetPos(vec3(-10, terrainA.Height(-10, 0)+0.2, 0))
	pointLightB.AddComponent(lights.NewPoint(0.0, 0.5, 1.0, 50))
	engine.AddObject(pointLightB)

	tSize := float32(terrain.Size)
	tHalfSize := tSize / 2
	for i := 0; i < 800; i++ {
		p := core.NewGameObject()
		p, err := loadModel("cube")
		handleError(err)
		x, z := rand.Float32()*tSize-tHalfSize, rand.Float32()*tSize-tHalfSize
		p.Transform().SetPos(vec3(x, terrainA.Height(x, z)+0.5, z))
		p.Transform().SetScale(vec3(0.5, 0.5, 0.5))
		p.Transform().Rotate(up(), rand.Float32()*math.Pi*2)
		engine.AddObject(p)
	}

	{
		bot, err := loadModel("bot")
		handleError(err)
		bot.AddComponent(components.NewRotator(vec3(0, -1, 0), 15))
		bot.Transform().SetPos(vec3(0, terrainA.Height(0, 0), 0))
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
