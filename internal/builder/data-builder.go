package builder

type Import struct {
	Alias string
	Path  string
}

type DataBuilder struct {
	Imports []Import
	Structs []DataStructure
}

var mapPackageStructs map[string][]string

func (tdb *DataBuilder) AddImportPath(alias, path string) {
	if tdb.CanAddImportPath(path) && !tdb.ContainsImportPath(path) {
		tdb.Imports = append(tdb.Imports, Import{
			Alias: alias,
			Path:  path,
		})
	}
}

func (tdb *DataBuilder) CanAddImportPath(path string) bool {
	return path != "main"
}

func (tdb *DataBuilder) ContainsImportPath(path string) bool {
	for _, v := range tdb.Imports {
		if v.Path == path {
			return true
		}
	}
	return false
}

func FindStructPackage(structName string) string {
	for packageName, structsDefinitions := range mapPackageStructs {
		for _, definition := range structsDefinitions {
			name, _ := getStructName(definition)
			if name == structName {
				return packageName
			}
		}
	}
	return ""
}

func GetAllStructsNames() []string {
	var names = []string{}
	for _, structsDefinitions := range mapPackageStructs {
		for _, definition := range structsDefinitions {
			name, _ := getStructName(definition)
			names = append(names, name)
		}
	}
	return names
}
