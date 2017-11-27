package loader

type vertexData struct {
	Type         dataType
	Declarations []*declaration

	meta map[dataType]string
}

func (f *vertexData) setMeta(t dataType, value string) {
	if f.meta == nil {
		f.meta = make(map[dataType]string)
	}
	f.meta[t] = value
}

func (f *vertexData) index(index int) *declaration {
	if index >= 0 && index <= 3 {
		for index >= len(f.Declarations) {
			f.Declarations = append(f.Declarations, &declaration{})
		}
		return f.Declarations[index]
	}
	return nil
}
