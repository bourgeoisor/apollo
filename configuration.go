package main

import (
    "log"
    "io/ioutil"
    "encoding/json"
)

type Configuration struct {
    Test string
    Testb bool
}

func createConfiguration() *Configuration {
    c := &Configuration{
        Test: "nope.avi",
        Testb: true,
    }

    c.load()
    c.save()

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
    cont, err := json.Marshal(c)
    if err != nil {
        log.Fatal(err)
    }

    err = ioutil.WriteFile("configuration.json", cont, 0644)
    if err != nil {
        log.Fatal(err)
    }
}
