package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/core"
	"github.com/stojg/graphics/lib/rendering/loader"
	"github.com/stojg/graphics/lib/resources"
)

var meshes map[string][]*resources.Mesh
var modelTextures map[string]*resources.Texture

func init() {
	meshes = make(map[string][]*resources.Mesh)
	modelTextures = make(map[string]*resources.Texture)
}

func loadModel(modelName string) (*core.GameObject, error) {
	var mi struct {
		Mesh     string   `json:"mesh"`
		Textures []string `json:"textures"`
	}
	d, err := ioutil.ReadFile(fmt.Sprintf("res/models/%s.json", modelName))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(d, &mi); err != nil {
		return nil, err
	}

	textures := mi.Textures
	var mtrls []*resources.Material
	for _, texture := range textures {
		mtrls = append(mtrls, loadMaterial(texture))
	}

	modelFile := fmt.Sprintf("res/meshes/%s/model.obj", mi.Mesh)
	if _, ok := meshes[modelFile]; !ok {
		meshes[modelFile] = loadMeshesFromObj(modelFile, mtrls)
	}

	gameObject := core.NewGameObject()
	for idx, m := range meshes[modelFile] {
		gameObject.AddComponent(components.NewModel(m, mtrls[idx]))
	}

	return gameObject, nil
}

func loadModelFromMesh(mesh components.Drawable, texture string) (*core.GameObject, error) {
	material := loadMaterial(texture)
	gameObject := core.NewGameObject()
	gameObject.AddComponent(components.NewModel(mesh, material))
	return gameObject, nil
}

func loadMaterial(texture string) *resources.Material {
	texturePath := fmt.Sprintf("./res/textures/%s", texture)
	material := resources.NewMaterial()
	txt, ok := modelTextures[texturePath+"/albedo.png"]
	if !ok {
		txt = resources.NewTexture(texturePath+"/albedo.png", true)
		modelTextures[texturePath+"/albedo.png"] = txt
	}
	material.AddTexture("albedo", txt)
	txt, ok = modelTextures[texturePath+"/metallic.png"]
	if !ok {
		txt = resources.NewMetallicTexture(texturePath + "/metallic.png")
		modelTextures[texturePath+"/metallic.png"] = txt
	}
	material.AddTexture("metallic", txt)
	txt, ok = modelTextures[texturePath+"/roughness.png"]
	if !ok {
		txt = resources.NewRoughnessTexture(texturePath + "/roughness.png")
		modelTextures[texturePath+"/roughness.png"] = txt
	}
	material.AddTexture("roughness", txt)
	txt, ok = modelTextures[texturePath+"/normal.png"]
	if !ok {
		txt = resources.NewTexture(texturePath+"/normal.png", false)
		modelTextures[texturePath+"/normal.png"] = txt
	}
	material.AddTexture("normal", txt)

	return material
}

func loadMeshesFromObj(obj string, material []*resources.Material) []*resources.Mesh {
	objVert, objInd, err := loader.Load(obj)
	if err != nil {
		fmt.Printf("Model loading failed: %v", err)
		os.Exit(1)
	}
	if len(objVert) != len(material) {
		fmt.Printf("Have %d meshes in object, but only %d materials\n", len(objVert), len(material))
	}
	var meshes []*resources.Mesh
	for i, data := range objVert {
		mesh := resources.NewMesh()
		mesh.SetVertices(resources.ConvertToVertices(data, objInd[i]), objInd[i])
		meshes = append(meshes, mesh)
	}
	return meshes
}
