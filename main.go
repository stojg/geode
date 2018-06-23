package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/go-gl/glfw/v3.2/glfw"
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

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(vec3(0, 1.8, 6))
	cameraObject.Transform().SetScale(vec3(0.1, 0.1, 0.1))
	cameraObject.AddComponent(components.NewCamera(75, width, height, 0.1, 2000))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	//cameraObject.AddComponent(&components.HeadHeight{})
	//loadModel(cameraObject, "res/meshes/sphere/model.obj", whiteMaterial)
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
	engine.AddObject(pointLight)

	{
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(-10, 0.3, 0))
		pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
		pointLight.AddComponent(lights.NewPoint(0.0, 0.5, 1.0, 50))
		lightMaterial := rendering.NewMaterial()
		lightMaterial.SetAlbedo(mgl32.Vec3{0.1, 0.05, 0.98})
		engine.AddObject(pointLight)
	}

	{
		pointLight := core.NewGameObject()
		pointLight.Transform().SetPos(vec3(2, 0.4, 4))
		pointLight.Transform().SetScale(vec3(0.05, 0.05, 0.05))
		pointLight.AddComponent(lights.NewPoint(1.0, 1.0, 1.0, 50))
		pointLight.AddComponent(components.NewTimeMove(mgl32.Vec3{-1, 0, 0}, func(elapsed float64) float64 {
			return math.Sin(glfw.GetTime())
		}))
		lightMaterial := rendering.NewMaterial()
		lightMaterial.SetAlbedo(mgl32.Vec3{50, 50, 50})
		engine.AddObject(pointLight)
	}

	var whiteMaterial []*rendering.Material
	plasticMtrl := rendering.NewMaterial()
	plasticMtrl.AddTexture("albedo", rendering.NewTexture("res/textures/scuffed-plastic/scuffed-plastic5-alb.png", true))
	plasticMtrl.AddTexture("metallic", rendering.NewMetallicTexture("res/textures/scuffed-plastic/scuffed-plastic-metal.png"))
	plasticMtrl.AddTexture("roughness", rendering.NewRoughnessTexture("res/textures/scuffed-plastic/scuffed-plastic-rough.png"))
	plasticMtrl.AddTexture("normal", rendering.NewTexture("res/textures/scuffed-plastic/scuffed-plastic-normal.png", false))
	whiteMaterial = append(whiteMaterial, plasticMtrl)

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			cube := core.NewGameObject()
			cube.Transform().SetScale(vec3(10, 0.04, 10))
			cube.Transform().SetPos(vec3(float32(i)*10, 0.00, float32(j)*10))
			loadModel(cube, "res/meshes/cube/model.obj", whiteMaterial)
			engine.AddObject(cube)
		}
	}

	{
		bot := core.NewGameObject()
		bot.Transform().SetPos(vec3(0, 0, 0))
		bot.AddComponent(components.NewRotator(vec3(0, -1, 0), 15))

		var mtrls []*rendering.Material
		outer := rendering.NewMaterial()
		outer.AddTexture("albedo", rendering.NewTexture("res/textures/sphere_bot/Robot_outerbody_Albedo.png", true))
		outer.AddTexture("metallic", rendering.NewMetallicTexture("res/textures/sphere_bot/Robot_outerbody_Metallic.png"))
		outer.AddTexture("roughness", rendering.NewRoughnessTexture("res/textures/sphere_bot/Robot_outerbody_Roughness.png"))
		outer.AddTexture("normal", rendering.NewTexture("res/textures/sphere_bot/Robot_outerbody_Normal.png", false))
		mtrls = append(mtrls, outer)

		inner := rendering.NewMaterial()
		inner.AddTexture("albedo", rendering.NewTexture("res/textures/sphere_bot/Robot_innerbody_Albedo.png", true))
		inner.AddTexture("metallic", rendering.NewMetallicTexture("res/textures/sphere_bot/Robot_innerbody_Metallic.png"))
		inner.AddTexture("roughness", rendering.NewRoughnessTexture("res/textures/sphere_bot/Robot_innerbody_Roughness.png"))
		inner.AddTexture("normal", rendering.NewTexture("res/textures/sphere_bot/Robot_innerbody_Normal.png", false))

		mtrls = append(mtrls, inner)

		loadModel(bot, "res/meshes/sphere_bot/model.obj", mtrls)
		engine.AddObject(bot)
	}

	engine.Start()
	return nil
}

func loadModel(g *core.GameObject, obj string, material []*rendering.Material) {
	objData, err := loader.Load(obj)
	if err != nil {
		fmt.Printf("Model loading failed: %v", err)
		os.Exit(1)
	}

	if _, ok := models[obj]; !ok {
		var meshes []*rendering.Mesh
		if len(objData) != len(material) {
			fmt.Printf("Have %d meshes in object, but only %d materials\n", len(objData), len(material))
		}
		for i, data := range objData {
			mesh := rendering.NewMesh()
			mesh.SetVertices(rendering.ConvertToVertices(data))
			meshes = append(meshes, mesh)
			g.AddComponent(components.NewMeshRenderer(mesh, material[i]))
		}
		models[obj] = meshes
	}
}

func vec3(x, y, z float32) mgl32.Vec3 {
	return mgl32.Vec3{x, y, z}
}

func up() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 0}
}
