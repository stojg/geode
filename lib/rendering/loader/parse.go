package loader

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func parse(src io.Reader) (*objectFile, int, error) {
	root := &objectFile{
		Geometry: newGeometry(),
		Objects:  make([]*object, 0),
	}

	scanner := bufio.NewScanner(src)
	currentLine := 0

	var (
		currentObject           *object
		currentObjectName       string
		currentObjectChildIndex int
		currentMaterial         string
		currentSmoothGroup      string
	)

	// create a "fakeObject" for grouping different materials in one object to several objects in the purpose of
	// minimising 3d drawcalls
	fakeObject := func(material string) *object {
		ot := tChildObject
		// @todo, this might actually always be true
		if currentObject != nil {
			ot = currentObject.Type
		}
		currentObjectChildIndex++
		name := fmt.Sprintf("%s_%d", currentObjectName, currentObjectChildIndex)
		return root.createObject(ot, name, material)
	}

	for scanner.Scan() {
		currentLine++
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		t, value := parseLineType(line)

		// Force GC and release mem to OS for >1 million
		// line source files, every million lines.
		//

		switch t {
		// comments - #
		case tComment:
			if currentObject == nil && len(root.MaterialLibraries) == 0 {
				root.Comments = append(root.Comments, value)
			} else if currentObject != nil && len(value) > 0 {
				currentObject.Comments = append(currentObject.Comments, value)
			}

			// mtl file ref - mtllib
		case tMtlLib:
			root.MaterialLibraries = append(root.MaterialLibraries, value)

			// geometry - v, vn, vt, vp
		case tVertext, tNormal, tUV, tParam:
			if _, err := root.Geometry.readValue(t, value, true); err != nil {
				return nil, currentLine, wrapErrorLine(err, currentLine)
			}

			// object, group - o, g
		case tChildObject, tChildGroup:
			currentObjectName = value
			currentObjectChildIndex = 0
			// inherit currently declared material
			currentObject = root.createObject(t, currentObjectName, currentMaterial)

			// object: material - usemtl
		case tMtlUse:
			// obj files can define multiple materials inside a single object/group.
			// usually these are small face groups that kill performance on 3D engines
			// as they have to render hundreds or thousands of meshes with the same material,
			// each mesh containing a few faces.
			//
			// this app will convert all these "multi material" objects into
			// separate object, later merging all meshes with the same material into
			// a single draw call geometry.
			//
			// this might be undesirable for certain users, renderers and authoring software,
			// in this case don't use this simplified on your obj files. simple as that.

			// only fake if an object has been declared and ...
			if currentObject != nil {
				// only fake if the current object has declared vertex data (faces etc.)
				// and the material name actually changed (encountering the same usemtl
				// multiple times in a row would be rare, but check for completeness)
				if len(currentObject.VertexData) > 0 && currentObject.Material != value {
					currentObject = fakeObject(value)
				}
			}

			// store material value for inheriting
			currentMaterial = value

			// set material to current object
			if currentObject != nil {
				currentObject.Material = currentMaterial
			}

			// object: faces - f, l, p
		case tFace, tLine, tPoint:
			// most tools support the file not defining a o/g prior to face declarations.
			// I'm not sure if the spec allows not declaring any o/g.
			// Our data structures and parsing however requires objects to put the faces into,
			// create a default object that is named after the input file (without suffix).
			if currentObject == nil {
				currentObject = root.createObject(tChildObject, "default", currentMaterial)
			}
			vd, vdErr := currentObject.readVertexData(t, value, true)
			if vdErr != nil {
				return nil, currentLine, wrapErrorLine(vdErr, currentLine)
			}
			// attach current smooth group and reset it
			if len(currentSmoothGroup) > 0 {
				vd.setMeta(tSmoothingGroup, currentSmoothGroup)
				currentSmoothGroup = ""
			}

			// smooth group - s
		case tSmoothingGroup:
			// smooth group can change mid vertex data declaration
			// so it is attached to the vertex data instead of current object directly
			currentSmoothGroup = value

			// unknown
		case tUnkown:
			return nil, currentLine, wrapErrorLine(fmt.Errorf("unsupported line '%s'", line), currentLine)
		default:
			return nil, currentLine, wrapErrorLine(fmt.Errorf("unsupported line '%s'", line), currentLine)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, currentLine, err
	}
	return root, currentLine, nil
}

func wrapErrorLine(err error, lineNum int) error {
	return fmt.Errorf("line:%d %s", lineNum, err.Error())
}

func parseLineType(str string) (dataType, string) {
	value := ""
	if i := strings.Index(str, " "); i != -1 {
		value = strings.TrimSpace(str[i+1:])
		str = str[0:i]
	}
	return typeFromString(str), value
}
