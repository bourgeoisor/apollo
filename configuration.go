package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
)

// Configuration is a map of all the configuration values.
type Configuration struct {
	options map[string]string
}

// NewConfiguration creates a new Configuration with default values and returns it.
func newConfiguration() *Configuration {
	options := map[string]string{
		"tabs-startup":   "",
		"rating-startup": "false",
		"debug":          "false",
	}

	c := &Configuration{
		options: options,
	}

	c.load()
	c.save()

	return c
}

// Load fetches the user's configuration options from a file.
func (c *Configuration) load() {
	path := configurationPath()
	cont, err := ioutil.ReadFile(path)
	if err != nil {
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

// Save takes the map of configuration options and saves it to a file.
func (c *Configuration) save() {
	cont, err := json.Marshal(c.options)
	if err != nil {
		log.Fatal(err)
	}

	path := configurationPath()
	err = ioutil.WriteFile(path, cont, 0644)
	if err != nil {
		log.Print(err)
	}
}

// Set changes the value of a specific configuration option. It then saves the configuration.
func (c *Configuration) set(option string, value string) error {
	if _, exist := c.options[option]; exist {
		c.options[option] = value
	} else {
		return errors.New("config: invalid option")
	}

	c.save()
	return nil
}

// Get returns the value of a specific configuration option.
func (c *Configuration) get(option string) string {
	if value, exist := c.options[option]; exist {
		return value
	}

	return ""
}

// Config returns the list of all the configuration values.
func (c *Configuration) config() []string {
	var s []string

	for key, value := range c.options {
		s = append(s, key+": "+value)
	}

	return s
}
