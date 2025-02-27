package cli

import (
	"encoding/json"
	"os"
	"strings"
)

type (
	TiledCommand struct {
		Arguments         string `json:"arguments"`
		Command           string `json:"command"`
		Enabled           bool   `json:"enabled"`
		Name              string `json:"name"`
		SaveBeforeExecute bool   `json:"saveBeforeExecute"`
		Shortcut          string `json:"shortcut"`
		ShowOutput        bool   `json:"showOutput"`
		WorkingDirectory  string `json:"workingDirectory"`
	}
	TiledCustomProperty struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Type string `json:"type"`
		// Class Properties
		Color    string             `json:"color,omitempty"`
		DrawFill bool               `json:"drawFill,omitempty"`
		Members  []TiledClassMember `json:"members,omitempty"`
		UseAs    []string           `json:"useAs,omitempty"`
		// Enum Properties
		StorageType   string   `json:"storageType,omitempty"`
		Values        []string `json:"values,omitempty"`
		ValuesAsFlags bool     `json:"valuesAsFlags,omitempty"`
	}
	TiledClassMember struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		PropertyType string `json:"propertyType"`
		Value        any    `json:"value"`
	}
	TiledCustomClass struct {
		Id       int                `json:"id"`
		Name     string             `json:"name"`
		Type     string             `json:"type"`
		Color    string             `json:"color"`
		DrawFill bool               `json:"drawFill"`
		Members  []TiledClassMember `json:"members"`
		UseAs    []string           `json:"useAs"`
	}
	TiledCustomEnum struct {
		Id            int      `json:"id"`
		Name          string   `json:"name"`
		Type          string   `json:"type"`
		StorageType   string   `json:"storageType"`
		Values        []string `json:"values"`
		ValuesAsFlags bool     `json:"valuesAsFlags"`
	}
	TiledProject struct {
		Commands      []TiledCommand        `json:"commands"`
		PropertyTypes []TiledCustomProperty `json:"propertyTypes"`
	}
)

func NewTiledClassMember(
	name string,
	memberType string,
	propertyType string,
	value any,
) TiledClassMember {
	return TiledClassMember{
		Name:         name,
		Type:         memberType,
		PropertyType: propertyType,
		Value:        value,
	}
}

func NewTiledCustomClass(id int, name string, members []TiledClassMember) TiledCustomClass {
	return TiledCustomClass{
		Id:       id,
		Name:     name,
		Type:     "class",
		Color:    "#ff0000",
		DrawFill: true,
		Members:  members,
		UseAs: []string{
			"property",
			"map",
			"layer",
			"object",
			"tile",
			"tileset",
			"wangcolor",
			"wangset",
			"project",
		},
	}
}

func NewTiledCustomEnum(id int, name string, values []string) TiledCustomEnum {
	return TiledCustomEnum{
		Id:            id,
		Name:          name,
		Type:          "enum",
		StorageType:   "int",
		Values:        values,
		ValuesAsFlags: false,
	}
}

func readTiledProject(path string) (*TiledProject, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	project := TiledProject{}
	decoder := json.NewDecoder(strings.NewReader(string(data)))
	decoder.Decode(&project)
	return &project, err
}

func writeTiledProject(path string, project *TiledProject) error {
	data, err := json.Marshal(project)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
