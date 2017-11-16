package loader

import "fmt"

type objectFile struct {
	Geometry          *geometry
	MaterialLibraries []string

	Objects  []*object
	Comments []string
}

func (o *objectFile) ObjectWithType(t dataType) (objects []*object) {
	for _, o := range o.Objects {
		if o.Type == t {
			objects = append(objects, o)
		}
	}
	return objects
}

func (o *objectFile) createObject(t dataType, name, material string) *object {
	if t != tChildObject && t != tChildGroup {
		fmt.Printf("createObject: invalid object type %s", t)
		return nil
	}
	child := &object{
		Type:     t,
		Name:     name,
		Material: material,
		parent:   o,
	}
	if child.Name == "" {
		child.Name = fmt.Sprintf("%s_%d", t.name(), len(o.ObjectWithType(t))+1)
	}
	o.Objects = append(o.Objects, child)
	return child
}
