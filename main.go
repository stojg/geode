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
	"github.com/stojg/graphics/lib/rendering/loader"
	"github.com/stojg/graphics/lib/rendering/terrain"
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

	var whiteMaterial []*rendering.Material
	plasticMtrl := rendering.NewMaterial()
	plasticMtrl.AddTexture("albedo", rendering.NewTexture("res/textures/scuffed-plastic/scuffed-plastic5-alb.png", true))
	plasticMtrl.AddTexture("metallic", rendering.NewMetallicTexture("res/textures/scuffed-plastic/scuffed-plastic-metal.png"))
	plasticMtrl.AddTexture("roughness", rendering.NewRoughnessTexture("res/textures/scuffed-plastic/scuffed-plastic-rough.png"))
	plasticMtrl.AddTexture("normal", rendering.NewTexture("res/textures/scuffed-plastic/scuffed-plastic-normal.png", false))
	whiteMaterial = append(whiteMaterial, plasticMtrl)

	cameraObject := core.NewGameObject()
	cameraObject.Transform().SetPos(vec3(0, 1.8, 6))
	cameraObject.Transform().SetScale(vec3(0.1, 0.1, 0.1))
	cameraObject.AddComponent(components.NewCamera(75, width, height, 0.1, 2000))
	cameraObject.AddComponent(&components.FreeMove{})
	cameraObject.AddComponent(components.NewFreelook(width, height))
	//cameraObject.AddComponent(&components.HeadHeight{})
	setMeshRenderer(cameraObject, "res/meshes/sphere/model.obj", whiteMaterial)
	engine.AddObject(cameraObject)

	directionalLight := lights.NewDirectional(10, 0.9, 0.9, 0.9, 1)
	dirLight := core.NewGameObject()
	dirLight.Transform().SetPos(vec3(1, 1, 0))
	dirLight.Transform().LookAt(vec3(0, 0, 0), up())
	dirLight.Transform().SetScale(vec3(0.5, 0.1, 0.5))
	dirLight.AddComponent(directionalLight)
	engine.AddObject(dirLight)

	spotLight := lights.NewSpot(0.9, 0.4, 0.1, 500, 65)
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
		//pointLight.AddComponent(components.NewTimeMove(mgl32.Vec3{-1, 0, 0}, func(elapsed float64) float64 {
		//	return math.Sin(glfw.GetTime())
		//}))
		lightMaterial := rendering.NewMaterial()
		lightMaterial.SetAlbedo(mgl32.Vec3{50, 50, 50})
		engine.AddObject(pointLight)
	}

	cubes := core.NewGameObject()
	setMeshInstanceRenderer(cubes, "res/meshes/cube/model.obj", whiteMaterial)
	for i := 0; i < 100; i++ {
		cube := core.NewGameObject()
		cube.Transform().SetScale(vec3(1, 1, 1))
		cube.Transform().SetPos(vec3(rand.Float32()*100-50, 1, rand.Float32()*100-50))
		cubes.AddChild(cube)
	}
	engine.AddObject(cubes)

	var dryDirt []*rendering.Material
	dryDirtMtrl := rendering.NewMaterial()
	dryDirtMtrl.AddTexture("albedo", rendering.NewTexture("res/textures/dry-dirt/albedo.png", true))
	dryDirtMtrl.AddTexture("metallic", rendering.NewMetallicTexture("res/textures/dry-dirt/metalness.png"))
	dryDirtMtrl.AddTexture("roughness", rendering.NewRoughnessTexture("res/textures/dry-dirt/roughness.png"))
	dryDirtMtrl.AddTexture("normal", rendering.NewTexture("res/textures/dry-dirt/normal2.png", false))
	dryDirt = append(dryDirt, dryDirtMtrl)

	{
		t := terrain.New(0, 0)
		test := core.NewGameObject()
		test.AddComponent(components.NewMeshRenderer(t.Mesh(), dryDirt[0]))
		engine.AddTerrain(test)
	}

	{
		t := terrain.New(-1, 0)
		test := core.NewGameObject()
		test.Transform().SetPos(vec3(t.X(), 0, t.Z()))
		test.AddComponent(components.NewMeshRenderer(t.Mesh(), dryDirt[0]))
		engine.AddTerrain(test)
	}

	{
		t := terrain.New(-1, -1)
		test := core.NewGameObject()
		test.Transform().SetPos(vec3(t.X(), 0, t.Z()))
		test.AddComponent(components.NewMeshRenderer(t.Mesh(), dryDirt[0]))
		engine.AddTerrain(test)
	}

	{
		t := terrain.New(0, -1)
		test := core.NewGameObject()
		test.Transform().SetPos(vec3(t.X(), 0, t.Z()))
		test.AddComponent(components.NewMeshRenderer(t.Mesh(), dryDirt[0]))
		engine.AddTerrain(test)
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

		setMeshRenderer(bot, "res/meshes/sphere_bot/model.obj", mtrls)
		engine.AddObject(bot)
	}

	engine.Start()
	return nil
}

func setMeshRenderer(g *core.GameObject, modelFile string, material []*rendering.Material) {

	if _, ok := models[modelFile]; !ok {
		models[modelFile] = loadObject(modelFile, material)
	}

	for idx, m := range models[modelFile] {
		g.AddComponent(components.NewMeshRenderer(m, material[idx]))
	}
}

func setMeshInstanceRenderer(g *core.GameObject, modelFile string, material []*rendering.Material) {

	if _, ok := models[modelFile]; !ok {
		models[modelFile] = loadObject(modelFile, material)
	}

	for idx, m := range models[modelFile] {
		g.AddComponent(components.NewMeshInstanceRenderer(m, material[idx]))
	}
}

func loadObject(obj string, material []*rendering.Material) []*rendering.Mesh {
	objVert, objInd, err := loader.Load(obj)
	if err != nil {
		fmt.Printf("Model loading failed: %v", err)
		os.Exit(1)
	}
	if len(objVert) != len(material) {
		fmt.Printf("Have %d meshes in object, but only %d materials\n", len(objVert), len(material))
	}
	var meshes []*rendering.Mesh
	for i, data := range objVert {
		mesh := rendering.NewMesh()
		mesh.SetVertices(rendering.ConvertToVertices(data, objInd[i]), objInd[i])
		meshes = append(meshes, mesh)
	}
	return meshes
}

func vec3(x, y, z float32) mgl32.Vec3 {
	return mgl32.Vec3{x, y, z}
}

func up() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 0}
}
