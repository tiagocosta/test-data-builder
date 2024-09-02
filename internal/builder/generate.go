package builder

import (
	"embed"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

func Generate(file embed.FS) {
	err := os.RemoveAll("testdatabuilder")
	if err != nil {
		log.Println("error creating test data builder folder", err)
		return
	}

	err = os.Mkdir("testdatabuilder", 0755)
	if err != nil {
		log.Println("error creating test data builder folder", err)
		return
	}

	dataBuilder := DataBuilder{}
	mapPackageStructs = make(map[string][]string)
	filesPaths, _ := loadPaths(".")
	rootPath, _ := os.Getwd()
	baseDir := filepath.Base(rootPath)
	r, _ := regexp.Compile(`type\s+[A-Za-z0-9_.]+\s+struct(?s)(.*?)}`)
	for _, filePath := range filesPaths {
		fileData, _ := loadFile(filePath)
		if r.MatchString(fileData) {
			pathSlice := strings.Split(filePath, "/")
			pathSlice = pathSlice[:len(pathSlice)-1]
			dataBuilder.AddImportPath("", baseDir+"/"+strings.Join(pathSlice, "/"))
			mapPackageStructs[getPackageName(fileData)] = r.FindAllString(fileData, -1)
		}
	}

	for packageName, structs := range mapPackageStructs {
		for _, st := range structs {
			ds := NewDataStructure(packageName, getStructName(st))
			ds.AddFields(st)
			dataBuilder.Structs = append(dataBuilder.Structs, *ds)
		}
	}

	funcMap := template.FuncMap{
		"ToLower":    strings.ToLower,
		"Contains":   strings.Contains,
		"TrimSuffix": strings.TrimSuffix,
		"TrimPrefix": strings.TrimPrefix,
	}

	var tmplFile1 = "templates/test-builder.tmpl"
	tmpl, err := template.New("test-builder.tmpl").Funcs(funcMap).ParseFS(file, tmplFile1)

	if err != nil {
		panic(err)
	}

	out, err := os.Create("testdatabuilder/test-data-builder.go")
	if err != nil {
		log.Println("error creating test data builder file", err)
		return
	}

	err = tmpl.Execute(out, dataBuilder)
	if err != nil {
		panic(err)
	}
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

func getPackageName(fileData string) string {
	r, _ := regexp.Compile(`package\s+[A-Za-z0-9_.]+`)
	matches := r.FindAllString(fileData, -1)
	return strings.TrimSpace(strings.Replace(matches[0], "package", "", 1))
}

func getStructName(structDefinition string) string {
	r, _ := regexp.Compile(`type\s+[A-Za-z0-9_.]+\s+struct`)
	matches := r.FindAllString(structDefinition, -1)
	return strings.TrimSpace(strings.Replace(strings.Replace(matches[0], "type", "", 1), "struct", "", 1))
}
