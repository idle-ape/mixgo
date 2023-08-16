package config

import (
	"fmt"
	"sync"
)

var DefaultConfigLoader = newMixgoConfigLoader()

// Config defines the common config interface
type Config interface {
	// RawData returns the raw data of then path
	RawData() []byte
}

// MixgoConfigLoader is a config loader with cache
type MixgoConfigLoader struct {
	cache map[string]Config
	rwl   sync.RWMutex
}

func newMixgoConfigLoader() *MixgoConfigLoader {
	return &MixgoConfigLoader{cache: map[string]Config{}, rwl: sync.RWMutex{}}
}

func (cl *MixgoConfigLoader) Load(path string, opts ...LoadOption) (Config, error) {
	mc := newMixgoConfig(path)
	for _, o := range opts {
		o(mc)
	}
	if mc.decoder == nil {
		return nil, fmt.Errorf("codec not exist")
	}
	if mc.p == nil {
		return nil, fmt.Errorf("provider not exist")
	}
	if mc.r == nil {
		return nil, fmt.Errorf("reciver not exist")
	}

	key := fmt.Sprintf("%s.%s.%s", mc.p.Name(), mc.decoder.Name(), path)
	cl.rwl.RLock()
	// get from cache
	if c, ok := cl.cache[key]; ok {
		cl.rwl.RUnlock()
		return c, nil
	}
	cl.rwl.RUnlock()

	if err := mc.Load(); err != nil {
		return nil, err
	}

	cl.rwl.Lock()
	cl.cache[key] = mc
	cl.rwl.Unlock()

	mc.p.Watch(func(p string, data []byte) {
		if p == path {
			fmt.Printf("file changed, p: %s, data: %s\n", p, string(data))
			cl.rwl.Lock()
			delete(cl.cache, key)
			cl.rwl.Unlock()
			if !mc.disableWatch {
				mc.Load()
			}
		}
	})
	return mc, nil
}

// Load returns the config specificed by input parameter.
func Load(path string, opts ...LoadOption) (Config, error) {
	return DefaultConfigLoader.Load(path, opts...)
}

func newMixgoConfig(path string) *MixgoConfig {
	return &MixgoConfig{
		p:       GetProvider("file"),
		path:    path,
		decoder: &yamlCodec{},
	}
}

type MixgoConfig struct {
	p            Provider
	r            interface{}
	path         string
	disableWatch bool
	decoder      Codec
	rawData      []byte
}

func (m *MixgoConfig) Load() error {
	data, err := m.p.Read(m.path)
	if err != nil {
		return fmt.Errorf("failed to load %s: %s", m.path, err.Error())
	}
	m.rawData = data
	if err := m.decoder.Unmarshal(m.rawData, m.r); err != nil {
		return fmt.Errorf("failed to parse %s: %s", m.path, err.Error())
	}
	return nil
}

// RawData returns the raw data of then path
func (m *MixgoConfig) RawData() []byte {
	return m.rawData
}
