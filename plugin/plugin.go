package plugin

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

var plugins map[string]Plugin

// Plugins config for plugin
type Plugins map[string]yaml.Node

type Plugin interface {
	Name() string
	Setup(yaml.Node) error
}

// Setup plugin setup
func (p Plugins) Setup() error {
	for name, node := range p {
		pl := Get(name)
		if pl == nil {
			return fmt.Errorf("plugin for %s not registered", name)
		}
		if err := pl.Setup(node); err != nil {
			return err
		}
	}
	return nil
}

// Register register a plugin
func Register(p Plugin) {
	if plugins == nil {
		plugins = map[string]Plugin{}
	}
	plugins[p.Name()] = p
}

// Get get a Plugin by name
func Get(name string) Plugin {
	return plugins[name]
}
