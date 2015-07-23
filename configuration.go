package main

import (
    "log"
    "io/ioutil"
    "encoding/json"
)

type Configuration struct {
    Test string
}

func createConfiguration() *Configuration {
    c := &Configuration{
        Test: "nope.avi",
    }

    c.load()

    return c
}

func (c *Configuration) load() {
    cont, err := ioutil.ReadFile("configuration.json")
    if err != nil {
        log.Fatal(err)
    }

    err = json.Unmarshal(cont, c)
    if err != nil {
        log.Fatal(err)
    }
}

func (c *Configuration) save() {

}
