package main

import (
    "log"
    "os"
    "errors"
    "io/ioutil"
    "encoding/json"
)

type Configuration struct {
    options map[string]string
}

func newConfiguration() *Configuration {
    options := map[string]string{
        "autotag": "false",
        "movies_tab": "false",
        "series_tab": "false",
        "games_tab": "false",
        "books_tab": "false",
        "debug": "false",
    }

    c := &Configuration{
        options: options,
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

    var options map[string]string
    err = json.Unmarshal(cont, &options)
    if err != nil {
        log.Fatal(err)
    }

    for key, value := range options {
        if _, exist := c.options[key]; exist {
            c.options[key] = value
        }
    }
}

func (c *Configuration) save() {
    cont, err := json.Marshal(c.options)
    if err != nil {
        log.Fatal(err)
    }

    path := os.Getenv("HOME") + "/.config/apollo/configuration.json"
    err = ioutil.WriteFile(path, cont, 0644)
    if err != nil {
        log.Print(err)
    }
}

func (c *Configuration) set(option string, value string) error {
    if _, exist := c.options[option]; exist {
        c.options[option] = value
    } else {
        return errors.New("config: invalid option")
    }

    c.save()
    return nil
}

func (c *Configuration) get(option string) string {
    if value, exist := c.options[option]; exist {
        return value
    }

    return ""
}

func (c *Configuration) config() []string {
    var s []string

    for key, value := range c.options {
        s = append(s, key + ": " + value)
    }

    return s
}
