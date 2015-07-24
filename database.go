package main

import (
    "log"
    "io/ioutil"
    "encoding/json"
)

type Movie struct {
    Title string
}

type Database struct {
    Movies []Movie
}

func createDatabase() *Database {
    d := &Database{}

    d.load()
    d.save()

    return d
}

func (d *Database) load() {
    cont, err := ioutil.ReadFile("database.json")
    if err != nil {
        log.Fatal(err)
    }

    err = json.Unmarshal(cont, d)
    if err != nil {
        log.Fatal(err)
    }
}

func (d *Database) save() {
    cont, err := json.Marshal(d)
    if err != nil {
        log.Fatal(err)
    }

    err = ioutil.WriteFile("database.json", cont, 0644)
    if err != nil {
        log.Fatal(err)
    }
}
