package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

const configFile = "./config.json"
const portEnvKey = "PORT"

type configData struct {
	Port              string
	DefaultPage       string
	AssetsDir         string
	UsersDb           string
	RedirectURL       string
	NumFailedAttempts uint
	FailedTimeout     uint
}

var once sync.Once
var cf *configData

func config() *configData {
	if cf == nil {
		once.Do(
			func() {
				cf = &configData{}
				cf.load()
			})
	}
	return cf
}

func (cf *configData) load() {
	bytes, err := ioutil.ReadFile(configFile)
	checkerr(err)
	err = json.Unmarshal(bytes, cf)
	// <TO DO> chang unmarshall to map[string]string
	// iterate and check each key is loaded to not empty
	// assign to 'cf' fields
	//log.Printf("cf = %v\n", cf)
	checkerr(err)
	port := os.Getenv(portEnvKey)
	if port != "" {
		cf.Port = port
	}
}
