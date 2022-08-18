package generate

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/schafer14/mlg/internal/ast"
)

type Gen struct {
	RootPath string
	Details  ast.AST
}

// GenerateMain a main file.
func (g *Gen) GenerateMain(pathTemp string, content []byte) error {

	filePathBytes, err := tempToString(pathTemp, g.Details)
	if err != nil {
		return err
	}
	filePath := string(filePathBytes)
	filePath = strings.ReplaceAll(filePath, "templates/src/", "")
	filePath = strings.ReplaceAll(filePath, ".txt", "")

	dirPath := path.Dir(filePath)
	if err := os.MkdirAll(dirPath, 0750); err != nil {
		log.Fatal(err)
	}

	fileContents, err := tempToString(string(content), g.Details)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filePath, fileContents, 0664)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File created: ", filePath)

	return nil
}

func tempToString(tmpStr string, d interface{}) ([]byte, error) {

	t, err := template.New("").Funcs(funcMap).Parse(tmpStr)
	if err != nil {
		return []byte{}, err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, d); err != nil {
		return []byte{}, err
	}

	return tpl.Bytes(), nil
}

var funcMap = map[string]interface{}{
	"awsPolicies":   awsPolicies,
	"hasDependency": hasDependency,
}

func awsPolicies(dep ast.Dependency) string {

	switch dep.Type {
	case "dynamo":
		pols := []string{}
		for _, p := range dep.Policies {
			pols = append(pols, fmt.Sprintf(`"dynamodb:%s"`, p))
		}
		return fmt.Sprintf(`{
			Action: [%s],
			Resource: [ table.arn ],
			Effect: "Allow",
		},`, strings.Join(pols, ", "))
	case "eventBridge":
		pols := []string{}
		for _, p := range dep.Policies {
			pols = append(pols, fmt.Sprintf(`"events:%s"`, p))
		}
		return fmt.Sprintf(`{
			Action: [%s],
			Resource: [ eventBus.arn ],
			Effect: "Allow",
		},`, strings.Join(pols, ", "))
	}

	return ""
}

func hasDependency(deps []ast.Dependency, depType string) bool {

	for _, d := range deps {
		if d.Type == depType {
			return true
		}
	}

	return false
}
