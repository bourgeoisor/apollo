package main

import (
    "log"
    "os"
    "io/ioutil"
    "encoding/json"
)

type Configuration struct {
    options map[string]string
}

func createConfiguration() *Configuration {
    c := &Configuration{
    }

    c.load()
    c.save()

    return c
}

func (c *Configuration) load() {
    path := os.Getenv("HOME") + "/.config/apollo/configuration.json"
    cont, err := ioutil.ReadFile(path)
    if err != nil {
        log.Print(err)
        return
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

    path := os.Getenv("HOME") + "/.config/apollo/configuration.json"
    err = ioutil.WriteFile(path, cont, 0644)
    if err != nil {
        log.Print(err)
    }
}

func (c *Configuration) set(option string, value string) {
    if _, ok := c.options[option]; ok {
        c.options[option] = value
    }
}

func (c *Configuration) get(option string) string {
    if val, ok := c.options[option]; ok {
        return val
    }
    return ""
}

func (c *Configuration) getAll() []string {
    var s []string

    return s
}
