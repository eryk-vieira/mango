package build

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/eryk-vieira/mango/internal/build/templates"
	"github.com/eryk-vieira/mango/internal/types"
)

const (
	frameworkTmpDir = "/mango/"
)

// TemplateData holds data to fill the ServerTemplate
type TemplateData struct {
	Imports []ImportData
	Routes  []RouteData
	Port    string
}

// ImportData represents data for each import
type ImportData struct {
	HandlerPackage string
	PackagePath    string
}

// RouteData represents data for each route
type RouteData struct {
	Method        string
	Pattern       string
	Handler       string
	HandlerMethod string
}

type serverBuilder struct {
	Settings  *types.Settings
	debugMode bool
}

func (b *serverBuilder) Build(routes []Route, errList *[]Errors) {
	tempPath := filepath.Join(os.TempDir(), frameworkTmpDir)

	defer os.RemoveAll(tempPath)

	isDebug, err := strconv.ParseBool(os.Getenv("NEXTGO_DEBUG_MODE"))

	if err != nil {
		b.debugMode = false
	} else if isDebug {
		b.debugMode = true
	}

	err = os.Mkdir(tempPath, fs.ModePerm)

	if err != nil {
		*errList = append(*errList, Errors{
			FilePath: "",
			Error:    err,
		})

		return
	}

	data := TemplateData{}

	for i, route := range routes {
		data.Imports = append(data.Imports, ImportData{
			PackagePath:    strings.TrimSuffix(route.FilePath, fmt.Sprintf("/%s.go", b.Settings.HTTP.HandlerName)),
			HandlerPackage: fmt.Sprintf("%s_%d", route.PackageName, i),
		})

		data.Routes = append(data.Routes, RouteData{
			Method:        route.Method,
			Pattern:       strings.TrimSuffix(route.Pattern, "/"),
			Handler:       fmt.Sprintf("%s_%d", route.PackageName, i),
			HandlerMethod: route.Method,
		})
	}

	data.Port = "8080"

	workDir, err := os.Getwd()

	if err != nil {
		*errList = append(*errList, Errors{
			FilePath: "",
			Error:    err,
		})

		return
	}

	tmpl := template.Must(template.New("server").Parse(templates.ServerTemplate))

	file, err := os.Create(filepath.Join(tempPath, "main.go"))

	if err != nil {
		*errList = append(*errList, Errors{
			FilePath: "",
			Error:    err,
		})

		return
	}

	defer file.Close()

	err = tmpl.Execute(file, data)

	if err != nil {
		*errList = append(*errList, Errors{
			FilePath: "",
			Error:    err,
		})

		return
	}

	err = b.copyFolder("./", tempPath)

	if err != nil {
		*errList = append(*errList, Errors{
			FilePath: "",
			Error:    err,
		})

		return
	}

	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = tempPath

	output, err := cmd.CombinedOutput()

	if err != nil {
		*errList = append(*errList, Errors{
			FilePath: "Error installing dependencies",
			Error:    errors.New(string(output)),
		})

		return
	}

	cmd = exec.Command("go", "build", "-o", filepath.Join(workDir, ".dist")+"/server", tempPath)
	cmd.Dir = tempPath

	if b.debugMode {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	output, err = cmd.CombinedOutput()

	if err != nil {
		*errList = append(*errList, Errors{
			FilePath: "Error building the application",
			Error:    errors.New(string(output)),
		})

		return
	}

}

func (*serverBuilder) copyFolder(src string, dst string) error {
	err := filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if strings.Contains(path, ".dist") {
			return nil
		}

		if info.Name() == ".dist" || info.Name() == "workspace" {
			return nil
		}

		if info.IsDir() {
			if err := os.MkdirAll(filepath.Join(dst, path), info.Mode()); err != nil {
				return err
			}

			return nil
		}

		go func() {
			copyFile(path, filepath.Join(dst, path))
		}()

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, sourceInfo.Mode())
}
