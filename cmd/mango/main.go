package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eryk-vieira/mango/internal"
	"github.com/eryk-vieira/mango/internal/types"
)

func main() {
	command, err := getFrameworkCommand()

	if err != nil {
		log.Fatalln(err)
		return
	}

	// Remove the index 1 of the array.
	// This should be done because the framework command does not follow the flag pattern
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)

	rootDir := flag.String("root", "./", "Project root folder")
	configFile := flag.String("config", "mango.json", "Config file")

	flag.Parse()

	if command == "help" {
		flag.Usage()

		return
	}

	err = os.Chdir(*rootDir)

	if err != nil {
		log.Fatalln(err)
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

func getFrameworkCommand() (string, error) {
	if len(os.Args) == 0 {
		return "", errors.New("at least one command should be specified")
	}

	return os.Args[1:][0], nil
}
