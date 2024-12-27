package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/eryk-vieira/mango/internal"
	"github.com/eryk-vieira/mango/internal/types"
)

func main() {
	command := getCommand(os.Args[1:])

	rootDir := flag.String("root", "./", "Project root folder")
	configFile := flag.String("config", "nextgo.json", "Config file")

	// Normalize os.Args removing the command.
	// This needs to be done because the command does not follow the flag pattern.
	// The flag should be parsed just after this execution
	normalizeArgs()

	flag.Parse()

	err := os.Chdir(*rootDir)

	if err != nil {
		panic(*rootDir + " Does not exists")
	}

	settings := parseSettings(*configFile)

	if command == "build" {
		internal.Build(settings)

		return
	}

	if command == "run" {
		internal.Run(settings)

		return
	}
}

func parseSettings(configPath string) *types.Settings {
	var settings types.Settings

	jsonFile, err := os.ReadFile(configPath)

	if err != nil {
		panic(fmt.Sprintf("%s file not found", configPath))
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

func normalizeArgs() {
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
}
