package main

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path"

	"github.com/schafer14/mlg/internal/ast"
	"github.com/schafer14/mlg/internal/editor"
	"github.com/schafer14/mlg/internal/generate"
	"golang.org/x/mod/modfile"
)

//go:embed templates/lambda-file-template.json5
var lambdaFileTemplate []byte

//go:embed templates/src
var srcDir embed.FS

func main() {

	// TODO: prompt if the user is okay with selection
	// TODO: loop on json5 error
	lambdaFile, err := editor.Read(lambdaFileTemplate)
	if err != nil {
		panic(err)
	}

	// TODO: check there is nothing staged.

	tree, err := ast.Unmarshal(lambdaFile)
	if err != nil {
		panic(err)
	}

	if tree.Name == "" {
		fmt.Println("No name provided")
		return
	}

	root, err := findRoot()
	if err != nil {
		panic(err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	goModPath := path.Join(workingDir, "go.mod")
	modFile, err := ioutil.ReadFile(goModPath)
	if err != nil {
		fmt.Println("must be run from directory containing a `go.mod` file")
		return
	}
	module, err := modfile.Parse("go.mod", modFile, nil)
	if err != nil {
		panic(fmt.Errorf("parsing go.mod file: %w", err))
	}
	tree.Module = module

	g := &generate.Gen{RootPath: root, Details: tree}

	fs.WalkDir(srcDir, ".", func(path string, d fs.DirEntry, err error) error {

		if err != nil {
			fmt.Println(err)
			return err
		}

		if d.IsDir() {
			return nil
		}

		fileContent, err := srcDir.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return g.GenerateMain(path, fileContent)
	})
}

// May want to make this more robust in the future.
func findRoot() (string, error) {
	return os.Getwd()
}
