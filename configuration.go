package main

import (
	"os"
	"encoding/json"
	"log"
	"strconv"
)

type Configuration struct {
	Url string
	Limit int
	DisplayProgress bool
	UrlsToIgnore []string
}

func getConfiguration() Configuration {
	config := Configuration{}

	if len(os.Args) == 2 {
		config = parseConfigFile(os.Args[1])
	} else {
		config = parseArgs(os.Args)
	}

	if config.Url[len(config.Url) - 1:] != "/" {
		config.Url = config.Url + "/"
	}

	return config
}

func parseConfigFile(filePath string) Configuration {
	file, _ := os.Open(filePath)
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)
	if err != nil {
		log.Fatal("error:", err)
	}
	return config
}

func parseArgs(args []string) Configuration {
	if len(args) != 3 && len(args) != 4 {
		log.Fatal("Invalid number of arguments.\nUsage example:\n\"go run link-checker.go http://example.com 100\"")
	}

	url := args[1]
	limit, _ := strconv.Atoi(args[2])
	displayProgress := true

	if len(args) > 3 {
		displayProgress = args[3] == "1" || args[3] == "true"
	}

	return Configuration {
		Url: url,
		Limit: limit,
		DisplayProgress: displayProgress,
	}
}
