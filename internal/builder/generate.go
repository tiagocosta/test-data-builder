package builder

import (
	"embed"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

const (
	TestDataBuilderFolder = "testdatabuilder"
	TestDataBuilderName   = "test-data-builder.go"
	TemplatesFolder       = "templates"
	TemplateName          = "test-builder.tmpl"
)

type Generator struct {
	File    embed.FS
	Builder *DataBuilder
}

func NewGenerator(file embed.FS) *Generator {
	return &Generator{
		File:    file,
		Builder: &DataBuilder{},
	}
}

func (gen *Generator) Generate() {

	err := removeOldFolder()
	if err != nil {
		panic(err)
	}

	err = createNewFolder()
	if err != nil {
		panic(err)
	}

	err = gen.findPackagesAndStructs()
	if err != nil {
		panic(err)
	}

	err = gen.createStructs()
	if err != nil {
		panic(err)
	}

	err = gen.createBuilderFromTemplate()
	if err != nil {
		panic(err)
	}
}

func (gen *Generator) findPackagesAndStructs() error {
	mapPackageStructs = make(map[string][]string)
	filesPaths, _ := loadPaths(".")
	rootPath, _ := os.Getwd()
	baseDir := filepath.Base(rootPath)
	r, err := regexp.Compile(`type\s+[A-Za-z0-9_.]+\s+struct(?s)(.*?)}`)
	if err != nil {
		return err
	}
	for _, filePath := range filesPaths {
		fileData, _ := loadFile(filePath)
		if r.MatchString(fileData) {
			pathSlice := strings.Split(filePath, "/")
			pathSlice = pathSlice[:len(pathSlice)-1]
			gen.Builder.AddImportPath("", baseDir+"/"+strings.Join(pathSlice, "/"))
			packageName, err := getPackageName(fileData)
			if err != nil {
				return err
			}
			mapPackageStructs[packageName] = r.FindAllString(fileData, -1)
		}
	}

	return nil
}

func (gen *Generator) createStructs() error {
	for packageName, structs := range mapPackageStructs {
		for _, st := range structs {
			structName, err := getStructName(st)
			if err != nil {
				return err
			}
			ds := NewDataStructure(packageName, structName)
			ds.AddFields(st)
			gen.Builder.Structs = append(gen.Builder.Structs, *ds)
		}
	}

	return nil
}

func (gen *Generator) createBuilderFromTemplate() error {
	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"Contains":   strings.Contains,
		"TrimSuffix": strings.TrimSuffix,
		"TrimPrefix": strings.TrimPrefix,
	}

	tmpl, err := template.New(TemplateName).Funcs(funcMap).ParseFS(gen.File, TemplatesFolder+"/"+TemplateName)

	if err != nil {
		panic(err)
	}

	out, err := os.Create(TestDataBuilderFolder + "/" + TestDataBuilderName)
	if err != nil {
		return err
	}

	err = tmpl.Execute(out, gen.Builder)
	if err != nil {
		return err
	}

	return nil
}

func removeOldFolder() error {
	err := os.RemoveAll(TestDataBuilderFolder)
	if err != nil {
		return err
	}

	return nil
}

func createNewFolder() error {
	err := os.Mkdir(TestDataBuilderFolder, 0755)
	if err != nil {
		return err
	}

	return nil
}

func loadPaths(rootPath string) ([]string, error) {
	var paths []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !avoidPath(path) && !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func avoidPath(path string) bool {
	return strings.Contains(path, "gerados") ||
		strings.Contains(path, "graphql") ||
		strings.Contains(path, "mock") ||
		strings.Contains(path, ".git") ||
		strings.Contains(path, "_test")
}

func loadFile(filePath string) (string, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(fileData), nil
}

func getPackageName(fileData string) (string, error) {
	r, err := regexp.Compile(`package\s+[A-Za-z0-9_.]+`)
	if err != nil {
		return "", err
	}
	matches := r.FindAllString(fileData, -1)
	return strings.TrimSpace(strings.Replace(matches[0], "package", "", 1)), nil
}

func getStructName(structDefinition string) (string, error) {
	r, err := regexp.Compile(`type\s+[A-Za-z0-9_.]+\s+struct`)
	if err != nil {
		return "", err
	}
	matches := r.FindAllString(structDefinition, -1)
	return strings.TrimSpace(strings.Replace(strings.Replace(matches[0], "type", "", 1), "struct", "", 1)), nil
}
