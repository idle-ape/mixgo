package config

import (
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	codecMap = make(map[string]Codec)
	l        = sync.RWMutex{}
)

func init() {
	RegisterCodec(&yamlCodec{})
}

// Codec defines codec interface.
type Codec interface {
	// Name codec name
	Name() string

	// Unmarshal unmarshal data
	Unmarshal([]byte, interface{}) error
}

type yamlCodec struct{}

func (y *yamlCodec) Name() string {
	return "yaml"
}

func (y *yamlCodec) Unmarshal(data []byte, out interface{}) error {
	return yaml.Unmarshal(data, out)
}

// RegisterCodec register codec
func RegisterCodec(c Codec) {
	l.Lock()
	codecMap[c.Name()] = c
	l.Unlock()
}

// GetCodec get a Codec by name
func GetCodec(name string) Codec {
	l.RLock()
	defer l.RUnlock()
	if c, ok := codecMap[name]; ok {
		return c
	}
	return nil
}
