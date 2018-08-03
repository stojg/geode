package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/profile"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/core"
	"github.com/stojg/graphics/lib/lights"
	"github.com/stojg/graphics/lib/rendering/terrain"
	"github.com/stojg/graphics/lib/resources"
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

	terrainA := terrain.New(float32(-0.5), float32(-0.5))
	terrainObj, _ := loadModelFromMesh(terrainA.Mesh(), "dry-dirt", components.R_TERRAIN|components.R_SHADOWED)
	txt := resources.NewTexture("res/textures/dry-dirt2/albedo.png", true)
	terrainObj.Model().Material().AddTexture("albedo2", txt)

	terrainObj.Transform().SetPos(vec3(terrainA.X(), 0, terrainA.Z()))
	engine.AddObject(terrainObj)

	cameraObject := core.NewGameObject(components.R_NA)
	cameraObject.SetPos(-10, terrainA.Height(10, -10)+20, -10)
	cameraObject.AddComponent(components.NewCamera(75, engine.Width(), engine.Height(), 0.1, 512))
	cameraObject.AddComponent(components.NewFreeMove(4))
	cameraObject.AddComponent(components.NewFreelook(engine.Width(), engine.Height()))
	cameraObject.AddComponent(components.NewHeadHeight(terrainA))
	engine.AddObject(cameraObject)

	p1 := core.NewParticleSystem(1000)
	p1.SetPos(10, terrainA.Height(10, terrainA.Height(10, 0)), 0)
	engine.AddObject(p1)

	sun := core.NewGameObject(components.R_LIGHT)
	sun.SetPos(1, 0.75, 0)
	sun.Transform().LookAt(vec3(0, 0, 0), up())
	sun.AddComponent(lights.NewDirectional(11, 0.996, 0.863, 0.533, 10))
	engine.AddObject(sun)

	spot := core.NewGameObject(components.R_LIGHT)
	spot.SetPos(3, 3.5, 4.6)
	spot.Transform().LookAt(vec3(0, 1, 0), up())
	spot.AddComponent(lights.NewSpot(0.9, 0.4, 0.1, 200, 65))
	engine.AddObject(spot)

	pointLightA := core.NewGameObject(components.R_LIGHT)
	pointLightA.SetPos(-2, terrainA.Height(-2, 10)+0.5, 10)
	pointLightA.AddComponent(lights.NewPoint(0, 0.5, 1.0, 50))
	engine.AddObject(pointLightA)

	pointLightB := core.NewGameObject(components.R_LIGHT)
	pointLightB.SetPos(-10, terrainA.Height(-10, 0)+0.5, 0)
	pointLightB.AddComponent(lights.NewPoint(0.0, 0.5, 1.0, 50))
	engine.AddObject(pointLightB)

	terrainSize := float32(terrain.Size)
	for i := 0; i < 200; i++ {
		p := core.NewGameObject(components.R_DEFAULT | components.R_SHADOWED)
		p, err := loadModel("cube")
		handleError(err)
		x, z := rand.Float32()*terrainSize-terrainSize/2, rand.Float32()*terrainSize-terrainSize/2
		p.SetPos(x, terrainA.Height(x, z)+0.25, z)
		p.SetScale(0.5, 0.5, 0.5)
		p.Rotate(up(), rand.Float32()*math.Pi*2)
		engine.AddObject(p)
	}

	for i := 0; i < 200; i++ {
		p := core.NewGameObject(components.R_DEFAULT | components.R_SHADOWED)
		p, err := loadModel("sphere")
		handleError(err)
		x, z := rand.Float32()*terrainSize-terrainSize/2, rand.Float32()*terrainSize-terrainSize/2
		p.SetPos(x, terrainA.Height(x, z)+0.5, z)
		p.SetScale(0.5, 0.5, 0.5)
		engine.AddObject(p)
	}

	for i := 0; i < 200; i++ {
		p := core.NewGameObject(components.R_DEFAULT | components.R_SHADOWED)
		p, err := loadModel("ico")
		handleError(err)
		x, z := rand.Float32()*terrainSize-terrainSize/2, rand.Float32()*terrainSize-terrainSize/2
		p.SetPos(x, terrainA.Height(x, z)+0.5, z)
		p.SetScale(0.5, 0.5, 0.5)
		p.Rotate(up(), rand.Float32()*math.Pi*2)
		engine.AddObject(p)
	}

	bot, err := loadModel("bot")
	handleError(err)
	bot.AddComponent(components.NewRotator(vec3(0, -1, 0), 15))
	bot.SetPos(0, terrainA.Height(0, 0), 0)
	engine.AddObject(bot)

	defer profile.Start().Stop()
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
