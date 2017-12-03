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
	"github.com/stojg/graphics/lib/rendering/loader"
)

var models map[string][]*rendering.Mesh

func main() {
	rand.Seed(19)
	l := newLogger("gl.log")
	models = make(map[string][]*rendering.Mesh)
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

	whiteMaterial := rendering.NewMaterial()
	whiteMaterial.SetAlbedo(mgl32.Vec3{0.961, 0.922, 0.898})

	tealMaterial := rendering.NewMaterial()
	tealMaterial.SetAlbedo(mgl32.Vec3{0.447, 0.792, 0.918})

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(vec3(0, 1.8, 6))
	cameraObject.Transform().SetScale(vec3(0.1, 0.1, 0.1))
	cameraObject.AddComponent(components.NewCamera(75, width, height, 0.01, 500))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	cameraObject.AddComponent(&components.HeadHeight{})
	loadModel(cameraObject, "res/meshes/sphere/model.obj", whiteMaterial)
	engine.AddObject(cameraObject)

	directionalLight := lights.NewDirectional(9, 0.9, 0.9, 0.9, 1)
	dirLight := core.NewGameObject()
	dirLight.Transform().SetPos(vec3(1, 1, 0))
	dirLight.Transform().LookAt(vec3(0, 0, 0), up())
	dirLight.Transform().SetScale(vec3(0.5, 0.1, 0.5))
	dirLight.AddComponent(directionalLight)
	engine.AddObject(dirLight)

	spotLight := lights.NewSpot(8, 0.9, 0.4, 0.1, 30, 65)
	spot := core.NewGameObject()
	spot.Transform().SetPos(vec3(3, 3.5, 4.6))
	spot.Transform().SetScale(vec3(0.05, 0.05, 0.3))
	spot.Transform().LookAt(vec3(0, 1, 0), up())
	spot.AddComponent(spotLight)
	engine.AddObject(spot)

	pointLight := core.NewGameObject()
	pointLight.Transform().SetPos(vec3(-2, 0.6, 2))
	pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
	pointLight.AddComponent(lights.NewPoint(0, 0.5, 1.0, 50))
	lightMaterial := rendering.NewMaterial()
	lightMaterial.SetAlbedo(mgl32.Vec3{23.47, 21.31, 20.79})
	loadModel(pointLight, "res/meshes/ico/model.obj", lightMaterial)
	engine.AddObject(pointLight)

	{
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(-10, 0.3, 0))
		pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
		pointLight.AddComponent(lights.NewPoint(0.0, 0.5, 1.0, 50))
		lightMaterial := rendering.NewMaterial()
		lightMaterial.SetAlbedo(mgl32.Vec3{0.1, 0.05, 0.98})
		loadModel(pointLight, "res/meshes/ico/model.obj", lightMaterial)
		engine.AddObject(pointLight)
	}
	//
	//{
	//	pointLight := core.NewGameObject()
	//	pointLight.Transform().SetPos(vec3(2, 0.4, 4))
	//	pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
	//	pointLight.AddComponent(lights.NewPoint(1.0, 1.0, 1.0, 50))
	//	pointLight.AddComponent(components.NewTimeMove(mgl32.Vec3{-1, 0, 0}, func(elapsed float64) float64 {
	//		return math.Sin(glfw.GetTime())
	//	}))
	//	lightMaterial := rendering.NewMaterial()
	//	lightMaterial.SetAlbedo(mgl32.Vec3{50, 50, 50})
	//	loadModel(pointLight, "res/meshes/ico/model.obj", lightMaterial)
	//	engine.AddObject(pointLight)
	//}

	//for i := 0; i < 3; i++ {
	//	pointLight := core.NewGameObject()
	//	pointLight.Transform().SetPos(vec3(rand.Float32()*30-15, 0.5, rand.Float32()*30-10))
	//	pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
	//	r, g, b := rand.Float32(), rand.Float32(), rand.Float32()
	//	pointLight.AddComponent(lights.NewPoint(r, g, b, 50))
	//	lightMaterial := rendering.NewMaterial()
	//	lightMaterial.SetAlbedo(mgl32.Vec3{r, g, b})
	//	loadModel(pointLight, "res/meshes/ico/model.obj", lightMaterial)
	//	engine.AddObject(pointLight)
	//}

	floor := core.NewGameObject()
	floor.Transform().SetScale(vec3(100, 0.01, 100))
	floor.Transform().SetPos(vec3(0, -0.005, 0))
	floorMaterial := rendering.NewMaterial()
	floorMaterial.SetAlbedo(mgl32.Vec3{0.8, 0.8, 0.8})
	floorMaterial.SetRoughness(0.9)
	floorMaterial.SetMetallic(0.02)
	loadModel(floor, "res/meshes/cube/model.obj", floorMaterial)
	engine.AddObject(floor)

	bot := core.NewGameObject()
	bot.Transform().SetPos(vec3(0, 0.2, 0))
	bot.AddComponent(components.NewRotator(vec3(0, -1, 0), 23))
	botMaterial := rendering.NewMaterial()
	//botMaterial.SetAlbedo(mgl32.Vec3{1, 0.765557, 0.336057}) // gold
	botMaterial.SetAlbedo(mgl32.Vec3{0.400, 0.249, 0.000})
	botMaterial.SetRoughness(0.25)
	botMaterial.SetMetallic(0)
	loadModel(bot, "res/meshes/sphere_bot/model.obj", botMaterial)
	engine.AddObject(bot)

	wallMaterial := rendering.NewMaterial()
	wallMaterial.SetAlbedo(mgl32.Vec3{0.8, 0.8, 0.8})
	wallMaterial.SetRoughness(0.9)
	wallMaterial.SetMetallic(0.02)

	{
		sphereMtrl := rendering.NewMaterial()
		sphereMtrl.SetAlbedo(mgl32.Vec3{0.202, 0.545, 0.147})
		sphereMtrl.SetRoughness(0.2)
		sphereMtrl.SetMetallic(0)
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(0.5, 0.5, 0.5))
		cube.Transform().SetPos(vec3(-4, 0.5, 0))
		loadModel(cube, "res/meshes/cube/model.obj", sphereMtrl)
		engine.AddObject(cube)
	}

	{
		sphereMtrl := rendering.NewMaterial()
		sphereMtrl.SetAlbedo(mgl32.Vec3{1, 0.765557, 0.336057})
		sphereMtrl.SetRoughness(0.3)
		sphereMtrl.SetMetallic(1)
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(0.2, 0.2, 0.2))
		cube.Transform().SetPos(vec3(-4, 0.2, 4))
		loadModel(cube, "res/meshes/sphere/model.obj", sphereMtrl)
		engine.AddObject(cube)
		{
			cube := core.NewGameObject()
			cube.Transform().SetScale(vec3(0.2, 0.2, 0.2))
			cube.Transform().SetPos(vec3(-5, 0.2, 4))
			loadModel(cube, "res/meshes/sphere/model.obj", sphereMtrl)
			engine.AddObject(cube)
		}
	}

	{ //podium
		cube := core.NewGameObject()
		podiumMtrl := rendering.NewMaterial()
		podiumMtrl.SetAlbedo(mgl32.Vec3{0.0, 0.0, 0.0})
		podiumMtrl.SetRoughness(0.03)
		podiumMtrl.SetMetallic(1)
		cube.Transform().SetScale(vec3(1.9, 0.1, 1.9))
		cube.Transform().SetPos(vec3(0, 0.1, 0))
		loadModel(cube, "res/meshes/cube/model.obj", podiumMtrl)
		engine.AddObject(cube)
	}

	{ // wall 1
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(1, 2, 8))
		cube.Transform().SetPos(vec3(4, 2, 2))
		loadModel(cube, "res/meshes/cube/model.obj", wallMaterial)
		engine.AddObject(cube)
	}

	{ // wall 2
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(1, 2, 8))
		cube.Transform().SetPos(vec3(-5, 2, -5))
		cube.Transform().SetRot(mgl32.QuatRotate(math.Pi/2, vec3(0, 1, 0)))
		loadModel(cube, "res/meshes/cube/model.obj", wallMaterial)
		engine.AddObject(cube)
	}

	{ // wall 3
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(1, 2, 8))
		cube.Transform().SetPos(vec3(-14, 2, 2))
		loadModel(cube, "res/meshes/cube/model.obj", wallMaterial)
		engine.AddObject(cube)
	}

	{ // roof
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(3, 0.1, 8))
		cube.Transform().SetPos(vec3(-12, 4, 2))
		loadModel(cube, "res/meshes/cube/model.obj", wallMaterial)
		engine.AddObject(cube)
	}

	{ // on top of wall 1
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(0.5, 0.5, 0.5))
		cube.Transform().SetPos(vec3(4, 4.45, 0))
		cubeMaterial := rendering.NewMaterial()
		cubeMaterial.SetAlbedo(mgl32.Vec3{0.0, 0.8, 0.0})
		cubeMaterial.SetRoughness(0.03)
		cubeMaterial.SetMetallic(0.02)
		loadModel(cube, "res/meshes/cube/model.obj", cubeMaterial)
		engine.AddObject(cube)
	}

	{ // pillar
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(0.25, 3, 0.25))
		cube.Transform().SetPos(vec3(-5, 2.99, 8))

		loadModel(cube, "res/meshes/cube/model.obj", wallMaterial)
		engine.AddObject(cube)
	}
	engine.Start()
	return nil
}

func loadModel(g *core.GameObject, obj string, material components.Material) {
	objData, err := loader.Load(obj)
	if err != nil {
		fmt.Printf("Model loading failed: %v", err)
		os.Exit(1)
	}

	if _, ok := models[obj]; !ok {
		var meshes []*rendering.Mesh
		for i, data := range objData {
			mesh := rendering.NewMesh()
			mesh.SetVertices(rendering.ConvertToVertices(data))
			fmt.Printf("loadModel: %s.%d has %d vertices\n", obj, i, mesh.NumVertices())
			meshes = append(meshes, mesh)
		}
		models[obj] = meshes
	}
	for _, mesh := range models[obj] {
		g.AddComponent(components.NewMeshRenderer(mesh, material))
	}

}

func vec3(x, y, z float32) mgl32.Vec3 {
	return mgl32.Vec3{x, y, z}
}

func up() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 0}
}
