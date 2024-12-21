package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHandlerParser(t *testing.T) {
	dir, _ := os.Getwd()

	t.Run("Should parse the package name", func(t *testing.T) {
		parser := HandlerParser{}
		signature, err := parser.Parse(filepath.Join(dir, "testdata", "handler_001.go"))

		if err != nil {
			t.Fatal(err)
			return
		}

		if signature.PackageName != "handler" {
			t.Errorf("Expected %s as package name got %s", "handler", signature.PackageName)
		}
	})

	t.Run("Should parse all the valid methods from the handler", func(t *testing.T) {
		parser := HandlerParser{}
		signature, err := parser.Parse(filepath.Join(dir, "testdata", "handler_002.go"))

		if err != nil {
			t.Fatal(err)
			return
		}

		if len(signature.Functions) != 5 {
			t.Errorf("Expected %s methods to be extracted got %d", "5", len(signature.Functions))
		}
	})
}
