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
	width := 800
	height := 600

	engine, err := core.NewEngine(width, height, "graphics")
	if err != nil {
		return err
	}

	whiteMaterial := rendering.NewMaterial()
	whiteMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/white.png"))

	tealMaterial := rendering.NewMaterial()
	tealMaterial.AddTexture("diffuse", rendering.NewTexture("res/textures/teal.png"))

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(vec3(0, 1.8, 6))
	cameraObject.Transform().SetScale(vec3(0.1, 0.1, 0.1))
	cameraObject.AddComponent(components.NewCamera(75, width, height, 0.01, 500))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	cameraObject.AddComponent(&components.HeadHeight{})
	loadModel(cameraObject, "res/meshes/sphere/model.obj", whiteMaterial)
	engine.AddObject(cameraObject)

	//directionalLight := lights.NewDirectional(9, 0.9, 0.9, 0.9, 1)
	//dirLight := core.NewGameObject()
	//dirLight.Transform().SetPos(vec3(1, 2, -1))
	//dirLight.Transform().LookAt(vec3(0, 0, 0), up())
	//dirLight.Transform().SetScale(vec3(0.5, 0.1, 0.5))
	//dirLight.AddComponent(directionalLight)
	//engine.AddObject(dirLight)
	//
	//spotLight := lights.NewSpot(8, 0.9, 0.4, 0.1, 20, 45)
	//spot := core.NewGameObject()
	//spot.Transform().SetPos(vec3(3, 3.5, 4.6))
	//spot.Transform().SetScale(vec3(0.05, 0.05, 0.3))
	//spot.Transform().LookAt(vec3(0, 1, 0), up())
	//spot.AddComponent(spotLight)
	//engine.AddObject(spot)

	pointLight := core.NewGameObject()
	pointLight.Transform().SetPos(vec3(-2, 0.6, 2))
	pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
	pointLight.AddComponent(lights.NewPoint(0.98, 0.05, 0.02, 8))
	loadModel(pointLight, "res/meshes/ico/model.obj", tealMaterial)
	engine.AddObject(pointLight)

	{
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(-10, 0.3, 0))
		pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
		pointLight.AddComponent(lights.NewPoint(0.1, 0.05, 0.98, 5))
		loadModel(pointLight, "res/meshes/ico/model.obj", tealMaterial)
		engine.AddObject(pointLight)
	}

	{
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(2, 0.4, 4))
		pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
		pointLight.AddComponent(lights.NewPoint(0.1, 0.8, 0.23, 5))
		loadModel(pointLight, "res/meshes/ico/model.obj", tealMaterial)
		engine.AddObject(pointLight)
	}

	for i := 0; i < 4; i++ {
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(rand.Float32()*20-15, 0.5, rand.Float32()*20-10))
		pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
		pointLight.AddComponent(lights.NewPoint(rand.Float32(), rand.Float32(), rand.Float32(), 1))
		loadModel(pointLight, "res/meshes/ico/model.obj", tealMaterial)
		engine.AddObject(pointLight)
	}

	floor := core.NewGameObject()
	floor.Transform().SetScale(vec3(15, 0.01, 15))
	floor.Transform().SetPos(vec3(0, -0.005, 0))
	loadModel(floor, "res/meshes/cube/model.obj", whiteMaterial)
	engine.AddObject(floor)

	bot := core.NewGameObject()
	bot.Transform().SetPos(vec3(0, 0.2, 0))
	bot.AddComponent(components.NewRotator(vec3(0, -1, 0), 23))
	loadModel(bot, "res/meshes/sphere_bot/model.obj", whiteMaterial)
	engine.AddObject(bot)

	{ //podium
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(1.9, 0.1, 1.9))
		cube.Transform().SetPos(vec3(0, 0.1, 0))
		loadModel(cube, "res/meshes/cube/model.obj", tealMaterial)
		engine.AddObject(cube)
	}

	{ // wall 1
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(1, 2, 8))
		cube.Transform().SetPos(vec3(4, 2, 2))
		loadModel(cube, "res/meshes/cube/model.obj", whiteMaterial)
		engine.AddObject(cube)
	}

	{ // wall 2
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(1, 2, 8))
		cube.Transform().SetPos(vec3(-5, 2, -5))
		cube.Transform().SetRot(mgl32.QuatRotate(math.Pi/2, vec3(0, 1, 0)))
		loadModel(cube, "res/meshes/cube/model.obj", whiteMaterial)
		engine.AddObject(cube)
	}

	{ // wall 3
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(1, 2, 8))
		cube.Transform().SetPos(vec3(-14, 2, 2))
		loadModel(cube, "res/meshes/cube/model.obj", whiteMaterial)
		engine.AddObject(cube)
	}

	{ // wall 3
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(3, 0.1, 8))
		cube.Transform().SetPos(vec3(-12, 4, 2))
		loadModel(cube, "res/meshes/cube/model.obj", whiteMaterial)
		engine.AddObject(cube)
	}

	{ // on top of wall 1
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(0.5, 0.5, 0.5))
		cube.Transform().SetPos(vec3(4, 4.45, 0))
		loadModel(cube, "res/meshes/cube/model.obj", whiteMaterial)
		engine.AddObject(cube)
	}

	{ // pillar
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(0.25, 3, 0.25))
		cube.Transform().SetPos(vec3(-5, 2.99, 8))
		loadModel(cube, "res/meshes/cube/model.obj", tealMaterial)
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
