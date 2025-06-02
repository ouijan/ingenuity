package cli

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

var commentTag = "ingenuity"

func parseSourceCode(sourcePath string) error {
	fmt.Printf("Parsing Source Code from %s\n", sourcePath)
	fset := token.NewFileSet()

	srcFile := "main.go"
	src, err := os.ReadFile(path.Join(sourcePath, srcFile))
	if err != nil {
		return err
	}

	f, err := parser.ParseFile(fset, srcFile, src, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Print(fset, f)

	for _, decl := range f.Decls {
		decl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		handleGenDecl(decl)
	}

	return nil
}

func handleGenDecl(decl *ast.GenDecl) {
	switch getTypeFromDoc(decl.Doc) {
	case "component":
		handleComponent(decl)
	case "enum":
		handleEnum(decl)
	}
}

func handleComponent(decl *ast.GenDecl) *TiledCustomClass {
	componentName := ""
	members := []TiledClassMember{}

	for _, spec := range decl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		fmt.Printf("Component: %v\n", typeSpec.Name)
		componentName = typeSpec.Name.Name
		for _, field := range typeSpec.Type.(*ast.StructType).Fields.List {
			tags := parseJsonTag(field)
			for _, name := range field.Names {
				fmt.Printf(" - %s %s = %v \n", name.Name, tags["json"], field.Type)
				// TODO: infer type, set proprty type if is an enum

				members = append(
					members,
					NewTiledClassMember(
						name.Name,
						fmt.Sprintf("%s", field.Type),
						tags["json"],
						nil,
					),
				)
			}
		}
	}

	component := NewTiledCustomClass(0, componentName, members)
	componentJson, err := json.Marshal(component)
	if err != nil {
		return nil
	}
	fmt.Printf("component json -> %s", componentJson)
	return &component
}

func findEnumType(decl *ast.GenDecl) string {
	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if valueSpec.Type != nil {
			return fmt.Sprintf("%s", valueSpec.Type)
		}
	}
	return ""
}

func handleEnum(decl *ast.GenDecl) *TiledCustomEnum {
	enumType := findEnumType(decl)
	if enumType == "" {
		return nil
	}

	values := []string{}

	fmt.Printf("Enum: %v\n", enumType)
	for _, spec := range decl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		for _, name := range valueSpec.Names {
			fmt.Printf(" - %s = %v\n", name.Name, name.Obj.Data)
			values = append(values, name.Name)
		}
	}

	enum := NewTiledCustomEnum(0, enumType, values)
	enumJson, err := json.Marshal(enum)
	if err != nil {
		return nil
	}
	fmt.Printf("enum json -> %s", enumJson)
	return &enum
}

func getTypeFromDoc(doc *ast.CommentGroup) string {
	if doc == nil {
		return ""
	}
	for _, comment := range doc.List {
		segments := parseComment(comment.Text)
		if len(segments) >= 2 {
			if segments[0] == commentTag {
				return segments[1]
			}
		}
	}
	return ""
}

func parseComment(comment string) []string {
	slashless := strings.Replace(comment, "//", "", -1)
	segments := strings.Split(slashless, ":")
	for i, segment := range segments {
		segments[i] = strings.ToLower(strings.TrimSpace(segment))
	}
	return segments
}

func parseJsonTag(field *ast.Field) map[string]string {
	if field.Tag == nil {
		return make(map[string]string)
	}
	tag := field.Tag.Value
	tag = strings.Trim(tag, "`")
	tags := strings.Split(tag, " ")
	tagMap := make(map[string]string)
	for _, tag := range tags {
		tag = strings.Trim(tag, " ")
		if tag == "" {
			continue
		}
		kv := strings.Split(tag, ":")
		if len(kv) == 2 {
			tagMap[kv[0]] = kv[1]
		}
	}
	return tagMap
}
