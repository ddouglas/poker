package poker

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed templates
var templates embed.FS

const templatesDirectory = "templates"

func TemplateFS(environment Environment) fs.FS {

	if !environment.IsProduction() {
		return getLocalTemplateFS()
	}

	subFS, err := fs.Sub(templates, "templates")
	if err != nil {
		panic(fmt.Sprintf("fs.Sub: %s", err))
	}

	return subFS

}

func getLocalTemplateFS() fs.FS {

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("failed to get current working directory: %s", err))
	}

	return os.DirFS(filepath.Join(cwd, templatesDirectory))

}
