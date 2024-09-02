package builder

import (
	"regexp"
	"strings"
	"test-data-builder/pkg/util"
)

type Field struct {
	Name string
	Type string
}

type DataStructure struct {
	Package string
	Name    string
	Fields  []Field
}

func NewDataStructure(packageName, name string) *DataStructure {
	return &DataStructure{
		Package: strings.Replace(packageName, "main.", "", 1),
		Name:    name,
	}
}

func (ds *DataStructure) AddFields(structDefinition string) {
	left := "{"
	right := "}"
	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(left) + `(.*?)` + regexp.QuoteMeta(right))
	matches := rx.FindAllString(structDefinition, -1)
	fieldsDefinition := strings.ReplaceAll(matches[0], "{", "")
	fieldsDefinition = strings.ReplaceAll(fieldsDefinition, "}", "")
	fields := strings.Fields(fieldsDefinition)

	for i := 0; i < len(fields); i += 2 {
		if i == len(fields)-1 {
			break
		}
		field := Field{}
		field.Name = fields[i]
		if util.IsReservedWord(field.Name) {
			field.Name += "_"
		}
		field.Type = fields[i+1]
		if !util.IsBasicType(field.Type) {
			field.AddPackage()
		}
		ds.Fields = append(ds.Fields, field)
	}
}

func (field *Field) AddPackage() {
	names := GetAllStructsNames()
	for _, name := range names {
		if strings.Contains(field.Type, name) {
			packageName := FindStructPackage(name)
			field.Type = strings.ReplaceAll(field.Type, name, packageName+"."+name)
			return
		}
	}
}
