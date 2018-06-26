package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/core"
	"github.com/stojg/graphics/lib/rendering"
	"github.com/stojg/graphics/lib/rendering/loader"
)

var meshes map[string][]*rendering.Mesh
var modelTextures map[string]*rendering.Texture

func init() {
	meshes = make(map[string][]*rendering.Mesh)
	modelTextures = make(map[string]*rendering.Texture)
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
	var mtrls []*rendering.Material
	for _, texture := range textures {
		mtrls = append(mtrls, loadMaterial(texture))
	}

	modelFile := fmt.Sprintf("res/meshes/%s/model.obj", mi.Mesh)
	if _, ok := meshes[modelFile]; !ok {
		meshes[modelFile] = loadObject(modelFile, mtrls)
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

func loadMaterial(texture string) *rendering.Material {
	texturePath := fmt.Sprintf("./res/textures/%s", texture)
	material := rendering.NewMaterial()
	txt, ok := modelTextures[texturePath+"/albedo.png"]
	if !ok {
		txt = rendering.NewTexture(texturePath+"/albedo.png", true)
		modelTextures[texturePath+"/albedo.png"] = txt
	}
	material.AddTexture("albedo", txt)
	txt, ok = modelTextures[texturePath+"/metallic.png"]
	if !ok {
		txt = rendering.NewMetallicTexture(texturePath + "/metallic.png")
		modelTextures[texturePath+"/metallic.png"] = txt
	}
	material.AddTexture("metallic", txt)
	txt, ok = modelTextures[texturePath+"/roughness.png"]
	if !ok {
		txt = rendering.NewRoughnessTexture(texturePath + "/roughness.png")
		modelTextures[texturePath+"/roughness.png"] = txt
	}
	material.AddTexture("roughness", txt)
	txt, ok = modelTextures[texturePath+"/normal.png"]
	if !ok {
		txt = rendering.NewTexture(texturePath+"/normal.png", false)
		modelTextures[texturePath+"/normal.png"] = txt
	}
	material.AddTexture("normal", txt)

	return material
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
