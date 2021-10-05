package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

const configFile = "./config.json"
const portEnvKey = "PORT"
const redirectPortEnvKey = "REDIRECTPORT"
const redirectHostEnvKey = "REDIRECTHOST"

type configData struct {
	Port              string
	DefaultPage       string
	AssetsDir         string
	UsersDb           string
	RedirectHost      string
	RedirectPort      string
	NumFailedAttempts uint
	FailedTimeout     uint
	CodeValiditySecs  uint
	redirectURL       string
}

var once sync.Once
var cf *configData

func config() *configData {
	if cf == nil {
		once.Do(
			func() {
				cf = &configData{}
				cf.load()
				cf.redirectURL = fmt.Sprintf("%s:%s",
					cf.RedirectHost, cf.RedirectPort)
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
	setEnvValue(&cf.Port, portEnvKey)
	setEnvValue(&cf.RedirectHost, redirectHostEnvKey)
	setEnvValue(&cf.RedirectPort, redirectPortEnvKey)
}

func setEnvValue(field *string, key string) {
	v := os.Getenv(key)
	if v != "" {
		*field = v
	}
}
