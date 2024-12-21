package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eryk-vieira/mango/internal"
	"github.com/eryk-vieira/mango/internal/types"
)

func main() {
	command := getCommand(os.Args[1:])
	workDir := os.Args[3:]

	if len(workDir) > 0 {
		err := os.Chdir(workDir[0])

		if err != nil {
			fmt.Printf("Error changing directory: %v\n", err)
			return
		}
	}

	settings := parseSettings(os.Args[2:])

	if command == "build" {
		internal.Build(settings)

		return
	}

	if command == "run" {
		internal.Run(settings)

		return
	}
}

func parseSettings(args []string) *types.Settings {
	var settings types.Settings

	var configFile string = "nextgo.json"

	if len(args) > 0 {
		configFile = args[0]
	}

	jsonFile, err := os.ReadFile(configFile)

	if err != nil {
		panic(fmt.Sprintf("%s file not found", configFile))
	}

	err = json.Unmarshal(jsonFile, &settings)

	if err != nil {
		panic(err)
	}

	return &settings
}

func getCommand(args []string) string {
	return args[0]
}
