package editor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// editor is the default editor used when $EDITOR is not assigned.
var editor = "vim"

func init() {
	if s := os.Getenv("EDITOR"); s != "" {
		editor = s
	}
}

// Read opens the default editor and returns the value.
func Read(content []byte) ([]byte, error) {
	return ReadEditor(editor, content)
}

// ReadEditor opens the editor and returns the value.
func ReadEditor(editor string, content []byte) ([]byte, error) {
	// tmpfile
	f, err := ioutil.TempFile("", ".LAMBDA_EDITOR-*-.json5")
	if err != nil {
		return nil, fmt.Errorf("creating tmpfile : %w", err)
	}
	defer os.Remove(f.Name())

	f.Write(content)

	// open editor
	cmd := exec.Command("sh", "-c", editor+" "+f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("executing : %w", err)
	}

	// read tmpfile
	b, err := ioutil.ReadFile(f.Name())
	if err != nil {
		return nil, fmt.Errorf("reading tmpfile : %w", err)
	}

	return b, nil
}
